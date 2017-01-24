package strutil

import "strings"

// SplitQueryField splits a query field into a key value pair.
//
// The query fields are separated by the first colon.
func SplitQueryField(field string) (string, string) {
	var k, v string
	split := strings.SplitN(field, ":", 2)
	if len(split) == 2 {
		k, v = split[0], split[1]
	} else {
		v = split[0]
	}
	return k, v
}
