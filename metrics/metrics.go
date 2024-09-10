package metrics

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/process"
	"github.com/wagslane/go-rabbitmq"
	"github.com/wavix/go-lib/amqp"
)

type MetricInstance struct {
	Connection *rabbitmq.Publisher
	Service    string
}

type MetricPayload struct {
	Service string `json:"service"`
	Key     string `json:"key"`
	Value   int    `json:"value"`
	Type    string `json:"type"`
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

	go instance.startSystemStats()

	return instance, nil
}

func (m *MetricInstance) IncrementCounter(key string) {
	if m == nil || m.Connection == nil {
		return
	}

	payload := MetricPayload{Service: m.Service, Key: key, Value: 1, Type: "counter"}
	data, err := json.Marshal(payload)
	if err != nil {
		log.Println("Failed to marshal counter")
		return
	}

	err = m.Connection.Publish([]byte(data), []string{RoutingKey}, rabbitmq.WithPublishOptionsContentType("application/json"))
	if err != nil {
		log.Println("Failed to publish counter")
	}
}

func (m *MetricInstance) IncrementCounterBy(key string, amount int) {
	if m == nil || m.Connection == nil {
		return
	}

	payload := MetricPayload{Service: m.Service, Key: key, Value: amount, Type: "counter"}
	data, err := json.Marshal(payload)
	if err != nil {
		log.Println("Failed to marshal counter")
		return
	}

	err = m.Connection.Publish([]byte(data), []string{RoutingKey}, rabbitmq.WithPublishOptionsContentType("application/json"))
	if err != nil {
		log.Println("Failed to publish counter")
	}
}

func (m *MetricInstance) SetGauge(key string, value int) {
	if m == nil || m.Connection == nil {
		return
	}

	payload := MetricPayload{Service: m.Service, Key: key, Value: value, Type: "gauge"}
	data, err := json.Marshal(payload)

	if err != nil {
		log.Println("Failed to marshal gauge")
		return
	}

	err = m.Connection.Publish([]byte(data), []string{RoutingKey}, rabbitmq.WithPublishOptionsContentType("application/json"))
	if err != nil {
		log.Println("Failed to publish gauge")
	}
}

func (m *MetricInstance) startSystemStats() {
	if m == nil || m.Connection == nil {
		return
	}

	for {
		cpu, mem, err := getCurrentProcessUsage()
		if err != nil {
			log.Println("Error in application cpu/mem stats: ", err)
			time.Sleep(10 * time.Second)
		}

		m.SetGauge("cpu", int(cpu))
		m.SetGauge("mem", int(mem))

		time.Sleep(10 * time.Second)
	}

}

func getCurrentProcessUsage() (int64, int64, error) {
	pid := os.Getpid()
	process, err := process.NewProcess(int32(pid))
	if err != nil {
		return 0, 0, err
	}

	cpuPercent, err := process.CPUPercent()
	if err != nil {
		return 0, 0, err
	}

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	vmem, err := mem.VirtualMemory()
	if err != nil {
		return 0, 0, err
	}

	memPercent := float32(memStats.Alloc) / float32(vmem.Total) * 100

	return formatFloat(float32(cpuPercent)), formatFloat(memPercent), nil
}

func formatFloat(num float32) int64 {
	formated := fmt.Sprintf("%.2f", num)
	result, _ := strconv.ParseFloat(formated, 32)
	return int64(result)
}
