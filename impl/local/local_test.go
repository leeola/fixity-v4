package local

import (
	"testing"

	"github.com/fatih/structs"
	"github.com/leeola/fixity"
	"github.com/leeola/fixity/util/fixityutil"
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
					fixity.Version{
						MultiJsonHash: fixity.MultiJsonHash{
							"json type": fixityutil.JsonHashWithMetaFields(
								fixity.Fields{{
									Field: "foo",
								}},
							),
						},
					},
					fixity.MultiJson{
						"json type": fixityutil.MustMarshalJsonWithMeta(&foo),
					},
				)
				So(err, ShouldBeNil)

				var fooField fixity.Field
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
