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
	fmt.Println("Starting Peril client...")
	connectionString := "amqp://guest:guest@localhost:5672/"
	connect, err := amqp.Dial(connectionString)
	defer connect.Close()

	if err != nil {
		fmt.Println("Error received while connecting to RabbitMQ server ", err)
		return
	}
	fmt.Println("RabbitMQ Server started successfully")
	username, _ := gamelogic.ClientWelcome()
	var simpleQueueType pubsub.SimpleQueueType = pubsub.Transient
	pubsub.DeclareAndBind(connect, routing.ExchangePerilDirect, routing.PauseKey+"."+username, routing.PauseKey, simpleQueueType.EnumIndex())

	// gamestate
	gameState := gamelogic.NewGameState(username)

	for {
		command := gamelogic.GetInput()
		if command[0] == "spawn" {
			err := gameState.CommandSpawn(command)
			if err != nil {
				log.Println(err.Error())
			}
			continue
		}

		if command[0] == "move" {
			armyMove, err := gameState.CommandMove(command)
			if err != nil {
				log.Println(err.Error())
			}
			log.Printf("Army moved successfully to %s", armyMove.ToLocation)
			continue
		}

		if command[0] == "status" {
			gameState.CommandStatus()
			continue
		}

		if command[0] == "help" {
			gamelogic.PrintClientHelp()
			continue
		}

		if command[0] == "spam" {
			log.Println("Spamming not allowed yet!")
			continue
		}

		if command[0] == "quit" {
			gamelogic.PrintQuit()
			continue
		}
		log.Println("Error!, invalid choice " + command[0])

	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Blocking, press CTRL+C to continue")
	<-done
	fmt.Println("Received CTRL+C, closing the connection")

}
