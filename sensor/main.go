package main

import (
	"context"
	"sync"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	wg := &sync.WaitGroup{}
	c := make(chan metric, 1)
	wg.Add(10)
	for i := 0; i < 5; i++ {
		go worker(ctx, wg, c)
		go temperature(ctx, wg, c)
	}
	wg.Wait()
}
