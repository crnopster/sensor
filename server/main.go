package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"sync"

	"github.com/crnopster/sensor/metric"
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

	wg.Add(2)
	srv := http.Server{Addr: ":8080"}
	// Graceful shutdown
	go shutdown(ctx, cancel, &srv, wg)

	i, r, mc := start(ctx, wg, *redisWorkers, *influxWorkers, *mqttWorkers, *topic)

	defer i.Client.Close()
	defer r.Client.Close()
	defer mc.Client.Disconnect(100)

	s := server{storages: []saver{i, r, mc}}

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
