package redissaver

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/crnopster/sensor/metric"
	"github.com/go-redis/redis/v7"
	jsoniter "github.com/json-iterator/go"
)

// RedisConn .
type RedisConn struct {
	Client *redis.Client
	C      chan metric.Metric
}

// Save metric to redisDB
func (r *RedisConn) Save(m metric.Metric) {
	r.C <- m
}

func (r *RedisConn) redisWorker(ctx context.Context, wg *sync.WaitGroup) {
	log.Println("redisWorker started")
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			log.Println("redisWorker stopped")
			return
		default:
			for m := range r.C {
				data := make(map[string]float32)
				data[m.Type] = m.Data
				uuid := fmt.Sprint(m.UUID)
				dataToSend, err := jsoniter.Marshal(data)
				if err != nil {
					log.Println(err.Error())
				}
				err = r.Client.Set(uuid, dataToSend, 10*time.Second).Err()
				if err != nil {
					log.Println(err.Error())
				}
				log.Println("sent to redis")
			}
		}
	}
}

// Worker calls redisWorker to save data into redisDB
func (r RedisConn) Worker(ctx context.Context, wg *sync.WaitGroup, workerCount int) {
	wg.Add(workerCount)
	for i := 0; i < workerCount; i++ {
		go r.redisWorker(ctx, wg)
	}

}

// NewRedisClient returns connect to redisDB & metric chan
func NewRedisClient() *RedisConn {
	return &RedisConn{
		Client: redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		}),
		C: make(chan metric.Metric, 1),
	}
}
