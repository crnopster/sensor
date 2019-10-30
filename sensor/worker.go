package main

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"sync"
)

func worker(ctx context.Context, wg *sync.WaitGroup, c chan result) {
	wg.Add(1)
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		for result := range c {
			clientID := result.clientID
			temperature := result.temperature
			humidity := result.humidity

			resp, err := http.PostForm("http://localhost:8080/",
				url.Values{"id": {clientID}, "temperature": {temperature}, "humidity": {humidity}})
			if err != nil {
				log.Fatal(err)
			}
			resp.Body.Close()
		}
	}
}
