package jwt

import (
	"testing"
	"time"

	"github.com/clickyab/services/config"
	"github.com/clickyab/services/initializer"
	. "github.com/smartystreets/goconvey/convey"
)

func TestJwt(t *testing.T) {
	config.Initialize("test", "test", "test")
	defer initializer.Initialize()()

	t.Log(publicPem.String())
	Convey("JWT should", t, func() {
		data := map[string]string{
			"Name":   "Tester",
			"Status": "immortal",
			"target": "services",
		}

		j := NewJWT()
		Convey(" be valid", func() {
			ll := j.Encode(data, 5*time.Minute)
			a, b, err := j.Decode([]byte(ll), []string{"Name", "Status", "target"}...)
			So(err, ShouldBeNil)
			So(a, ShouldBeFalse)
			So(len(b), ShouldEqual, 3)
			So(b["Name"], ShouldEqual, "Tester")
			So(b["name"], ShouldNotEqual, "Tester")
			So(b["Status"], ShouldEqual, "immortal")
			So(b["target"], ShouldEqual, "services")
		})
		Convey(" expired", func() {
			ll := j.Encode(data, -5*time.Minute)
			a, b, err := j.Decode([]byte(ll), []string{"Name", "Status", "target"}...)
			So(err, ShouldBeNil)
			So(a, ShouldBeFalse)
			So(len(b), ShouldEqual, 3)
			So(b["Name"], ShouldEqual, "Tester")
			So(b["name"], ShouldNotEqual, "Tester")
			So(b["Status"], ShouldEqual, "immortal")
			So(b["target"], ShouldEqual, "services")
		})
		Convey(" return err", func() {
			ll := j.Encode(data, -5*time.Minute)
			a, b, err := j.Decode([]byte(ll), []string{"FirstName", "Status", "target"}...)
			So(err, ShouldNotBeNil)
			So(b, ShouldBeNil)
			So(a, ShouldBeFalse)
		})

		Convey("get all", func() {
			ll := j.Encode(data, 5*time.Minute)
			a, b, err := j.Decode([]byte(ll))
			So(err, ShouldBeNil)
			So(a, ShouldBeFalse)
			So(len(b), ShouldEqual, 3)
			So(b["Name"], ShouldEqual, "Tester")
			So(b["name"], ShouldNotEqual, "Tester")
			So(b["Status"], ShouldEqual, "immortal")
			So(b["target"], ShouldEqual, "services")
		})
	})
}
