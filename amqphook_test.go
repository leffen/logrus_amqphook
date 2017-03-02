package logrus_amqphook

import (
	"github.com/Sirupsen/logrus"
	"github.com/Sirupsen/logrus/formatters/logstash"
	"testing"
	"time"
)

func TestHook(t *testing.T) {
	hook := NewAmqpHook("TestApp", "amqps://ca.paxtonagent.logger:k7qRMWe5wa@mq1.cloudaccess.io", "ca.logging.e.prod", "logrus.")
	logrus.SetFormatter(new(logstash.LogstashFormatter))
	logrus.AddHook(hook)
	logrus.Errorf("Something broke...")
	time.Sleep(1 * time.Second)
}
