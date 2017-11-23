package logrus_amqphook

import (
	"log"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

// VERSION of the application
const VERSION = "1.0.4"

const (
	bufferSize        = 1000
	sleepBetweenFails = 10 * time.Second
)

type AmqpHook struct {
	connString    string
	exchangeName  string
	routingKey    string
	logInputChan  chan *logrus.Entry
	logOutputChan chan *logrus.Entry
	amqpChan      *amqp.Channel
	Formatter     *GelfJsonFormatter
}

func NewAmqpHook(connString, exchangeName, routingKey string) *AmqpHook {
	hook := &AmqpHook{
		connString:    connString,
		exchangeName:  exchangeName,
		routingKey:    routingKey,
		logInputChan:  make(chan *logrus.Entry),
		logOutputChan: make(chan *logrus.Entry, bufferSize),
		Formatter:     NewFormatter(),
	}

	rb := newRingBuffer(hook.logInputChan, hook.logOutputChan)
	go rb.Run()
	go hook.handle()

	return hook
}

func (hook *AmqpHook) Fire(entry *logrus.Entry) error {
	//hook.logInputChan <- entry
	// hook.logOutputChan <- entry
	hook.sendEvent(entry)
	return nil
}

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

func (hook *AmqpHook) handle() {
	for msg := range hook.logOutputChan {
		hook.sendEvent(msg)
	}
}

func (hook *AmqpHook) sendEvent(entry *logrus.Entry) {
	logEntry, _ := hook.Formatter.Format(entry)
	attempt := 0

	for {
		attempt++
		if attempt > 1 {
			time.Sleep(sleepBetweenFails) // Let the amqp server rest a little
		}

		if hook.amqpChan == nil {
			c, err := hook.buildChannel()
			if err != nil {
				log.Println(err)
				continue
			}
			hook.amqpChan = c
		}

		if err := hook.amqpChan.Publish(hook.exchangeName, hook.routingKey, false, false, amqp.Publishing{
			Body:         []byte(logEntry),
			DeliveryMode: amqp.Persistent,
		}); err != nil {
			log.Println("Unable to publish:", err)
			continue
		}
		break
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

	err = amqpChan.ExchangeDeclare(hook.exchangeName, "fanout", true, true, true, true, nil)
	if err != nil {
		return nil, err
	}

	// Clear amqp channel if connection to server is lost
	amqpErrorChan := make(chan *amqp.Error)
	amqpChan.NotifyClose(amqpErrorChan)
	go func(h *AmqpHook, ec chan *amqp.Error) {
		for msg := range ec {
			log.Println(msg)
			h.amqpChan = nil
		}
	}(hook, amqpErrorChan)
	return amqpChan, err
}
