package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func shutdown(ctx context.Context, cancel context.CancelFunc, srv *http.Server, wg *sync.WaitGroup) {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	d := <-done

	log.Println("Server stopped with signal : ", d)
	cancel()
	time.Sleep(1 * time.Second)
	err := srv.Shutdown(ctx)
	if err != nil {
		log.Println(err.Error())
	}
	wg.Done()
}
