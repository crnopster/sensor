package main

import (
	"log"
	"sync"
	"time"

	client "github.com/influxdata/influxdb1-client/v2"
)

func saveToInfluxDB(clientInflux client.Client, tags map[string]string, ci chan result, wg *sync.WaitGroup) {
	fields := make(map[string]interface{})

	defer wg.Done()

	bp, _ := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  "sensor",
		Precision: "",
	})

	for res := range ci {
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
		err = clientInflux.Write(bp)
		if err != nil {
			log.Println(err.Error())
		}
	}
}
