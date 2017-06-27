package local

import (
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

func TestBlockchainAppendContent(t *testing.T) {
	Convey("Scenario: Appending content to blockchain", t, func() {
		bc := mustTestLocal().Blockchain()

		Convey("Given an empty blockchain", func() {
			Convey("When content is appended", func() {
				c := fixity.Content{Hash: "foo"}
				b, err := bc.AppendContent(c)
				So(err, ShouldBeNil)

				Convey("Then it should extend the blockchain", func() {
					So(b.Block, ShouldEqual, 1)
				})
				Convey("Then it should store the content hash", func() {
					So(b.ContentBlock, ShouldNotBeNil)
					So(b.ContentBlock.Hash, ShouldEqual, c.Hash)
				})
			})
		})

		Convey("Given an non-empty blockchain", func() {
			c := fixity.Content{Hash: "foo"}
			_, err := bc.AppendContent(c)
			So(err, ShouldBeNil)
			Convey("When content is appended", func() {
				c := fixity.Content{Hash: "bar"}
				b, err := bc.AppendContent(c)
				So(err, ShouldBeNil)

				Convey("Then it should extend the blockchain", func() {
					So(b.Block, ShouldEqual, 2)
				})
				Convey("Then it should store the content hash", func() {
					So(b.ContentBlock, ShouldNotBeNil)
					So(b.ContentBlock.Hash, ShouldEqual, c.Hash)
				})
			})
		})
	})
}
