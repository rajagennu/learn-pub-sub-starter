package pubsub

import (
	"context"
	"encoding/json"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishJSON[T any](ch *amqp.Channel, exchange, key string, val T) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	jsonVal, _ := json.Marshal(val)
	msg := amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		Timestamp:    time.Now(),
		ContentType:  "application/json",
		Body:         jsonVal,
	}

	err := ch.PublishWithContext(ctx, exchange, key, false, false, msg)

	if err != nil {
		return err
	}
	return nil
}

func DeclareAndBind(
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	simpleQueueType int, // an enum to represent "durable" or "transient"
) (*amqp.Channel, amqp.Queue, error) {
	ch, _ := conn.Channel()
	durable := false

	if simpleQueueType == 1 {
		durable = true
	}

	autoDelete := false
	exclusive := false
	if simpleQueueType == 2 {
		autoDelete = true
		exclusive = true
	}
	noWait := false

	queue, _ := ch.QueueDeclare(queueName, durable, autoDelete, exclusive, noWait, nil)
	err := ch.QueueBind(queue.Name, key, exchange, noWait, nil)
	return ch, queue, err

}
