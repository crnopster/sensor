package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	uuid "github.com/google/uuid"
)

var defaulttemperature float32 = 20
var defaulthumidity float32 = 50

type result struct {
	clientID    string
	temperature float32
	humidity    float32
}

func sensor(ctx context.Context, wg *sync.WaitGroup, c chan result) {
	temperature := defaulttemperature
	humidity := defaulthumidity

	defer wg.Done()
	clientID := fmt.Sprint(uuid.New())

	for {
		hourTemperature := float32(rand.Intn(41))
		hourHimidity := float32(rand.Intn(70) + 20)
		x := (hourTemperature - temperature) / 360
		y := (hourHimidity - humidity) / 360
		for i := 0; i < 360; i++ {
			select {
			case <-ctx.Done():
				time.Sleep(time.Millisecond * 100)
				return
			default:
			}
			temperature = temperature + x
			humidity = humidity + y
			r := &result{
				clientID:    clientID,
				temperature: temperature,
				humidity:    humidity,
			}
			c <- *r
			time.Sleep(time.Second * 10)
		}
	}
}
