package router

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestReverse(t *testing.T) {
	Convey("Test reverse", t, func() {
		reverse = map[string]string{
			"static": "/test/static/param",
			"param":  "/test/:static/:param",
			"catch":  "/test/:static/:param/*all",
		}

		d, err := Path("static", map[string]string{"ignore": "me"})
		So(err, ShouldBeNil)
		So(d, ShouldEqual, mountPoint.String()+"/test/static/param")

		d, err = Path("param", map[string]string{"static": "hi", "param": "bye"})
		So(err, ShouldBeNil)
		So(d, ShouldEqual, mountPoint.String()+"/test/hi/bye")

		d, err = Path("param", map[string]string{"static": "hi"})
		So(err, ShouldNotBeNil)

		d, err = Path("catch", map[string]string{"static": "hi", "param": "bye"}, "aaa", "bbb")
		So(err, ShouldBeNil)
		So(d, ShouldEqual, mountPoint.String()+"/test/hi/bye/aaa/bbb")

	})
}
