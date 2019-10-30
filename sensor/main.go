package main

import (
	"context"
	"flag"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	ctx := context.Background()
	wg := &sync.WaitGroup{}

	sensorCount := flag.Int("sensorCount", 100, "sensor count")
	workerCount := flag.Int("workerCount", 5, "worker count")
	flag.Parse()

	c := make(chan result)

	for a := 0; a < *workerCount; a++ {
		go worker(ctx, wg, c)
	}
	for i := 0; i < *sensorCount; i++ {
		go sensor(ctx, wg, c)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	d := <-done
	ctx.Done()
	log.Println("Sensor emulator stopped. Signal: ", d)
	for i := 0; i < *sensorCount+*workerCount; i++ {
		wg.Done()
	}
	wg.Wait()
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
