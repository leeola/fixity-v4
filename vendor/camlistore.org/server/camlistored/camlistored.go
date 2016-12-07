/*
Copyright 2011 Google Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// The camlistored binary is the Camlistore server.
package main // import "camlistore.org/server/camlistored"

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"camlistore.org/pkg/buildinfo"
	"camlistore.org/pkg/env"
	"camlistore.org/pkg/httputil"
	"camlistore.org/pkg/netutil"
	"camlistore.org/pkg/osutil"
	"camlistore.org/pkg/serverinit"
	"camlistore.org/pkg/webserver"

	// VM environments:
	"camlistore.org/pkg/osutil/gce" // for init side-effects + LogWriter

	// Storage options:
	_ "camlistore.org/pkg/blobserver/b2"
	"camlistore.org/pkg/blobserver/blobpacked"
	_ "camlistore.org/pkg/blobserver/cond"
	_ "camlistore.org/pkg/blobserver/diskpacked"
	_ "camlistore.org/pkg/blobserver/encrypt"
	_ "camlistore.org/pkg/blobserver/google/cloudstorage"
	_ "camlistore.org/pkg/blobserver/google/drive"
	_ "camlistore.org/pkg/blobserver/localdisk"
	_ "camlistore.org/pkg/blobserver/mongo"
	_ "camlistore.org/pkg/blobserver/proxycache"
	_ "camlistore.org/pkg/blobserver/remote"
	_ "camlistore.org/pkg/blobserver/replica"
	_ "camlistore.org/pkg/blobserver/s3"
	_ "camlistore.org/pkg/blobserver/shard"
	// Indexers: (also present themselves as storage targets)
	"camlistore.org/pkg/index"
	// KeyValue implementations:
	_ "camlistore.org/pkg/sorted/kvfile"
	_ "camlistore.org/pkg/sorted/leveldb"
	_ "camlistore.org/pkg/sorted/mongo"
	_ "camlistore.org/pkg/sorted/mysql"
	_ "camlistore.org/pkg/sorted/postgres"
	"camlistore.org/pkg/sorted/sqlite" // for sqlite.CompiledIn()

	// Handlers:
	_ "camlistore.org/pkg/search"
	_ "camlistore.org/pkg/server" // UI, publish, etc

	// Importers:
	_ "camlistore.org/pkg/importer/allimporters"

	// Licence:
	_ "camlistore.org/pkg/camlegal"

	"go4.org/legal"
	"go4.org/wkfs"

	"cloud.google.com/go/logging"
	"golang.org/x/crypto/acme/autocert"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

var (
	flagVersion    = flag.Bool("version", false, "show version")
	flagHelp       = flag.Bool("help", false, "show usage")
	flagLegal      = flag.Bool("legal", false, "show licenses")
	flagConfigFile = flag.String("configfile", "",
		"Config file to use, relative to the Camlistore configuration directory root. "+
			"If blank, the default is used or auto-generated. "+
			"If it starts with 'http:' or 'https:', it is fetched from the network.")
	flagListen      = flag.String("listen", "", "host:port to listen on, or :0 to auto-select. If blank, the value in the config will be used instead.")
	flagOpenBrowser = flag.Bool("openbrowser", true, "Launches the UI on startup")
	flagReindex     = flag.Bool("reindex", false, "Reindex all blobs on startup")
	flagRecovery    = flag.Bool("recovery", false, "Recovery mode: rebuild the blobpacked meta index. The tasks performed by the recovery mode might change in the future.")
	flagPollParent  bool
)

// For logging on Google Cloud Logging when not running on Google Compute Engine
// (for debugging).
var (
	flagGCEProjectID string
	flagGCELogName   string
	flagGCEJWTFile   string
)

func init() {
	if debug, _ := strconv.ParseBool(os.Getenv("CAMLI_DEBUG")); debug {
		flag.BoolVar(&flagPollParent, "pollparent", false, "Camlistored regularly polls its parent process to detect if it has been orphaned, and terminates in that case. Mainly useful for tests.")
		flag.StringVar(&flagGCEProjectID, "gce_project_id", "", "GCE project ID; required by --gce_log_name.")
		flag.StringVar(&flagGCELogName, "gce_log_name", "", "log all messages to that log name on Google Cloud Logging as well.")
		flag.StringVar(&flagGCEJWTFile, "gce_jwt_file", "", "Filename to the GCE Service Account's JWT (JSON) config file; required by --gce_log_name.")
	}
}

func exitf(pattern string, args ...interface{}) {
	if !strings.HasSuffix(pattern, "\n") {
		pattern = pattern + "\n"
	}
	fmt.Fprintf(os.Stderr, pattern, args...)
	osExit(1)
}

func slurpURL(urls string, limit int64) ([]byte, error) {
	res, err := http.Get(urls)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	return ioutil.ReadAll(io.LimitReader(res.Body, limit))
}

// loadConfig returns the server's parsed config file, locating it using the provided arg.
//
// The arg may be of the form:
// - empty, to mean automatic (will write a default high-level config if
//   no cloud config is available)
// - a filepath absolute or relative to the user's configuration directory,
// - a URL
func loadConfig(arg string) (conf *serverinit.Config, isNewConfig bool, err error) {
	if strings.HasPrefix(arg, "http://") || strings.HasPrefix(arg, "https://") {
		contents, err := slurpURL(arg, 256<<10)
		if err != nil {
			return nil, false, err
		}
		conf, err = serverinit.Load(contents)
		return conf, false, err
	}
	var absPath string
	switch {
	case arg == "":
		absPath = osutil.UserServerConfigPath()
		_, err = wkfs.Stat(absPath)
		if err != nil {
			if !os.IsNotExist(err) {
				return
			}
			conf, err = serverinit.DefaultEnvConfig()
			if err != nil || conf != nil {
				return
			}
			err = wkfs.MkdirAll(osutil.CamliConfigDir(), 0700)
			if err != nil {
				return
			}
			log.Printf("Generating template config file %s", absPath)
			if err = serverinit.WriteDefaultConfigFile(absPath, sqlite.CompiledIn()); err == nil {
				isNewConfig = true
			}
		}
	case filepath.IsAbs(arg):
		absPath = arg
	default:
		absPath = filepath.Join(osutil.CamliConfigDir(), arg)
	}
	conf, err = serverinit.LoadFile(absPath)
	return
}

// If cert/key files are specified, and found, use them.
// If cert/key files are specified, not found, and the default values, generate
// them (self-signed CA used as a cert), and use them.
// If cert/key files are not specified, use Let's Encrypt.
func setupTLS(ws *webserver.Server, config *serverinit.Config, hostname string) {
	cert, key := config.OptionalString("httpsCert", ""), config.OptionalString("httpsKey", "")
	if !config.OptionalBool("https", true) {
		return
	}
	if (cert != "") != (key != "") {
		exitf("httpsCert and httpsKey must both be either present or absent")
	}

	defCert := osutil.DefaultTLSCert()
	defKey := osutil.DefaultTLSKey()
	const hint = "You must add this certificate's fingerprint to your client's trusted certs list to use it. Like so:\n\"trustedCerts\": [\"%s\"],"
	if cert == defCert && key == defKey {
		_, err1 := wkfs.Stat(cert)
		_, err2 := wkfs.Stat(key)
		if err1 != nil || err2 != nil {
			if os.IsNotExist(err1) || os.IsNotExist(err2) {
				sig, err := httputil.GenSelfTLSFiles(hostname, defCert, defKey)
				if err != nil {
					exitf("Could not generate self-signed TLS cert: %q", err)
				}
				log.Printf(hint, sig)
			} else {
				exitf("Could not stat cert or key: %q, %q", err1, err2)
			}
		}
	}
	if cert == "" && key == "" {
		// Use Let's Encrypt if no files are specified, and we have a usable hostname.
		if netutil.IsFQDN(hostname) {
			m := autocert.Manager{
				Prompt:     autocert.AcceptTOS,
				HostPolicy: autocert.HostWhitelist(hostname),
				Cache:      autocert.DirCache(osutil.DefaultLetsEncryptCache()),
			}
			log.Print("TLS enabled, with Let's Encrypt")
			ws.SetTLS(webserver.TLSSetup{
				CertManager: m.GetCertificate,
			})
			return
		}
		// Otherwise generate new certificates
		sig, err := httputil.GenSelfTLSFiles(hostname, defCert, defKey)
		if err != nil {
			exitf("Could not generate self signed creds: %q", err)
		}
		log.Printf(hint, sig)
		cert = defCert
		key = defKey
	}
	data, err := wkfs.ReadFile(cert)
	if err != nil {
		exitf("Failed to read pem certificate: %s", err)
	}
	sig, err := httputil.CertFingerprint(data)
	if err != nil {
		exitf("certificate error: %v", err)
	}
	log.Printf("TLS enabled, with SHA-256 certificate fingerprint: %v", sig)
	ws.SetTLS(webserver.TLSSetup{
		CertFile: cert,
		KeyFile:  key,
	})
}

var osExit = os.Exit // testing hook

func handleSignals(shutdownc <-chan io.Closer) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP)
	signal.Notify(c, syscall.SIGINT)
	for {
		sig := <-c
		sysSig, ok := sig.(syscall.Signal)
		if !ok {
			log.Fatal("Not a unix signal")
		}
		switch sysSig {
		case syscall.SIGHUP:
			log.Print("SIGHUP: restarting camli")
			err := osutil.RestartProcess()
			if err != nil {
				log.Fatal("Failed to restart: " + err.Error())
			}
		case syscall.SIGINT:
			log.Print("Got SIGINT: shutting down")
			donec := make(chan bool)
			go func() {
				cl := <-shutdownc
				if err := cl.Close(); err != nil {
					exitf("Error shutting down: %v", err)
				}
				donec <- true
			}()
			select {
			case <-donec:
				log.Printf("Shut down.")
				osExit(0)
			case <-time.After(2 * time.Second):
				exitf("Timeout shutting down. Exiting uncleanly.")
			}
		default:
			log.Fatal("Received another signal, should not happen.")
		}
	}
}

// listenAndBaseURL finds the configured, default, or inferred listen address
// and base URL from the command-line flags and provided config.
func listenAndBaseURL(config *serverinit.Config) (listen, baseURL string) {
	baseURL = config.OptionalString("baseURL", "")
	listen = *flagListen
	listenConfig := config.OptionalString("listen", "")
	// command-line takes priority over config
	if listen == "" {
		listen = listenConfig
		if listen == "" {
			exitf("\"listen\" needs to be specified either in the config or on the command line")
		}
	}
	return
}

func redirectFromHTTP(base string) {
	http.ListenAndServe(":80", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, base, http.StatusFound)
	}))
}

// certHostname figures out the name to use for the TLS certificates, using baseURL
// and falling back to the listen address if baseURL is empty or invalid.
func certHostname(listen, baseURL string) (string, error) {
	hostPort, err := netutil.HostPort(baseURL)
	if err != nil {
		hostPort = listen
	}
	hostname, _, err := net.SplitHostPort(hostPort)
	if err != nil {
		return "", fmt.Errorf("failed to find hostname for cert from address %q: %v", hostPort, err)
	}
	return hostname, nil
}

func maybeSetupGoogleCloudLogging() {
	if flagGCEProjectID == "" && flagGCELogName == "" && flagGCEJWTFile == "" {
		return
	}
	if flagGCEProjectID == "" || flagGCELogName == "" || flagGCEJWTFile == "" {
		exitf("All of --gce_project_id, --gce_log_name, and --gce_jwt_file must be specified for logging on Google Cloud Logging.")
	}
	jsonSlurp, err := ioutil.ReadFile(flagGCEJWTFile)
	if err != nil {
		exitf("Error reading --gce_jwt_file value: %v", err)
	}
	jwtConf, err := google.JWTConfigFromJSON(jsonSlurp, logging.Scope)
	if err != nil {
		exitf("Error reading --gce_jwt_file value: %v", err)
	}
	logc, err := logging.NewClient(context.Background(), flagGCEProjectID, flagGCELogName, option.WithHTTPClient(jwtConf.Client(context.Background())))
	if err != nil {
		exitf("Error creating GCL client: %v", err)
	}
	log.SetOutput(io.MultiWriter(os.Stderr, logc.Writer(logging.Debug)))
}

// main wraps Main so tests (which generate their own func main) can still run Main.
func main() {
	Main(nil, nil)
}

// Main sends on up when it's running, and shuts down when it receives from down.
func Main(up chan<- struct{}, down <-chan struct{}) {
	flag.Parse()

	if *flagVersion {
		fmt.Fprintf(os.Stderr, "camlistored version: %s\nGo version: %s (%s/%s)\n",
			buildinfo.Version(), runtime.Version(), runtime.GOOS, runtime.GOARCH)
		return
	}
	if *flagHelp {
		flag.Usage()
		os.Exit(0)
	}
	if *flagLegal {
		for _, l := range legal.Licenses() {
			fmt.Fprintln(os.Stderr, l)
		}
		return
	}
	if *flagRecovery {
		blobpacked.SetRecovery()
	}
	if env.OnGCE() {
		log.SetOutput(gce.LogWriter())
	} else {
		maybeSetupGoogleCloudLogging()
	}

	if *flagReindex {
		index.SetImpendingReindex()
	}

	log.Printf("Starting camlistored version %s; Go %s (%s/%s)", buildinfo.Version(), runtime.Version(),
		runtime.GOOS, runtime.GOARCH)

	shutdownc := make(chan io.Closer, 1) // receives io.Closer to cleanly shut down
	go handleSignals(shutdownc)

	// In case we're running in a Docker container with no
	// filesytem from which to load the root CAs, this
	// conditionally installs a static set if necessary. We do
	// this before we load the config file, which might come from
	// an https URL.
	httputil.InstallCerts()

	config, isNewConfig, err := loadConfig(*flagConfigFile)
	if err != nil {
		exitf("Error loading config file: %v", err)
	}

	ws := webserver.New()
	listen, baseURL := listenAndBaseURL(config)

	hostname, err := certHostname(listen, baseURL)
	if err != nil {
		exitf("Bad baseURL or listen address: %v", err)
	}
	setupTLS(ws, config, hostname)

	err = ws.Listen(listen)
	if err != nil {
		exitf("Listen: %v", err)
	}

	if baseURL == "" {
		baseURL = ws.ListenURL()
	}

	shutdownCloser, err := config.InstallHandlers(ws, baseURL, *flagReindex, nil)
	if err != nil {
		exitf("Error parsing config: %v", err)
	}
	shutdownc <- shutdownCloser

	urlToOpen := baseURL
	if !isNewConfig {
		// user may like to configure the server at the initial startup,
		// open UI if this is not the first run with a new config file.
		urlToOpen += config.UIPath
	}

	if *flagOpenBrowser {
		go osutil.OpenURL(urlToOpen)
	}

	go ws.Serve()
	if flagPollParent {
		osutil.DieOnParentDeath()
	}

	if err := config.UploadPublicKey(); err != nil {
		exitf("Error uploading public key on startup: %v", err)
	}

	if err := config.StartApps(); err != nil {
		exitf("StartApps: %v", err)
	}

	for appName, appURL := range config.AppURL() {
		addr, err := netutil.HostPort(appURL)
		if err != nil {
			log.Printf("Could not get app %v address: %v", appName, err)
			continue
		}
		if err := netutil.AwaitReachable(addr, 5*time.Second); err != nil {
			log.Printf("Could not reach app %v: %v", appName, err)
		}
	}
	log.Printf("Available on %s", urlToOpen)

	if env.OnGCE() && strings.HasPrefix(baseURL, "https://") {
		go redirectFromHTTP(baseURL)
	}

	// Block forever, except during tests.
	up <- struct{}{}
	<-down
	osExit(0)
}
