package snail

import (
	"os"
	"testing"

	"github.com/leeola/fixity"
	"github.com/leeola/fixity/q"
	"github.com/leeola/fixity/util/testutil"
	. "github.com/smartystreets/goconvey/convey"
)

func newSnail() (s *Snail, tmpDir string) {
	tmp := testutil.MustTempDir("snail")

	c := Config{
		Path:   tmp,
		Bucket: "testing",
	}

	s, err := New(c)
	if err != nil {
		panic(err)
	}

	return s, tmp
}

func TestIndexing(t *testing.T) {
	Convey("Scenario: Indexing", t, func() {
		s, tmp := newSnail()
		defer os.RemoveAll(tmp)
		defer s.Close()

		Convey("Given an empty index", func() {
			Convey("When multiple repeat data is indexed", func() {
				fields := fixity.Fields{{Field: "field1", Value: "value1"}}
				var errs []error
				errs = append(errs, s.Index("foo", "bar", fields))
				errs = append(errs, s.Index("foo", "bar", fields))
				errs = append(errs, s.Index("foo", "bar", fields))
				Convey("Then it should not error", func() {
					for _, err := range errs {
						So(err, ShouldBeNil)
					}
				})
			})

			Convey("When multiple unique data is indexed", func() {
				var errs []error
				errs = append(errs, s.Index("foo", "bar",
					fixity.Fields{{Field: "field1", Value: "value1"}}))
				errs = append(errs, s.Index("foo", "bar",
					fixity.Fields{{Field: "field2", Value: "value2"}}))
				errs = append(errs, s.Index("foo", "bar",
					fixity.Fields{{Field: "field3", Value: "value3"}}))
				Convey("Then it should not error", func() {
					for _, err := range errs {
						So(err, ShouldBeNil)
					}
				})
			})
		})
	})
}

func TestSearch(t *testing.T) {
	Convey("Scenario: Searching", t, func() {
		s, tmp := newSnail()
		defer os.RemoveAll(tmp)
		defer s.Close()

		Convey("Given an empty index", func() {
			Convey("When queried", func() {
				keys, err := s.Search(q.New().Const(q.Eq("foo", "bar")))
				So(err, ShouldBeNil)
				Convey("Then it should not match anything", func() {
					So(keys, ShouldHaveLength, 0)
				})
			})
		})

		Convey("Given one document", func() {
			err := s.Index("hash", "id", fixity.Fields{{Field: "field1", Value: "value1"}})
			So(err, ShouldBeNil)
			Convey("When queried with a matching value", func() {
				keys, err := s.Search(q.New().Const(q.Eq("field1", "value1")))
				So(err, ShouldBeNil)
				Convey("Then it should respond with the expected data", func() {
					So(keys, ShouldHaveLength, 1)
					So(keys[0], ShouldEqual, "hash")
				})
			})
		})
	})
}
