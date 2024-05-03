package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"syscall"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril server...")
	connectionString := "amqp://guest:guest@localhost:5672/"
	connect, err := amqp.Dial(connectionString)
	defer connect.Close()

	if err != nil {
		fmt.Println("Error received while connecting to RabbitMQ server ", err)
		return
	}
	fmt.Println("RabbitMQ Server started successfully")

	ch, _ := connect.Channel()
	defer ch.Close()
	data := routing.PlayingState{
		IsPaused: true,
	}
	err = pubsub.PublishJSON(ch, routing.ExchangePerilDirect, routing.PauseKey, data)
	if err != nil {
		log.Fatalf(err.Error())
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Blocking, press CTRL+C to continue")
	<-done
	fmt.Println("Received CTRL+C, closing the connection")

}
