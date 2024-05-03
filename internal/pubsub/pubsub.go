package pubsub

import (
	"encoding/json"
	"context"
	"time"
)

func PublishJSON[T any](ch *amqp.Channel, exchange, key string, val T) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	jsonVal, _ := json.Marshal(val)
	msg := ch.Publishing{
		DeliveryMode: amqp.Persistent,
		Timestamp:    time.Now(),
		ContentType:  "application/json",
		Body:         jsonVal,
	}

	err := ch.PublishWithContext(ctx, exchange, key, false, false,msg) 

	if err != nil {
		return err
	}
	return nil
}

