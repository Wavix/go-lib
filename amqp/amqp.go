package amqp

import (
	"log"

	"github.com/wagslane/go-rabbitmq"
)

const topic = "wavix.topic"

func OpenAmqpConnection(amqp string) *rabbitmq.Conn {
	conn, err := rabbitmq.NewConn(
		amqp,
		rabbitmq.WithConnectionOptionsLogging,
	)
	if err != nil {
		log.Println(err)
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
	)
	if err != nil {
		log.Println(err)
		return nil
	}

	err = consumer.Run(handler)
	if err != nil {
		log.Println(err)
		return nil
	}

	return consumer
}

func StartAmqpPublisher(conn *rabbitmq.Conn) *rabbitmq.Publisher {
	publisher, err := rabbitmq.NewPublisher(
		conn,
		rabbitmq.WithPublisherOptionsLogging,
		rabbitmq.WithPublisherOptionsExchangeName(topic),
	)
	if err != nil {
		log.Println(err)
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
		log.Println(err)
	}

	return err != nil
}
