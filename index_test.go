package fixity

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestBlockchainAppendContent(t *testing.T) {
	Convey("Scenario: Testing equality of Fields", t, func() {
		Convey("When fields are equal", func() {
			fieldsss := [][]Fields{
				[]Fields{
					Fields{{
						Field: "foo",
						Value: "bar",
					}},
					Fields{{
						Field: "foo",
						Value: "bar",
					}},
				},
			}

			Convey("Then Equal() should return true", func() {
				for _, fieldss := range fieldsss {
					a := fieldss[0]
					b := fieldss[1]
					So(a.Equal(b), ShouldBeTrue)
					So(b.Equal(a), ShouldBeTrue)
				}
			})
		})

		Convey("When fields are not equal", func() {
			fieldsss := [][]Fields{
				[]Fields{
					Fields{{
						Field: "foo",
						Value: "bar",
					}},
					Fields{{
						Field: "foo",
						Value: "baz",
					}},
				},
			}

			Convey("Then Equal() should return false", func() {
				for _, fieldss := range fieldsss {
					a := fieldss[0]
					b := fieldss[1]
					So(a.Equal(b), ShouldBeFalse)
					So(b.Equal(a), ShouldBeFalse)
				}
			})
		})
	})
}
