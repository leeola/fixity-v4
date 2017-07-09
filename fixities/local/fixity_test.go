package local

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/leeola/fixity"
	"github.com/leeola/fixity/indexes/nopindex"
	"github.com/leeola/fixity/stores/memory"
	. "github.com/smartystreets/goconvey/convey"
)

func mustTestLocal() fixity.Fixity {
	c := Config{
		Index: nopindex.New(),
		Store: memory.New(),
		Db:    newMemoryDb(),
	}
	f, err := New(c)
	if err != nil {
		panic(err)
	}
	return f
}

func TestFixityIgnoreDuplicateBlob(t *testing.T) {
	Convey("Scenario: Writing with IgnoreDuplicateBlob", t, func() {
		f := mustTestLocal()
		req := fixity.NewWrite("foo", ioutil.NopCloser(strings.NewReader("bar")))
		req.IgnoreDuplicateBlob = true

		Convey("Given the blob didn't exist previously", func() {
			Convey("When the written", func() {
				_, err := f.WriteRequest(req)

				Convey("Then it should not error", func() {
					So(err, ShouldBeNil)
				})
			})
		})
		Convey("Given the blob existed previously", func() {
			previousContent, err := f.WriteRequest(req)
			So(err, ShouldBeNil)

			Convey("When the written", func() {
				c, err := f.WriteRequest(req)

				Convey("Then it should not error", func() {
					So(err, ShouldBeNil)
				})
				Convey("Then the returned content should be the previous Content", func() {
					So(c.Hash, ShouldEqual, previousContent.Hash)
				})
			})
		})
	})
}
