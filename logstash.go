package logstash

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/Sirupsen/logrus"
)

// Formatter generates json in logstash format.
// Logstash site: http://logstash.net/
type LogstashFormatter struct {
	Source          string
	Type            string // If not empty use for logstash type field.
	TimestampFormat string // TimestampFormat sets the format used for timestamps.
}

func NewLogstashFormatter(logstashType string) *LogstashFormatter {
	hostname, _ := os.Hostname()

	return &LogstashFormatter{
		Source:          hostname,
		TimestampFormat: time.RFC3339Nano,
		Type:            logstashType,
	}
}

func (f *LogstashFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	entry.Data["@version"] = 1

	if f.TimestampFormat == "" {
		f.TimestampFormat = logrus.DefaultTimestampFormat
	}

	entry.Data["timestamp"] = entry.Time.Format(f.TimestampFormat)

	// set message field
	entry.Data["message"] = entry.Message

	// set level field
	entry.Data["level"] = entry.Level.String()

	// set type field
	if f.Type != "" {
		entry.Data["type"] = f.Type
	}

	if f.Source != "" {
		entry.Data["source"] = f.Source
	}

	serialized, err := json.Marshal(entry.Data)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal fields to JSON, %v", err)
	}
	return append(serialized, '\n'), nil
}
