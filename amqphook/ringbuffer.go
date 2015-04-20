package amqphook

import (
	"github.com/Sirupsen/logrus"
)

type RingBuffer struct {
	inputChannel  <-chan *logrus.Entry
	outputChannel chan *logrus.Entry
}

func NewRingBuffer(inputChannel <-chan *logrus.Entry, outputChannel chan *logrus.Entry) *RingBuffer {
	return &RingBuffer{inputChannel, outputChannel}
}

func (r *RingBuffer) Run() {
	for v := range r.inputChannel {
		select {
		case r.outputChannel <- v:
		default:
			<-r.outputChannel
			r.outputChannel <- v
		}
	}
	close(r.outputChannel)
}
