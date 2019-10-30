package main

import (
	"log"
	"net/http"
	"sync"
	"time"

	client "github.com/influxdata/influxdb1-client/v2"
)

func main() {
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
	fields := make(map[string]interface{})
	c := make(chan result)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, c)
	})
	go saveToInfluxDB(con, tags, fields, c, wg)

	wg.Wait()

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

type result struct {
	id          string
	temperature string
	humidity    string
}

func handler(w http.ResponseWriter, r *http.Request, c chan result) {

	if r.Method == "POST" {
		id := r.FormValue("id")
		temperature := r.FormValue("temperature")
		humidity := r.FormValue("humidity")
		time := time.Now()
		log.Println(id, temperature, humidity, time)
		r := new(result)
		r.id = id
		r.temperature = temperature
		r.humidity = humidity
		c <- *r

	}
}

func saveToInfluxDB(con client.Client, tags map[string]string, fields map[string]interface{}, c chan result, wg *sync.WaitGroup) {

	wg.Add(1)
	defer wg.Done()

	bp, _ := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  "sensor",
		Precision: "",
	})

	for r := range c {
		id := r.id
		temperature := r.temperature
		humidity := r.humidity
		fields["id"] = id
		fields["temperature"] = temperature
		fields["humidity"] = humidity
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
