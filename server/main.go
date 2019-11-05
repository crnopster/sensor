package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/go-redis/redis/v7"
	client "github.com/influxdata/influxdb1-client/v2"
)

func main() {
	influxSaver := flag.Int("influxSaver", 10, "number of saveToInfluxDB goroutines")
	redisSaver := flag.Int("redisSaver", 10, "number of saveToRedisDB goroutines")
	chanSender := flag.Int("chanSender", 10, "number of chanSend goroutines")
	flag.Parse()

	wg := &sync.WaitGroup{}

	clientRedis := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	pong, err := clientRedis.Ping().Result()
	fmt.Println(pong, err)

	clientInflux, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     "http://localhost:8086",
		Username: "sensors",
		Password: "pass",
	})
	if err != nil {
		log.Println(err.Error())
	}
	defer clientInflux.Close()

	tags := make(map[string]string)
	c := make(chan result)
	ci := make(chan result)
	cr := make(chan result)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, c)
	})

	wg.Add(*chanSender)
	for i := 0; i < *chanSender; i++ {
		go chanSend(c, ci, cr, wg)
	}

	wg.Add(*influxSaver)
	for i := 0; i < *influxSaver; i++ {
		go saveToInfluxDB(clientInflux, tags, ci, wg)
	}
	wg.Add(*redisSaver)
	for i := 0; i < *redisSaver; i++ {
		go saveToRedisDB(*clientRedis, wg, cr)
	}

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}

	wg.Wait()

}

type result struct {
	ID          string
	Temperature float32
	Humidity    float32
}

func handler(w http.ResponseWriter, r *http.Request, c chan result) {

	if r.Method == "POST" {
		var res result
		err := json.NewDecoder(r.Body).Decode(&res)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		c <- res

	}
}

func chanSend(c chan result, ci chan result, cr chan result, wg *sync.WaitGroup) {
	defer wg.Done()
	for res := range c {
		ci <- res
		cr <- res
	}
}
