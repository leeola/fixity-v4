package client

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"

	"github.com/leeola/errors"
	"github.com/leeola/kala/contenttype"
	"github.com/leeola/kala/node/handlers"
)

func (c *Client) GetDownloadMetaExport(h string) (contenttype.Changes, error) {
	u, err := url.Parse(c.kalaAddr)
	if err != nil {
		return nil, err
	}
	u.Path = path.Join(u.Path, "download", h, "meta", "export")

	res, err := c.httpClient.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var cRes handlers.ChangesResponse
	if err := json.Unmarshal(b, &cRes); err != nil {
		return nil, err
	}

	if cRes.Error != "" {
		return nil, errors.Errorf("error: %s", cRes.Error)
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.Errorf("unexpected kala response: %d %q",
			res.StatusCode, res.Status)
	}

	return cRes.Changes, nil
}

func (c *Client) GetDownloadBlob(h string) (io.ReadCloser, error) {
	u, err := url.Parse(c.kalaAddr)
	if err != nil {
		return nil, err
	}
	u.Path = path.Join(u.Path, "download", h, "blob")

	res, err := c.httpClient.Get(u.String())
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		if res.Body != nil {
			// We're going to read the body due to the bad status code, close the
			// body when we're done.
			defer res.Body.Close()

			b, err := ioutil.ReadAll(res.Body)
			if err != nil {
				return nil, err
			}

			var cRes handlers.ErrorResponse
			if err := json.Unmarshal(b, &cRes); err != nil {
				return nil, err
			}

			if cRes.Error != "" {
				return nil, errors.Errorf("kalanode: %s", cRes.Error)
			}
		}

		// if we didn't have a good errorresponse, return an unexpected response err
		return nil, errors.Errorf("unexpected kala response: %d %q",
			res.StatusCode, res.Status)
	}

	return res.Body, nil
}
