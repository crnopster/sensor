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

	redisWorkers := flag.Int("redisWorkers", 1, "number of redisworker goroutines")
	influxWorkers := flag.Int("influxWorkers", 1, "number of influxworker goroutines")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	// connect to redis
	r := redissaver.NewRedisClient()
	defer r.Client.Close()
	// start redisworkers
	r.Worker(ctx, wg, *redisWorkers)

	// connect to influx
	i := influxsaver.NewInfluxClient()
	defer i.Client.Close()
	// start influxworkers
	i.Worker(ctx, wg, *influxWorkers)

	wg.Add(2)

	srv := http.Server{Addr: ":8080"}
	// Graceful shutdown
	go shutdown(ctx, cancel, &srv, wg)

	s := server{storages: []saver{i}}

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
