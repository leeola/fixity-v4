package client

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strconv"

	"github.com/leeola/errors"
	"github.com/leeola/kala/index"
)

func (c *Client) Query(q index.Query) (index.Results, error) {
	u, err := url.Parse(c.kalaAddr)
	if err != nil {
		return index.Results{}, err
	}
	u.Path = path.Join("index", "query")

	v := u.Query()
	if q.Limit != 0 {
		v.Add("limit", strconv.Itoa(q.Limit))
	}
	if q.FromEntry != 0 {
		v.Add("fromEntry", strconv.Itoa(q.FromEntry))
	}
	if q.IndexVersion != "" {
		v.Add("indexVersion", q.IndexVersion)
	}
	u.RawQuery = v.Encode()

	res, err := c.httpClient.Get(u.String())
	if err != nil {
		return index.Results{}, err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return index.Results{}, index.ErrNoQueryResults
	}

	if res.StatusCode != http.StatusOK {
		return index.Results{}, errors.Errorf("unexpected kala response: %d %q",
			res.Status, res.StatusCode)
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return index.Results{}, errors.Wrap(err, "failed to read request body")
	}

	var results index.Results
	if err := json.Unmarshal(b, &results); err != nil {
		return index.Results{}, errors.Wrap(err, "failed to unmarshal request body")
	}

	return results, nil
}
