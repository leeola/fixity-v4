package clijson

import (
	"fmt"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCliJson(t *testing.T) {
	Convey("Scenario: parsing input", t, func() {
		type test struct {
			Input    []string
			Expected string
		}

		tests := []test{
			{
				[]string{"foo=bar"},
				`{"foo":"bar"}`,
			},
			{
				[]string{"foo=bar", "1"},
				`[{"foo":"bar"},1]`,
			},
			{
				[]string{"[", "foo=bar", "]"},
				`[{"foo":"bar"}]`,
			},
			{
				[]string{"{", "foo=bar", "baz=bat"},
				`{"baz":"bat","foo":"bar"}`,
			},
			{
				[]string{"{", "foo=bar", "baz=bat", "}"},
				`{"baz":"bat","foo":"bar"}`,
			},
			{
				[]string{"{", "foo=bar", "baz=[", "1", "str", "]", "}"},
				`{"baz":[1,"str"],"foo":"bar"}`,
			},
			{
				[]string{"{", "foo=bar", "baz=[", "1", "str"},
				`{"baz":[1,"str"],"foo":"bar"}`,
			},
			{
				[]string{"[", "foo=bar", "1", "bat=baz", "boo=[", "1", "str"},
				`[{"foo":"bar"},1,{"bat":"baz","boo":[1,"str"]}]`,
			},
		}

		for _, test := range tests {
			Convey(fmt.Sprintf("Given input: %s", strings.Join(test.Input, " ")), func() {
				Convey("When run against CliJson", func() {
					outB, err := CliJson(test.Input)
					So(err, ShouldBeNil)
					Convey("Then it should return the expected output", func() {
						So(string(outB), ShouldEqual, test.Expected)
					})
				})
			})
		}
	})
}
