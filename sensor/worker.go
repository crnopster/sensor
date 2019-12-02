package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	uuid "github.com/google/uuid"
)

type metricToSend struct {
	UUID      uuid.UUID `json:"uuid"`
	Type      string    `json:"type"`
	Data      float32   `json:"data"`
	Timestamp time.Time `json:"timestamp"`
}

func worker(ctx context.Context, wg *sync.WaitGroup, c chan metric) {
	log.Println("worker started")
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			log.Println("worker stopped")
			return
		default:
		}

		for result := range c {
			b := new(bytes.Buffer)
			r := &metricToSend{
				UUID:      result.ID,
				Type:      result.Type,
				Data:      result.Data,
				Timestamp: result.Timestamp,
			}

			err := json.NewEncoder(b).Encode(r)
			if err != nil {
				log.Println(err.Error())
			}

			resp, err := http.Post("http://localhost:8080/", "application/json", b)
			if err != nil {
				log.Fatal(err)
			}

			defer resp.Body.Close()
		}
	}
}
