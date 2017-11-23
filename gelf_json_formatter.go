package logrus_amqphook

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"

	"path"
	"sync"
)

type GelfJsonFormatter struct {
	Source   string
	Type     string
	HostName string
	Facility string
	wg       sync.WaitGroup
	mu       sync.RWMutex
}

func NewFormatter() *GelfJsonFormatter {
	hostname, _ := os.Hostname()

	return &GelfJsonFormatter{
		Source:   hostname,
		HostName: Hostname(),
		Facility: path.Base(os.Args[0]),
	}
}

func Hostname() string {
	name := os.Getenv("HOST_HOSTNAME")
	if len(name) > 0 {
		return name
	}
	name, _ = os.Hostname()
	return name
}

func (f *GelfJsonFormatter) Format(entry *logrus.Entry) ([]byte, error) {

	f.mu.RLock() // Claim the mutex as a RLock - allowing multiple go routines to log simultaneously
	defer f.mu.RUnlock()

	newData := make(map[string]interface{})
	for k, v := range entry.Data {
		newData[fmt.Sprintf("_%s", k)] = v
	}

	newData["version"] = "1.1"
	newData["timestamp"] = float64(time.Now().UnixNano()/1000000) / 1000.
	newData["short_message"] = entry.Message
	newData["host"] = f.HostName
	newData["facility"] = f.Facility
	newData["level"] = entry.Level
	newData["_level_str"] = entry.Level.String()
	newData["_source"] = f.Source

	serialized, err := json.Marshal(newData)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal fields to JSON, %v", err)
	}

	return append(serialized, '\n'), nil
}
