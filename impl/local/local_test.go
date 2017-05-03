package local

import (
	"testing"

	"github.com/fatih/structs"
	"github.com/leeola/kala"
	"github.com/leeola/kala/util/kalautil"
	. "github.com/smartystreets/goconvey/convey"
)

func TestMakeFields(t *testing.T) {
	k := &Local{}

	Convey("Scenario: Unmarshalling a field", t, func() {
		Convey("Given field to be unmarshalled", func() {
			Convey("When makeFields() is called", func() {
				foo := struct {
					Foo string `json:"foo"`
				}{Foo: "foo value"}
				fields, err := k.makeFields(
					kala.Version{
						JsonMeta: &kala.JsonMeta{
							IndexedFields: kala.Fields{{
								Field: "foo",
							}},
						},
					},
					kala.Json{
						Json: kalautil.MustMarshalJson(&foo).Json,
					},
				)
				So(err, ShouldBeNil)

				var fooField kala.Field
				for _, f := range fields {
					if f.Field == "foo" {
						fooField = f
						break
					}
				}

				Convey("Then it should create the foo indexField", func() {
					So(structs.IsZero(fooField), ShouldBeFalse)
				})
				Convey("Then the indexField should have the proper value", func() {
					v, _ := fooField.Value.(string)
					So(v, ShouldEqual, foo.Foo)
				})
			})
		})
	})
}
