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

func (c *Client) Query(q index.Query, ss []index.SortBy) (index.Results, error) {
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
	if q.Metadata != nil {
		for k, kv := range q.Metadata {
			// TODO(leeola): implement meaning method for taking non-string values
			s, ok := kv.(string)
			if !ok {
				return index.Results{}, errors.Errorf(
					"unhandled non-string metadata query: %s=%s", k, kv)
			}
			v.Set(k, s)
		}
	}
	if len(ss) > 0 {
		for _, s := range ss {
			if s.Descending {
				v.Add("sortDescending", s.Field)
			} else {
				v.Add("sortAscending", s.Field)
			}
		}
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
