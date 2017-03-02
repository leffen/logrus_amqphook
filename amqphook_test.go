package logrus_amqphook

import (
	"github.com/Sirupsen/logrus"
	"testing"
	"time"
)

func TestHook(t *testing.T) {
	hook := NewAmqpHook( "amqps://ca.paxtonagent.logger:k7qRMWe5wa@mq1.cloudaccess.io", "ca.logging.e.prod","")
	logrus.SetFormatter(new(LogstashFormatter))
	logrus.AddHook(hook)
	logrus.Errorf("Something broke...")
	time.Sleep(1 * time.Second)
}
