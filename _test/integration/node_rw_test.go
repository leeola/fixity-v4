package integration

import (
	"io/ioutil"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/leeola/kala/client"
)

// TODO: Use testutil tempdir instead of hardcoded tempdir,
// or at least in combination with the _test/tmp dir.
const TestTmpDir = "../tmp"

func TestReadWrite(t *testing.T) {
	os.RemoveAll(TestTmpDir)

	for _, nt := range NodeTests(TestTmpDir) {
		t.Log(DescribeNodeTest(nt))

		ts := httptest.NewServer(nt.Node)
		defer ts.Close()

		client, err := client.New(client.Config{KalaAddr: ts.URL})
		if err != nil {
			t.Fatal(err)
		}

		source := "foo bar baz"
		h, err := client.Write([]byte(source))
		if err != nil {
			t.Fatal(err)
		}

		read, err := client.Read(h)
		if err != nil {
			t.Fatal(err)
		}

		readB, err := ioutil.ReadAll(read)
		if err != nil {
			t.Fatal(err)
		}

		if source != string(readB) {
			t.Fatalf(
				"Write and Read bytes do not equal. got:%q, want:%q",
				string(readB), source,
			)
		}

		h, err = client.PostBlob(strings.NewReader(source))
		if err != nil {
			t.Fatal(err)
		}

		read, err = client.GetBlob(h)
		if err != nil {
			t.Fatal(err)
		}

		readB, err = ioutil.ReadAll(read)
		if err != nil {
			t.Fatal(err)
		}
		if source != string(readB) {
			t.Fatalf(
				"PostBlob and GetBlob bytes do not equal. got:%q, want:%q",
				string(readB), source,
			)
		}

		// TODO(leeola): Write a random bytes generator to test upload and download
		// repeatedly with varying sizes of data.
	}
}
