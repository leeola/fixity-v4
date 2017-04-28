package bleve

import (
	"os"
	"testing"

	"github.com/leeola/kala"
	"github.com/leeola/kala/impl/local"
	"github.com/leeola/kala/q"
	"github.com/leeola/kala/util/testutil"
	. "github.com/smartystreets/goconvey/convey"
)

func newKala(rootDir string) kala.Kala {
	iConf := Config{
		Path: rootDir,
	}
	i, err := New(iConf)
	if err != nil {
		panic(err)
	}

	kConf := local.Config{
		Store: &testutil.NoopStore{},
		Index: i,
	}
	k, err := local.New(kConf)
	if err != nil {
		panic(err)
	}

	return k
}

func TestBleve(t *testing.T) {
	tmp := testutil.MustTempDir("kala-bleve")
	k := newKala(tmp)
	defer os.RemoveAll(tmp)

	Convey("Scenario: Basic querying", t, func() {
		Convey("Given no other entries", func() {
			Convey("When we create a new entry", func() {
				createdHashes, err := k.Write(
					kala.Commit{},
					kala.Json{Meta: kala.JsonMeta{
						IndexedFields: kala.Fields{{
							Field: "field",
							Value: "foo",
						}},
					}},
					nil,
				)
				So(err, ShouldBeNil)
				So(createdHashes, ShouldHaveLength, 2)
				createdVersionHash := createdHashes[1]
				Convey("Then it should show up in search results", func() {
					r, err := k.Search(q.New().Const(q.Eq("field", "foo")))
					So(err, ShouldBeNil)
					So(r, ShouldHaveLength, 1)
					So(r[0], ShouldEqual, createdVersionHash)
				})
			})
		})
	})
}
