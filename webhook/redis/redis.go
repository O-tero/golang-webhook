package redis

import (
	"context"
	"encoding/json"
	"log"

	"github.com/go-redis/redis/v8"
)

// WebhookPayLoad defines the structure of the data expeted
// tobe recieved from Redis including UL webhook ID, and relevant data
type WebhookPayLoad struct {
	Url			string		`json:"url"`
	WebhookId	string		`json:"webhookId"`
	Data		struct	{
		Id 		string		`json:"id"`
		Payment	string		`json:"payment"`
		Event	string		`json:"event"`
		Date	string		`json:"date"`
	}	`json:"data"`
}

func Subscribe(ctx context.Context, client *redis.Client, webhookQueue chan WebhookPayLoad) error {
	// Subscribe to the "webhooks" channel in Redis
	pubSub := client.Subscribe(ctx, "payments")

	// Ensure that the PubSub connecton is closed when the function exits
	defer func(pubSub *redis.PubSub) {
		if err := pubSub.Close(); err != nil {
			log.Println("Error closing PubSub:", err)
		}
	}(pubSub)

	var payload WebhookPayLoad

	// Infinite lop to continously recieve messages from the "webhook" channel
	for {
		// Recieve the message
		msg, err := pubSub.ReceiveMessage(ctx)
		if err != nil {
			return err // Return the error if there's an issue receiving the message 
		}

		// Unmarshal the JSON string into a WebhookPayLoad struct
		err = json.Unmarshal([]byte(msg.Payload), &payload)
		if err != nil {
			log.Println("Error unmarshalling payload:", err)
			continue // Continue with the next message if there's an error unmarshalling
		}

		// Send the webhook
		webhookQueue <- payload //Sending the payload to the channel
		}
}
