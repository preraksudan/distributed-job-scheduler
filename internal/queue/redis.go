package queue

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

var Client *redis.Client

func ConnectRedis() {

	Client = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	_, err := Client.Ping(context.Background()).Result()

	if err != nil {
		log.Fatal("Unable to connect to Redis:", err)
	}

	log.Println("Connected to Redis")
}
