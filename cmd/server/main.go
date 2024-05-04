package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"syscall"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
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

	var simpleQueueType pubsub.SimpleQueueType = pubsub.Durable
	pubsub.DeclareAndBind(connect, routing.ExchangePerilTopic, routing.GameLogSlug, routing.GameLogSlug+"1", simpleQueueType.EnumIndex())

	gamelogic.PrintServerHelp()
	for {
		data := routing.PlayingState{
			IsPaused: true,
		}
		words := gamelogic.GetInput()
		if len(words) == 0 {
			continue
		}
		if words[0] == "pause" {
			log.Println("sending the pause message")
			data = routing.PlayingState{
				IsPaused: true,
			}

			err = pubsub.PublishJSON(ch, routing.ExchangePerilDirect, routing.PauseKey, data)
			if err != nil {
				log.Fatalf(err.Error())
			}
			continue
		}

		if words[0] == "resume" {
			log.Println("sending the resume message")
			data = routing.PlayingState{
				IsPaused: false,
			}

			err = pubsub.PublishJSON(ch, routing.ExchangePerilDirect, routing.PauseKey, data)
			if err != nil {
				log.Fatalf(err.Error())
			}
			continue
		}

		if words[0] == "quit" {
			log.Println("Exiting ")
			break
		}
		log.Println(" I do not understand the command")

	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Blocking, press CTRL+C to continue")
	<-done
	fmt.Println("Received CTRL+C, closing the connection")

}
