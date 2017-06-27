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
	Convey("Scenario: Appending content to the blockchain", t, func() {
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

				Convey("Then it update the head", func() {
					head, err := bc.Head()
					So(err, ShouldBeNil)
					So(b, ShouldResemble, head)
				})
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

func TestBlockchainDeleteContent(t *testing.T) {
	Convey("Scenario: Deleting content from the blockchain", t, func() {
		bc := mustTestLocal().Blockchain()
		Convey("Given a single content to be deleted", func() {
			c := fixity.Content{Hash: "foo"}
			b, err := bc.AppendContent(c)
			So(err, ShouldBeNil)
			bHash := b.Hash
			b, err = bc.AppendContent(fixity.Content{Hash: "bar"})
			So(err, ShouldBeNil)
			Convey("When the content is deleted", func() {
				b, err := bc.DeleteContent(c)
				So(err, ShouldBeNil)

				Convey("Then it update the head", func() {
					head, err := bc.Head()
					So(err, ShouldBeNil)
					So(b, ShouldResemble, head)
				})
				Convey("Then it should extend the blockchain", func() {
					So(b.Block, ShouldEqual, 3)
				})
				Convey("It should mark the content hash for deletion", func() {
					So(b.DeleteBlock, ShouldNotBeNil)
					So(bHash, ShouldBeIn, b.DeleteBlock.Hashes)
				})
				Convey("It should only mark the content hash for deletion", func() {
					So(b.DeleteBlock, ShouldNotBeNil)
					So(b.DeleteBlock.Hashes, ShouldHaveLength, 1)
				})
			})
		})
		Convey("Given multiple contents to be deleted", func() {
			c1 := fixity.Content{Hash: "foo"}
			b, err := bc.AppendContent(c1)
			So(err, ShouldBeNil)
			b1Hash := b.Hash
			_, err = bc.AppendContent(fixity.Content{Hash: c1.Hash + " ignore"})
			So(err, ShouldBeNil)
			c2 := fixity.Content{Hash: "bar"}
			_, err = bc.AppendContent(c2)
			So(err, ShouldBeNil)
			b2Hash := b.Hash
			_, err = bc.AppendContent(fixity.Content{Hash: c1.Hash + " ignore"})
			So(err, ShouldBeNil)
			Convey("When the content is deleted", func() {
				b, err := bc.DeleteContent(c1, c2)
				So(err, ShouldBeNil)

				Convey("Then it update the head", func() {
					head, err := bc.Head()
					So(err, ShouldBeNil)
					So(b, ShouldResemble, head)
				})
				Convey("Then it should extend the blockchain", func() {
					So(b.Block, ShouldEqual, 5)
				})
				Convey("It should mark the content hash for deletion", func() {
					So(b.DeleteBlock, ShouldNotBeNil)
					So(b1Hash, ShouldBeIn, b.DeleteBlock.Hashes)
					So(b2Hash, ShouldBeIn, b.DeleteBlock.Hashes)
				})
				Convey("It should only mark the content hash for deletion", func() {
					So(b.DeleteBlock, ShouldNotBeNil)
					So(b.DeleteBlock.Hashes, ShouldHaveLength, 2)
				})
			})
		})
	})
}
