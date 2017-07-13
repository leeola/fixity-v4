package sync

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSyncPathResolveDirs(t *testing.T) {
	Convey("Scenario: Resolving sync directory structure", t, func() {
		// pairs of input/output.
		// if all three expected strings are nil, error is expected
		scenarios := [][]string{
			// 0
			[]string{"dir", "file", ""}, // input
			[]string{"", "", ""},        // error: dir must be absolute
			// 1
			[]string{"/", "", ""}, // input
			[]string{"", "", ""},  // error: can't resolve folder with no parent
			// 2
			[]string{"/foo", "", "/foo"}, // input
			[]string{"", "", ""},         // error: folder can't be absolute
			// 3
			[]string{"/foo/dir", "", ""},
			[]string{"/foo/dir", "/foo/dir", "dir"},
			// 4
			[]string{"/foo", "bar", ""},         // input
			[]string{"/foo", "/foo/bar", "foo"}, // expected
			// 5
			[]string{"/foo", "bar", "folder"},      // input
			[]string{"/foo", "/foo/bar", "folder"}, // expected
			// 6
			[]string{"/foo/bar/baz", "", ""},                // input
			[]string{"/foo/bar/baz", "/foo/bar/baz", "baz"}, // expected
		}

		for i, inScenario := range scenarios {
			if i%2 != 0 {
				continue
			}
			testNum := i / 2
			exScenario := scenarios[i+1]
			inDir, inFile, inFolder := inScenario[0], inScenario[1], inScenario[2]
			exTrim, exPath, exFolder := exScenario[0], exScenario[1], exScenario[2]
			expectsError := exTrim == "" && exPath == "" && exFolder == ""

			conveyName := fmt.Sprintf("Given test #%d with a dir:%q, file:%q and folder:%q",
				testNum, inDir, inFile, inFolder)
			Convey(conveyName, func() {
				Convey("When run", func() {
					trimPath, path, folder, err := resolveDirs(inDir, inFile, inFolder)

					Convey("Then it should resolve in a trimPath of "+exTrim, func() {
						So(trimPath, ShouldEqual, exTrim)
					})
					Convey("Then it should resolve in a path of "+exPath, func() {
						So(path, ShouldEqual, exPath)
					})
					Convey("Then it should resolve in a folder of "+exFolder, func() {
						So(folder, ShouldEqual, exFolder)
					})
					if expectsError {
						Convey("Then it should return an error", func() {
							So(err, ShouldNotBeNil)
						})
					} else {
						Convey("Then it should not return an error", func() {
							So(err, ShouldBeNil)
						})
					}
				})
			})
		}
	})
}
