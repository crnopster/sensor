package main

import (
	"context"
	"flag"
	"sync"
)

func main() {
	tSensorCount := flag.Int("tSensor", 10, "Number of temperature sensors")
	hSensorCount := flag.Int("hSensor", 10, "Number of humidity sensors")
	oSensorCount := flag.Int("oSensor", 10, "Number of oxygen sensors")
	workerCount := flag.Int("worker", 1, "Number of workers")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg := &sync.WaitGroup{}

	c := make(chan metric, 1)

	start(ctx, wg, c, *tSensorCount, *hSensorCount, *oSensorCount, *workerCount)

	wg.Wait()
}
