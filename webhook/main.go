package main

import (
	"context"
	"log"
	"os"

	redisClient "webhook/redis"

	"webhook/queue"

	"github.com/go-redis/redis/v8"
)

func main() {
	// Create  context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize the redis client
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: "",
		DB:		  0,
	})

	// Create a channel to act as the queue
	webhookQueue := make(chan redisClient.WebhookPayLoad, 100)

	go queue.ProcessWebhooks(ctx, webhookQueue)

	// Subscribe to the "transactions" channel
	err := redisClient.Subscribe(ctx, client, webhookQueue)

	if err != nil {
		log.Println("Error:", err)
	}

	select {}

}
