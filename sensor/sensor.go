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
	temperature string
	humidity    string
}

func sensor(ctx context.Context, wg *sync.WaitGroup, c chan result) {
	temperature := defaulttemperature
	humidity := defaulthumidity

	wg.Add(1)
	defer wg.Done()
	clientID := fmt.Sprint(uuid.New())
	r := new(result)
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
			r.clientID = clientID
			r.temperature = fmt.Sprintf("%v", temperature)
			r.humidity = fmt.Sprintf("%v", humidity)
			c <- *r
			time.Sleep(time.Second * 10)
		}
	}
}
