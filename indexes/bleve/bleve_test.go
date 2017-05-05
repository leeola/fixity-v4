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

	Convey("Scenario: Basic querying", t, func() {
		k := newKala(tmp)
		defer os.RemoveAll(tmp)
		Convey("Given a single entry", func() {
			createdHashes, err := k.Write(
				kala.Commit{
					JsonMeta: &kala.JsonMeta{
						IndexedFields: kala.Fields{{
							Field: "field",
							Value: "foo bar baz",
						}},
					},
				},
				kala.Json{Json: []byte("{}")},
				nil,
			)
			So(err, ShouldBeNil)
			So(createdHashes, ShouldHaveLength, 2)
			createdVersionHash := createdHashes[1]
			Convey("When the correct value is queried", func() {
				r, err := k.Search(q.New().Const(q.Eq("field", "foo bar baz")))
				So(err, ShouldBeNil)
				Convey("Then it should show up in search results", func() {
					So(r, ShouldHaveLength, 1)
					So(r[0], ShouldEqual, createdVersionHash)
				})
			})
			Convey("When the incorrect value is queried", func() {
				r, err := k.Search(q.New().Const(q.Eq("field", "incorrect")))
				So(err, ShouldBeNil)
				Convey("Then it should not show up in search results", func() {
					So(r, ShouldHaveLength, 0)
				})
			})
			Convey("When the a substring of the full value is queried", func() {
				r, err := k.Search(q.New().Const(q.Eq("field", "bar")))
				So(err, ShouldBeNil)
				Convey("Then it should not show up in search results", func() {
					So(r, ShouldHaveLength, 0)
				})
			})
		})
	})

	Convey("Scenario: Multi field querying", t, func() {
		k := newKala(tmp)
		defer os.RemoveAll(tmp)
		Convey("Given multiple entries", func() {
			createdHashes, err := k.Write(
				kala.Commit{
					JsonMeta: &kala.JsonMeta{
						IndexedFields: kala.Fields{
							{
								Field: "fielda",
								Value: "foo",
							},
							{
								Field: "fieldb",
								Value: "bar",
							},
						},
					},
				},
				kala.Json{Json: []byte("{}")},
				nil,
			)
			So(err, ShouldBeNil)
			So(createdHashes, ShouldHaveLength, 2)
			createdVersionHash := createdHashes[1]

			Convey("When we query fielda with the correct value", func() {
				r, err := k.Search(q.New().Const(q.Eq("fielda", "foo")))
				So(err, ShouldBeNil)
				Convey("Then it should show up in search results", func() {
					So(r, ShouldHaveLength, 1)
					So(r[0], ShouldEqual, createdVersionHash)
				})
			})
			Convey("When we query fielda with the incorrect value", func() {
				r, err := k.Search(q.New().Const(q.Eq("fielda", "bar")))
				So(err, ShouldBeNil)
				Convey("Then it should not show up in search results", func() {
					So(r, ShouldHaveLength, 0)
				})
			})
			Convey("When we query fieldb with the correct value", func() {
				r, err := k.Search(q.New().Const(q.Eq("fieldb", "bar")))
				So(err, ShouldBeNil)
				Convey("Then it should show up in search results", func() {
					So(r, ShouldHaveLength, 1)
					So(r[0], ShouldEqual, createdVersionHash)
				})
			})
		})
	})

	Convey("Scenario: querying fulltextsearch", t, func() {
		k := newKala(tmp)
		defer os.RemoveAll(tmp)
		Convey("Given multiple entries", func() {
			createdHashes, err := k.Write(
				kala.Commit{
					JsonMeta: &kala.JsonMeta{
						// fts is default with bleve.
						IndexedFields: kala.Fields{
							{
								Field: "field",
								Value: "this is a fts field, with foo in it",
							},
						},
					},
				},
				kala.Json{Json: []byte("{}")},
				nil,
			)
			So(err, ShouldBeNil)
			So(createdHashes, ShouldHaveLength, 2)
			createdVersionHash := createdHashes[1]

			Convey("When we query fielda with the correct value", func() {
				r, err := k.Search(q.New().Const(q.Eq("field", "foo")))
				So(err, ShouldBeNil)
				Convey("Then it should show up in search results", func() {
					So(r, ShouldHaveLength, 1)
					So(r[0], ShouldEqual, createdVersionHash)
				})
			})
			Convey("When we query with an incorrect value", func() {
				r, err := k.Search(q.New().Const(q.Eq("field", "bar")))
				So(err, ShouldBeNil)
				Convey("Then it should not show up in search results", func() {
					So(r, ShouldHaveLength, 0)
				})
			})
		})
	})

	// Note that this test is attempting to ignore sort order. Eg, this only tests
	// that the skipping is consistent and works, not what the order of the fields
	// are.
	Convey("Scenario: Result skipping", t, func() {
		k := newKala(tmp)
		defer os.RemoveAll(tmp)
		Convey("Given 5 entries", func() {
			for i := 0; i < 5; i++ {
				_, err := k.Write(
					kala.Commit{
						JsonMeta: &kala.JsonMeta{
							IndexedFields: kala.Fields{{
								Field: "field",
								Value: "foo",
							}},
						},
					},
					kala.Json{Json: []byte("{}")},
					nil,
				)
				So(err, ShouldBeNil)
			}

			Convey("When we query the first two repeatedly", func() {
				query := q.New().Const(q.Eq("field", "foo")).Limit(2)
				a, err := k.Search(query)
				So(err, ShouldBeNil)
				b, err := k.Search(query)
				So(err, ShouldBeNil)
				Convey("Then it should return with the same results both times", func() {
					So(a, ShouldHaveLength, 2)
					So(b, ShouldHaveLength, 2)
					for i, h := range a {
						So(h, ShouldEqual, b[i])
					}
				})
			})

			Convey("When we query the second two repeatedly", func() {
				firstTwoQ := q.New().Const(q.Eq("field", "foo")).Limit(2)
				firstTwo, err := k.Search(firstTwoQ)
				So(err, ShouldBeNil)
				secondTwoQ := q.New().Const(q.Eq("field", "foo")).Limit(2).Skip(2)
				a, err := k.Search(secondTwoQ)
				So(err, ShouldBeNil)
				b, err := k.Search(secondTwoQ)
				So(err, ShouldBeNil)
				Convey("Then it should return with the same results both times", func() {
					So(a, ShouldHaveLength, 2)
					So(b, ShouldHaveLength, 2)
					for i, h := range a {
						So(h, ShouldEqual, b[i])
					}
				})
				Convey("Then it should not return with the first two", func() {
					for _, fh := range firstTwo {
						for _, ah := range a {
							So(fh, ShouldNotEqual, ah)
						}
					}
				})
			})
		})
	})
}
