package main

import (
	"log"
	"sync"
	"time"

	"github.com/go-redis/redis/v7"
	jsoniter "github.com/json-iterator/go"
)

func saveToRedisDB(clientRedis redis.Client, wg *sync.WaitGroup, cr chan result) {
	defer wg.Done()
	m := make(map[float32]float32)
	for res := range cr {
		m[res.Temperature] = res.Humidity
		b, err := jsoniter.Marshal(m)
		if err != nil {
			log.Println(err.Error())
		}
		m[res.Temperature] = res.Humidity
		cmd := clientRedis.Set(res.ID, b, 10*time.Second)
		log.Println(cmd)
		err = cmd.Err()
		if err != nil {
			log.Println(err)
		}
	}
}
