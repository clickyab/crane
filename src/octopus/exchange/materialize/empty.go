package materialize

import "github.com/Sirupsen/logrus"

type job struct {
	data  []byte
	topic string
	key   string
}

func (j job) Encode() ([]byte, error) {
	return j.data, nil
}

func (j job) Length() int {
	return len(j.data)
}

func (j job) Topic() string {
	return j.topic
}

func (j job) Key() string {
	return j.key
}

func (j job) Report() func(error) {
	return func(err error) {
		logrus.Debug(err)
	}
}
