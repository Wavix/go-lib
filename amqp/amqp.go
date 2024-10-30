package amqp

import (
	"os"

	"github.com/wavix/go-lib/logger"

	"github.com/wagslane/go-rabbitmq"
)

const topic = "wavix.topic"

type LoggerType interface {
	Fatalf(string, ...interface{})
	Errorf(string, ...interface{})
	Warnf(string, ...interface{})
	Infof(string, ...interface{})
	Debugf(string, ...interface{})
}

type CustomLogger struct{}

var log = logger.New("AMQP")

func (l CustomLogger) Fatalf(format string, v ...interface{}) {
	log.Error().Msgf(format, v...)
	os.Exit(1)
}

func (l CustomLogger) Errorf(format string, v ...interface{}) {
	log.Error().Msgf(format, v...)
}

func (l CustomLogger) Warnf(format string, v ...interface{}) {
	log.Warn().Msgf(format, v...)
}

func (l CustomLogger) Infof(format string, v ...interface{}) {
	log.Info().Msgf(format, v...)
}

func (l CustomLogger) Debugf(format string, v ...interface{}) {
	log.Debug().Msgf(format, v...)
}

func OpenAmqpConnection(amqp string) *rabbitmq.Conn {
	conn, err := rabbitmq.NewConn(
		amqp,
		rabbitmq.WithConnectionOptionsLogger(CustomLogger{}),
	)
	if err != nil {
		log.Error().Msgf("Failed to open AMQP connection: %v", err)
		return nil
	}

	return conn
}

func StartAmqpConsumer(conn *rabbitmq.Conn, handler func(rabbitmq.Delivery) rabbitmq.Action, routingKey string, queue string) *rabbitmq.Consumer {
	consumer, err := rabbitmq.NewConsumer(
		conn,
		queue,
		rabbitmq.WithConsumerOptionsRoutingKey(routingKey),
		rabbitmq.WithConsumerOptionsExchangeName(topic),
		rabbitmq.WithConsumerOptionsConsumerAutoAck(true),
		rabbitmq.WithConsumerOptionsLogger(CustomLogger{}),
	)
	if err != nil {
		log.Error().Msgf("Failed to start AMQP consumer: %v", err)
		return nil
	}

	err = consumer.Run(handler)
	if err != nil {
		log.Error().Msgf("Failed to run AMQP consumer: %v", err)
		return nil
	}

	return consumer
}

func StartAmqpPublisher(conn *rabbitmq.Conn) *rabbitmq.Publisher {
	publisher, err := rabbitmq.NewPublisher(
		conn,
		rabbitmq.WithPublisherOptionsLogger(CustomLogger{}),
		rabbitmq.WithPublisherOptionsExchangeName(topic),
	)
	if err != nil {
		log.Error().Msgf("Failed to start AMQP publisher: %v", err)
		return nil
	}

	return publisher
}

func PublishAmqpMessage(publisher *rabbitmq.Publisher, routingKey string, message []byte) bool {
	err := publisher.Publish(
		message,
		[]string{routingKey},
		rabbitmq.WithPublishOptionsContentType("application/json"),
		rabbitmq.WithPublishOptionsExchange(topic),
	)
	if err != nil {
		log.Error().Msgf("Failed to publish AMQP message: %v", err)
	}

	return err != nil
}
