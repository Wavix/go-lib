package metrics

import (
	"encoding/json"
	"log"

	"ci.unitedline.net/wavix/go-private-sdk/amqp"
	"github.com/wagslane/go-rabbitmq"
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

func Init(amqpDsn string, service string) *MetricInstance {
	connection := amqp.OpenAmqpConnection(amqpDsn)
	if connection == nil {
		log.Println("Failed to open AMQP connection")
		return nil
	}

	publisher := amqp.StartAmqpPublisher(connection)

	if publisher == nil {
		log.Println("Failed to start AMQP publisher")
	}

	return &MetricInstance{
		Connection: publisher,
		Service:    service,
	}
}

func (m *MetricInstance) IncrementCounter(key string) {
	payload := MetricPayload{Key: key, Value: 1, Type: "counter"}
	data, _ := json.Marshal(payload)

	err := m.Connection.Publish([]byte(data), []string{RoutingKey}, rabbitmq.WithPublishOptionsContentType("application/json"))
	if err != nil {
		log.Println("Failed to publish counter")
	}
}

func (m *MetricInstance) SetGauge(key string, value int) {
	payload := MetricPayload{Key: key, Value: value, Type: "gauge"}
	data, _ := json.Marshal(payload)

	err := m.Connection.Publish([]byte(data), []string{RoutingKey}, rabbitmq.WithPublishOptionsContentType("application/json"))
	if err != nil {
		log.Println("Failed to publish gauge")
	}
}
