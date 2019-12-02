package influxsaver

import (
	"context"
	"log"
	"sync"

	"github.com/crnopster/sensor/metric"
	influx "github.com/influxdata/influxdb1-client/v2"
)

// InfluxConn .
type InfluxConn struct {
	Client influx.Client
	C      chan metric.Metric
}

// Save metric to influxDB
func (i *InfluxConn) Save(m metric.Metric) {
	i.C <- m
}

func (i *InfluxConn) influxWorker(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Println("influxWorker started")

	bp, _ := influx.NewBatchPoints(influx.BatchPointsConfig{
		Database:  "sensor",
		Precision: "",
	})

	for {
		fields := make(map[string]interface{})
		select {
		case <-ctx.Done():
			log.Println("influxWorker stopped")
			return
		case m := <-i.C:
			fields["UUID"] = m.UUID
			fields["Type"] = m.Type
			fields["Data"] = m.Data
			newPoint, err := influx.NewPoint(
				"NEWSENSOR",
				nil,
				fields,
				m.Timestamp,
			)
			if err != nil {
				log.Println(err.Error())
			}
			bp.AddPoint(newPoint)
			err = i.Client.Write(bp)
			if err != nil {
				log.Println(err.Error())
			}
			log.Println("sent to influx")
		}
	}
}

// Worker calls influxWorker to save data into influxDB
func (i InfluxConn) Worker(ctx context.Context, wg *sync.WaitGroup, workerCount int) {
	wg.Add(workerCount)
	for a := 0; a < workerCount; a++ {
		go i.influxWorker(ctx, wg)
	}
}

// NewInfluxClient returns connect to influxDB & metric chan
func NewInfluxClient() *InfluxConn {
	cli, err := influx.NewHTTPClient(influx.HTTPConfig{
		Addr:     "http://localhost:8086",
		Username: "sensors",
		Password: "pass",
	})
	if err != nil {
		log.Println(err.Error())
	}
	return &InfluxConn{
		Client: cli,
		C:      make(chan metric.Metric, 1),
	}
}
