package logrus_amqphook

import (
	"github.com/Sirupsen/logrus"
	"testing"
	"time"
	"os"
	_ "github.com/joho/godotenv/autoload"
)

func TestHook(t *testing.T) {
	hook := NewAmqpHook( os.Getenv("TEST_CONNECTION"), "graylog_test","#")
	logrus.SetFormatter(new(GelfJsonFormatter))
	logrus.AddHook(hook)
	logrus.Errorf("Sender en feil her gitt")
	time.Sleep(1 * time.Second)
}
