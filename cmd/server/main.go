package main

import ( 
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"os/signal"
	"syscall"
	"os"
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
	defer ch.close()
        


	done := make(chan os.Signal, 1)
	signal.Notify( done, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("Blocking, press CTRL+C to continue")
	<- done
	fmt.Println("Received CTRL+C, closing the connection")

}


