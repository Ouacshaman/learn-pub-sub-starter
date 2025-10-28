package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril server...")

	conn_str := "amqp://guest:guest@localhost:5672/"

	conn, err := amqp.Dial(conn_str)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer conn.Close()

	fmt.Println("Connection is sucessfully established")

	ch, err := conn.Channel()
	if err != nil {
		fmt.Println(err)
		return
	}

	err = ch.Confirm(false)
	if err != nil {
		fmt.Println(err)
		return
	}

	play_state := routing.PlayingState{IsPaused: true}

	err = pubsub.PublishJSON(ch, routing.ExchangePerilDirect, routing.PauseKey, play_state)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan

	fmt.Println("Connection Closing")
}
