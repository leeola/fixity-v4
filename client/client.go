package client

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"

	"github.com/leeola/errors"
	"github.com/leeola/kala/store"
)

type Config struct {
	// The kala server address that this Client will be talking to.
	KalaAddr string

	// Optional. The http client that this Kala Client will use.
	HttpClient *http.Client
}

type Client struct {
	kalaAddr   *url.URL
	httpClient *http.Client
}

func New(c Config) (*Client, error) {
	if c.KalaAddr == "" {
		return nil, errors.New("missing Config field: KalaAddr")
	}

	u, err := url.Parse(c.KalaAddr)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse KalaAddr")
	}

	if c.HttpClient == nil {
		c.HttpClient = &http.Client{}
	}

	return &Client{
		kalaAddr:   u,
		httpClient: c.HttpClient,
	}, nil
}

func (c *Client) Exists(h string) (bool, error) {
	u := &*c.kalaAddr
	u.Path = path.Join("content", h)
	res, err := c.httpClient.Head(u.String())
	if err != nil {
		return false, err
	}

	return res.StatusCode == http.StatusOK, nil
}

func (c *Client) Read(h string) (io.ReadCloser, error) {
	u := &*c.kalaAddr
	u.Path = path.Join("content", h)
	res, err := c.httpClient.Get(u.String())
	if err != nil {
		return nil, err
	}

	if res.StatusCode == http.StatusNotFound {
		return nil, store.HashNotFoundErr
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.Errorf("unexpected kala response: %d %q",
			res.Status, res.StatusCode)
	}

	return res.Body, nil
}

func (c *Client) Write(b []byte) (string, error) {
	r := bytes.NewReader(b)
	u := &*c.kalaAddr
	u.Path = "content"
	// TODO(leeola): decide the best content type for the http api
	res, err := c.httpClient.Post(u.String(), "application/json", r)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", errors.Errorf("unexpected kala response: %d %q",
			res.Status, res.StatusCode)
	}

	resB, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", errors.Wrap(err, "failed to read body")
	}

	return string(resB), nil
}

func (c *Client) WriteHash(h string, b []byte) error {
	bLen := len(b)
	r := bytes.NewReader(b)

	u := &*c.kalaAddr
	u.Path = path.Join("content", h)

	req, err := http.NewRequest(http.MethodPut, u.String(), r)
	if err != nil {
		return errors.Wrap(err, "failed to construct request")
	}
	// TODO(leeola): decide the best content type for the http api
	req.Header["Content-Type"] = []string{"application/json"}
	req.ContentLength = int64(bLen)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to make request")
	}
	res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return errors.Errorf("unexpected kala response: %d %q",
			res.Status, res.StatusCode)
	}

	return nil
}
