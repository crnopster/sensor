package main

import (
	"context"
	"math/rand"
	"sync"
	"time"

	uuid "github.com/google/uuid"
)

var defaultTemperature float32 = 20
var defaultHumidity float32 = 50
var defaultO2 float32 = 25

type metric struct {
	ID        uuid.UUID
	Data      float32
	Type      string
	Timestamp time.Time
}

func temperature(ctx context.Context, wg *sync.WaitGroup, c chan metric) {
	t := "temperature"
	temp := defaultTemperature
	ID := uuid.New()
	defer wg.Done()

	for {
		hourTemperature := float32(rand.Intn(41))
		step := (hourTemperature - temp) / 360
		for i := 0; i < 360; i++ {
			select {
			case <-ctx.Done():
				time.Sleep(time.Second)
				return
			default:
				temp = temp + step
				result := &metric{
					ID:        ID,
					Data:      temp,
					Type:      t,
					Timestamp: time.Now(),
				}
				c <- *result
				time.Sleep(time.Second * 10)
			}
		}
	}
}

func humidity(ctx context.Context, wg *sync.WaitGroup, c chan metric) {
	t := "humidity"
	hum := defaultHumidity
	ID := uuid.New()
	defer wg.Done()

	for {
		hourHumidity := float32(rand.Intn(71) + 20)
		step := (hourHumidity - hum) / 360
		for i := 0; i < 360; i++ {
			select {
			case <-ctx.Done():
				time.Sleep(time.Second)
				return
			default:
				hum = hum + step
				result := &metric{
					ID:        ID,
					Data:      hum,
					Type:      t,
					Timestamp: time.Now(),
				}
				c <- *result
				time.Sleep(time.Second * 10)
			}
		}
	}
}

func oxygen(ctx context.Context, wg *sync.WaitGroup, c chan metric) {
	t := "oxygen"
	O2 := defaultO2
	ID := uuid.New()
	defer wg.Done()

	for {
		hourO2 := float32(rand.Intn(5) + 20)
		step := (hourO2 - O2) / 360
		for i := 0; i < 360; i++ {
			select {
			case <-ctx.Done():
				time.Sleep(time.Second)
				return
			default:
				O2 = O2 + step
				result := &metric{
					ID:        ID,
					Data:      O2,
					Type:      t,
					Timestamp: time.Now(),
				}
				c <- *result
				time.Sleep(time.Second * 10)

			}
		}
	}
}
