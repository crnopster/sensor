package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"
)

type resultToSend struct {
	ID          string
	Temperature float32
	Humidity    float32
}

func worker(ctx context.Context, wg *sync.WaitGroup, c chan result) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		for result := range c {
			b := new(bytes.Buffer)
			r := &resultToSend{
				ID:          result.clientID,
				Temperature: result.temperature,
				Humidity:    result.humidity,
			}

			err := json.NewEncoder(b).Encode(r)
			if err != nil {
				log.Println(err.Error())
			}

			log.Println(b)

			resp, err := http.Post("http://localhost:8080/", "application/json;charset=utf-8", b)
			if err != nil {
				log.Fatal(err)
			}
			defer resp.Body.Close()
		}
	}
}
