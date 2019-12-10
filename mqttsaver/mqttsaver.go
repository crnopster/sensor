package mqttsaver

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/crnopster/sensor/metric"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// MqttConn .
type MqttConn struct {
	Client mqtt.Client
	C      chan metric.Metric
}

// Save metric to mqtt broker
func (mc *MqttConn) Save(m metric.Metric) {
	mc.C <- m
}

func (mc *MqttConn) mqttWorker(ctx context.Context, wg *sync.WaitGroup, topic string) {
	log.Println("mqttWorker started")
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			log.Println("mqttWorker stopped")
			return
		default:
			for m := range mc.C {
				mc.Client.Publish(topic, 16, false, m)
			}
		}
	}
}

// Worker calls mqttWorker to publish data into topic
func (mc *MqttConn) Worker(ctx context.Context, wg *sync.WaitGroup, workerCount int, topic string) {
	wg.Add(workerCount)
	for i := 0; i < workerCount; i++ {
		go mc.mqttWorker(ctx, wg, topic)
	}

}

// NewMqttClient returns connect to local mqtt broker & metric chan
func NewMqttClient() *MqttConn {
	opts := mqtt.NewClientOptions()
	opts.AddBroker("tcp://localhost:1883")
	opts.SetUsername("sensor")

	client := mqtt.NewClient(opts)
	token := client.Connect()
	for !token.WaitTimeout(3 * time.Second) {
	}
	if err := token.Error(); err != nil {
		log.Fatal(err)
	}
	c := make(chan metric.Metric, 1)
	return &MqttConn{
		Client: client,
		C:      c,
	}
}
