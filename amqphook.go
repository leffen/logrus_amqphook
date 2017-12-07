package logrus_amqphook

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

// VERSION of the application
const VERSION = "1.1.0"

const (
	sleepBetweenFails = time.Second
	maxAttemts        = 10
)

// AmqpHook handles connection properties
type AmqpHook struct {
	connString    string
	exchangeName  string
	routingKey    string
	logInputChan  chan *logrus.Entry
	logOutputChan chan *logrus.Entry
	amqpChan      *amqp.Channel
	Formatter     *GelfJsonFormatter

	AutoDeleteExchange bool
	InternalExchange   bool
	NowaitExchange     bool
}

// NewAmqpHook creates a new hook to logrus
func NewAmqpHook(connString, exchangeName, routingKey string) *AmqpHook {
	hook := &AmqpHook{
		connString:   connString,
		exchangeName: exchangeName,
		routingKey:   routingKey,
		Formatter:    NewFormatter(),
	}

	return hook
}

// Fire delivers entry to amqp
func (hook *AmqpHook) Fire(entry *logrus.Entry) error {
	if hook.amqpChan == nil {
		c, err := hook.buildChannel()
		if err != nil {
			return fmt.Errorf("AmqpEhook.sendEvent>Unable to build channel %s with error %s", hook.exchangeName, err)
		}
		hook.amqpChan = c
	}

	var err error
	logEntry, _ := hook.Formatter.Format(entry)
	attempt := 0

	for {
		attempt++
		if attempt > 1 {
			time.Sleep(sleepBetweenFails) // Let the amqp server rest a little
		}

		if attempt > maxAttemts {
			return fmt.Errorf("Max retries exceeded. Last error %s", err)
		}

		err = hook.amqpChan.Publish(hook.exchangeName, hook.routingKey, false, false, amqp.Publishing{Body: logEntry, DeliveryMode: amqp.Persistent})

		if err != nil {
			logrus.Errorf("AmqpEhook.sendEvent>Unable to publish to %s with error: %s", hook.exchangeName, err)
			continue
		}
		break
	}

	return nil
}

// Levels is available logging levels
func (hook *AmqpHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
		logrus.DebugLevel,
	}
}

func (hook *AmqpHook) buildChannel() (*amqp.Channel, error) {
	conn, err := amqp.Dial(hook.connString)
	if err != nil {
		return nil, err
	}
	amqpChan, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	err = amqpChan.ExchangeDeclare(hook.exchangeName, "fanout", true, hook.AutoDeleteExchange, hook.InternalExchange, hook.NowaitExchange, nil)
	if err != nil {
		return nil, err
	}

	// Clear amqp channel if connection to server is lost
	amqpErrorChan := make(chan *amqp.Error)
	amqpChan.NotifyClose(amqpErrorChan)
	go func(h *AmqpHook, ec chan *amqp.Error) {
		for msg := range ec {
			logrus.Errorf("AmqpHook.buildChannel> Channel Cleanup %s\n", msg)
			h.amqpChan = nil
		}
	}(hook, amqpErrorChan)
	return amqpChan, nil
}
