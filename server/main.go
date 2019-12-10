package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"sync"

	"github.com/crnopster/sensor/influxsaver"
	"github.com/crnopster/sensor/metric"
	"github.com/crnopster/sensor/mqttsaver"
	"github.com/crnopster/sensor/redissaver"
)

type saver interface {
	//save calls save method for all storages
	Save(metric.Metric)
}

type server struct {
	storages []saver
}

func main() {

	redisWorkers := flag.Int("redisWorkers", 5, "number of redisworker goroutines")
	influxWorkers := flag.Int("influxWorkers", 5, "number of influxworker goroutines")
	mqttWorkers := flag.Int("mqttWorkers", 5, "number of mqttWorker goroutines")
	topic := flag.String("topic", "test", "mqtt topic")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	// connect to influx
	i := influxsaver.NewInfluxClient()
	defer i.Client.Close()
	// start influxWorkers
	i.Worker(ctx, wg, *influxWorkers)

	// connect to redis
	r := redissaver.NewRedisClient()
	defer r.Client.Close()
	// start redisWorkers
	r.Worker(ctx, wg, *redisWorkers)

	// connect to mqtt broker
	mc := mqttsaver.NewMqttClient()
	defer mc.Client.Disconnect(100)
	// start mqttWorkers
	mc.Worker(ctx, wg, *mqttWorkers, *topic)

	wg.Add(2)
	srv := http.Server{Addr: ":8080"}
	// Graceful shutdown
	go shutdown(ctx, cancel, &srv, wg)

	s := server{storages: []saver{i, r}}

	http.HandleFunc("/", s.handler)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

	wg.Wait()
}

func (s *server) handler(w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		var m metric.Metric
		err := json.NewDecoder(req.Body).Decode(&m)
		if err != nil {
			log.Println(err.Error())
		}
		for _, storage := range s.storages {
			storage.Save(m)
		}
	}
}
