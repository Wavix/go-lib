package metrics

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/wagslane/go-rabbitmq"
	"github.com/wavix/go-lib/amqp"
)

type MetricInstance struct {
	Connection *rabbitmq.Publisher
	Service    string
}

type MetricPayload struct {
	Key   string `json:"key"`
	Value int    `json:"value"`
	Type  string `json:"type"`
}

const RoutingKey = "metrics_exporter"

func Init(amqpDsn string, service string) (*MetricInstance, error) {
	connection := amqp.OpenAmqpConnection(amqpDsn)
	if connection == nil {
		return nil, errors.New("failed to open AMQP connection")
	}

	publisher := amqp.StartAmqpPublisher(connection)

	if publisher == nil {
		return nil, errors.New("failed to start AMQP publisher")
	}

	instance := &MetricInstance{
		Connection: publisher,
		Service:    service,
	}

	return instance, nil
}

func (m *MetricInstance) IncrementCounter(key string) {
	if m == nil || m.Connection == nil {
		return
	}

	payload := MetricPayload{Key: key, Value: 1, Type: "counter"}
	data, _ := json.Marshal(payload)

	err := m.Connection.Publish([]byte(data), []string{RoutingKey}, rabbitmq.WithPublishOptionsContentType("application/json"))
	if err != nil {
		log.Println("Failed to publish counter")
	}
}

func (m *MetricInstance) SetGauge(key string, value int) {
	if m == nil || m.Connection == nil {
		return
	}

	payload := MetricPayload{Key: key, Value: value, Type: "gauge"}
	data, _ := json.Marshal(payload)

	err := m.Connection.Publish([]byte(data), []string{RoutingKey}, rabbitmq.WithPublishOptionsContentType("application/json"))
	if err != nil {
		log.Println("Failed to publish gauge")
	}
}
