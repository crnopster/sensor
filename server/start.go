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

	// connect to redis
	r := redissaver.NewRedisClient()

	// connect to mqtt broker
	mc := mqttsaver.NewMqttClient()

	r.Worker(ctx, wg, redisWorkers)

	i.Worker(ctx, wg, influxWorkers)

	mc.Worker(ctx, wg, mqttWorkers, topic)

	return i, r, mc

}
