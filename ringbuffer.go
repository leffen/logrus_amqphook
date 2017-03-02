package amqphook

import (
	"github.com/Sirupsen/logrus"
)

type ringBuffer struct {
	inputChannel  <-chan *logrus.Entry
	outputChannel chan *logrus.Entry
}

func newRingBuffer(inputChannel <-chan *logrus.Entry, outputChannel chan *logrus.Entry) *ringBuffer {
	return &ringBuffer{inputChannel, outputChannel}
}

func (r *ringBuffer) Run() {
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
