package logrus_amqphook

import (
	"os"
	"testing"
	"time"

	"github.com/Sirupsen/logrus"
	_ "github.com/joho/godotenv/autoload"
)

func TestHook(t *testing.T) {
	hook := NewAmqpHook(os.Getenv("TEST_CONNECTION"), "graylog_test", "#")
	logrus.SetFormatter(new(GelfJsonFormatter))
	logrus.AddHook(hook)
	logrus.Errorf("Sender en feil her gitt")
	logrus.WithFields(logrus.Fields{
		"Server UID":  "UID",
		"Camera":      "Name",
		"Camera IP":   "SrvIP",
		"Camera port": "SrvPort",
	}).Errorf("Sender en PARAM feil her gitt")
	time.Sleep(1 * time.Second)
}
