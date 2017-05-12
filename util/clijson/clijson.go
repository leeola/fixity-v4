package clijson

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
)

// CliJson converts the given os.Args array into Json bytes.
//
//
// WARNING: This package was sort of hacked together while experimenting
// with the syntax. A proper implementation is needed.
func CliJson(args []string) ([]byte, error) {
	v := process(args)
	if v == nil {
		return nil, errors.New("empty args")
	}

	return json.Marshal(v)
}

func process(args []string) interface{} {
	if len(args) == 0 {
		return nil
	}

	var values []interface{}
	for len(args) > 0 {
		var v interface{}

		switch {
		case args[0] == "[":
			v, args = list(args[1:])

		case args[0] == "{":
			v, args = object(args[1:])

		case isKey(args[0]):
			v, args = object(args)

		default:
			v = value(args[0])
			args = args[1:]
		}

		values = append(values, v)
	}

	// remove the implicit list if there is only one value.
	if len(values) == 1 {
		return values[0]
	}

	return values
}

func list(args []string) ([]interface{}, []string) {
	// always return at least an empty list, never nil.
	values := []interface{}{}

	if len(args) == 0 {
		return values, nil
	}

	for len(args) > 0 {
		var v interface{}
		switch {
		case args[0] == "]":
			return values, args[1:]

		case args[0] == "[":
			v, args = list(args[1:])

		case args[0] == "{":
			v, args = object(args[1:])

		case isKey(args[0]):
			v, args = object(args)

		default:
			v = value(args[0])
			args = args[1:]
		}

		values = append(values, v)
	}

	return values, nil
}

func object(args []string) (map[string]interface{}, []string) {
	// always return at least an empty map, never nil.
	values := map[string]interface{}{}

	if len(args) == 0 {
		return values, nil
	}

	for len(args) > 0 {
		var (
			k string
			v interface{}
		)

		if args[0] == "}" {
			return values, args[1:]
		}

		// if it's not a key, end the creation of the object
		if !isKey(args[0]) {
			// don't crop the args
			return values, args
		}

		kv := strings.SplitN(args[0], "=", 2)

		k = kv[0]

		switch kv[1] {
		case "[":
			v, args = list(args[1:])
		case "{":
			v, args = object(args[1:])
		default:
			v = value(kv[1])
			args = args[1:]
		}

		values[k] = v
	}

	return values, nil
}

func value(s string) interface{} {
	if v, err := strconv.Atoi(s); err == nil {
		return int(v)
	}
	if v, err := strconv.ParseBool(s); err == nil {
		return v
	}
	return s
}

func isKey(s string) bool {
	return strings.Contains(s, "=")
}
