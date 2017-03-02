package logrus_amqphook

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/Sirupsen/logrus"

  "path"
)


type GelfJsonFormatter struct {
	Source          string
	Type            string // If not empty use for logstash type field.
	TimestampFormat string // TimestampFormat sets the format used for timestamps.
  HostName        string
  Facility        string
}

func NewFormatter() *GelfJsonFormatter {
	hostname, _ := os.Hostname()

	return &GelfJsonFormatter{
		Source:          hostname,
		TimestampFormat: time.RFC3339Nano,
    HostName:         hostname,
    Facility:         path.Base(os.Args[0]),
	}
}

func (f *GelfJsonFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	entry.Data["version"] = "1.1"
	entry.Data["timestamp"] = float64(time.Now().UnixNano()/1000000) / 1000.
	entry.Data["short_message"] = entry.Message
	entry.Data["host"] = f.HostName
	entry.Data["level"] = entry.Level.String()
	if f.Type != "" { entry.Data["type"] = f.Type }
  entry.Data["source"] = f.Source

	serialized, err := json.Marshal(entry.Data)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal fields to JSON, %v", err)
	}

	return append(serialized, '\n'), nil
}
