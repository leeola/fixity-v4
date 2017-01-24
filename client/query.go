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

type HashResult struct {
	index.Result
	Error error
}

func (c *Client) QueryChan(cancel chan struct{}, q index.Query, ss ...index.SortBy) <-chan HashResult {
	// maxResults is the max number of results to return. Ie, the expected behavior of
	// q.Limit.
	maxResults := q.Limit

	// Paginate the results by X.
	q.Limit = 5
	// If the page size is larger than the max desired results, set our page size
	// to the max results.
	if maxResults > 0 && q.Limit > maxResults {
		q.Limit = maxResults
	}

	// TODO(leeola): only add the indexEntry sort if it doesn't already exist
	ss = append(ss, index.SortBy{
		Field: "indexEntry",
	})

	// Setting the channel buffer to the page size for no specific reason,
	// it just seems reasonable.
	ch := make(chan HashResult, q.Limit)

	go func() {
		totalResults := 0
		for {
			select {
			case <-cancel:
				close(ch)
				return
			default:
			}

			results, err := c.Query(q, ss...)
			if err != nil {
				ch <- HashResult{Error: err}
				close(ch)
				return
			}

			if len(results.Hashes) == 0 {
				close(ch)
				return
			}

			for _, h := range results.Hashes {
				select {
				case <-cancel:
					close(ch)
					return
				default:
				}

				q.FromEntry = h.Entry + 1
				ch <- HashResult{Result: index.Result{
					IndexVersion: results.IndexVersion,
					Hash:         h,
				}}

				totalResults++
				if maxResults > 0 && totalResults >= maxResults {
					close(ch)
					return
				}
			}
		}
	}()

	return ch
}

func (c *Client) Query(q index.Query, ss ...index.SortBy) (index.Results, error) {
	u, err := url.Parse(c.kalaAddr)
	if err != nil {
		return index.Results{}, err
	}
	u.Path = path.Join("index", "query")

	v := u.Query()
	if q.Limit != 0 {
		v.Add("limit", strconv.Itoa(q.Limit))
	}
	if q.SearchVersions {
		v.Add("searchVersions", strconv.FormatBool(q.SearchVersions))
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
