package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"sync"
	"time"

	client "github.com/influxdata/influxdb1-client/v2"
)

func main() {
	influxSaver := flag.Int("influxSaver", 10, "number of saveToInfluxDB goroutines")
	flag.Parse()

	wg := &sync.WaitGroup{}

	con, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     "http://localhost:8086",
		Username: "sensors",
		Password: "pass",
	})
	if err != nil {
		log.Println(err.Error())
	}
	defer con.Close()

	tags := make(map[string]string)
	c := make(chan result)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, c)
	})

	wg.Add(*influxSaver)
	for i := 0; i < *influxSaver; i++ {
		go saveToInfluxDB(con, tags, c, wg)
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

func saveToInfluxDB(con client.Client, tags map[string]string, c chan result, wg *sync.WaitGroup) {
	fields := make(map[string]interface{})

	defer wg.Done()

	bp, _ := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  "sensor",
		Precision: "",
	})

	for res := range c {
		fields["id"] = res.ID
		fields["temperature"] = res.Temperature
		fields["humidity"] = res.Humidity
		log.Println(res.ID, res.Temperature, res.Humidity)
		newPoint, err := client.NewPoint(
			"sensor",
			tags,
			fields,
			time.Now(),
		)
		if err != nil {
			log.Println(err.Error())
		}
		bp.AddPoint(newPoint)
		err = con.Write(bp)
		if err != nil {
			log.Println(err.Error())
		}
	}
}
