package main

import "context"

import "sync"

func start(ctx context.Context, wg *sync.WaitGroup, c chan metric, temperatureSensorCount, humiditySensorCount, oxygenSensorCount, workerCount int) {
	wg.Add(temperatureSensorCount + humiditySensorCount + oxygenSensorCount + workerCount)

	for i := 0; i < workerCount; i++ {
		go worker(ctx, wg, c)
	}

	for i := 0; i < temperatureSensorCount; i++ {
		go temperature(ctx, wg, c)
	}
	for i := 0; i < humiditySensorCount; i++ {
		go humidity(ctx, wg, c)
	}
	for i := 0; i < oxygenSensorCount; i++ {
		go oxygen(ctx, wg, c)
	}
}
