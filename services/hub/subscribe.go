package hub

import "github.com/olebedev/emitter"

var engine = emitter.Emitter{}

// Subscribe in a topic
func Subscribe(topic string, middlewares ...func(*emitter.Event)) <-chan emitter.Event {
	return engine.On(topic, middlewares...)
}

// Publish a message in a topic
func Publish(topic string, msg ...interface{}) {
	<-engine.Emit(topic, msg...)
}
