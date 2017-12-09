package healthz

import (
	"errors"
	"testing"

	"context"
	"net/http"
	"sync"

	"net/http/httptest"

	"io/ioutil"

	"github.com/sirupsen/logrus"
	"github.com/smartystreets/goconvey/convey"
)

type mysqlHealth struct {
	Err error
}

type redisHealth struct {
	Err error
}

type brokerHealth struct {
	Err error
}

func (h *mysqlHealth) Healthy(context.Context) error {
	return h.Err
}

func (h *redisHealth) Healthy(context.Context) error {
	return h.Err
}

func (h *brokerHealth) Healthy(context.Context) error {
	return h.Err
}

type MyHandler struct {
	sync.Mutex
}

func (h *MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Lock()
	defer h.Unlock()
	rout := route{}
	ctx := context.Background()
	rout.check(ctx, w, r)

}

func TestHealth(t *testing.T) {
	handler := &MyHandler{}
	server := httptest.NewServer(handler)
	defer server.Close()
	m := &mysqlHealth{
		Err: nil,
	}
	r := &mysqlHealth{
		Err: nil,
	}
	b := &mysqlHealth{
		Err: nil,
	}
	Register(m, r, b)
	convey.Convey("test with all ok", t, func() {

		resp, err := http.Get(server.URL)
		if err != nil {
			t.Fatal(err)
		}
		//logrus.Fatal(m.Err,r.Err,b.Err)
		convey.So(resp.StatusCode, convey.ShouldEqual, http.StatusOK)

	})

	convey.Convey("mysql error", t, func() {
		m.Err = errors.New("mysql error here")
		logrus.Warn(m.Err)
		resp, err := http.Get(server.URL)
		if err != nil {
			t.Fatal(err)
		}
		msg, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		convey.So(resp.StatusCode, convey.ShouldEqual, http.StatusInternalServerError)
		convey.So(string(msg), convey.ShouldEqual, "mysql error here")

	})

	convey.Convey("mysql and redis error", t, func() {
		m.Err = errors.New("mysql error here")
		r.Err = errors.New("redis error here")
		resp, err := http.Get(server.URL)
		if err != nil {
			t.Fatal(err)
		}
		msg, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		convey.So(resp.StatusCode, convey.ShouldEqual, http.StatusInternalServerError)
		convey.So(string(msg), convey.ShouldEqual, "mysql error hereredis error here")

	})

}
