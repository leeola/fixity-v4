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
	"github.com/leeola/kala/node"
)

func (c *Client) Upload(r io.Reader, mc contenttype.MetaChanges) ([]string, error) {
	u, err := url.Parse(c.kalaAddr)
	if err != nil {
		return nil, err
	}
	u.Path = path.Join(u.Path, "upload")

	q := u.Query()
	for k, v := range mc {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()

	res, err := c.httpClient.Post(u.String(), "plain/text", r)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, errors.Errorf("unexpected kala response: %d %q",
			res.StatusCode, res.Status)
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var hashesRes node.HashesResponse
	if err := json.Unmarshal(b, &hashesRes); err != nil {
		return nil, err
	}

	return hashesRes.Hashes, nil
}

func (c *Client) UploadMeta(mc contenttype.MetaChanges) ([]string, error) {
	u, err := url.Parse(c.kalaAddr)
	if err != nil {
		return nil, err
	}
	u.Path = path.Join(u.Path, "upload", "meta")

	q := u.Query()
	for k, v := range mc {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()

	res, err := c.httpClient.Post(u.String(), "plain/text", nil)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, errors.Errorf("unexpected kala response: %d %q",
			res.StatusCode, res.Status)
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var hashesRes node.HashesResponse
	if err := json.Unmarshal(b, &hashesRes); err != nil {
		return nil, err
	}

	return hashesRes.Hashes, nil
}
