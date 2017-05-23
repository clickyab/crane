package workers

import (
	"fmt"
	"strconv"
	"time"

	"services/assert"
)

// TODO get this from config
const limit = 1000

type Acknowledger interface {
	Ack(multiple bool) error
	// Nack negatively acknowledge the delivery of message(s) identified by the delivery tag from either the client or server.
	Nack(multiple, requeue bool) error
	// Reject delegates a negatively acknowledgement through the Acknowledger interface.
	Reject(requeue bool) error
}

func genID(t time.Time, s ...string) string {
	m := t.Format("2006010203")
	return fmt.Sprint(m, s)
}

func factID(tm time.Time) int {

	layout := "2006-01-02T15:04:05.000Z"
	str := "2017-03-21T00:00:00.000Z"
	t, err := time.Parse(layout, str)

	if err != nil {
		fmt.Println(err)
	}
	assert.Nil(err)
	return int(tm.Sub(t).Hours()) + 1

}

func timestampToTime(s string) time.Time {
	i, err := strconv.ParseInt(s, 10, 0)
	assert.Nil(err)
	return time.Unix(i, 0)
}
