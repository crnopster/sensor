package main

import (
	"context"
	"sync"

	"github.com/crnopster/sensor/influxsaver"
	"github.com/crnopster/sensor/mqttsaver"
	"github.com/crnopster/sensor/redissaver"
)

func start(ctx context.Context, wg *sync.WaitGroup, redisWorkers, influxWorkers, mqttWorkers int, topic string) (*influxsaver.InfluxConn, *redissaver.RedisConn, *mqttsaver.MqttConn) {
	// connect to influx
	i := influxsaver.NewInfluxClient()
	defer i.Client.Close()

	// connect to redis
	r := redissaver.NewRedisClient()
	defer r.Client.Close()

	// connect to mqtt broker
	mc := mqttsaver.NewMqttClient()
	defer mc.Client.Disconnect(100)

	r.Worker(ctx, wg, redisWorkers)

	i.Worker(ctx, wg, influxWorkers)

	mc.Worker(ctx, wg, mqttWorkers, topic)

	return i, r, mc

}
