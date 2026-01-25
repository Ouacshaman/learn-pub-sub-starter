package main

import (
	"fmt"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
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

	gamelogic.PrintServerHelp()

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

	_, _, err = pubsub.DeclareAndBind(conn, routing.ExchangePerilTopic, routing.GameLogSlug, "game_logs.*", pubsub.Durable)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = pubsub.SubscribeGob(conn, routing.ExchangePerilTopic, routing.GameLogSlug, routing.GameLogSlug+".*", pubsub.Durable, handlerLog())
	if err != nil {
		fmt.Println(err)
		return
	}

	for true {
		input := gamelogic.GetInput()
		if len(input) == 0 {
			continue
		}
		switch input[0] {
		case "pause":
			fmt.Println("Publishing paused game state")
			play_state := routing.PlayingState{IsPaused: true}
			err = pubsub.PublishJSON(ch, routing.ExchangePerilDirect, routing.PauseKey, play_state)
			if err != nil {
				fmt.Println(err)
				continue
			}

		case "resume":
			fmt.Println("Publishing resume game state")
			play_state := routing.PlayingState{IsPaused: false}
			err = pubsub.PublishJSON(ch, routing.ExchangePerilDirect, routing.PauseKey, play_state)
			if err != nil {
				fmt.Println(err)
				continue
			}

		case "quit":
			fmt.Println("Logging Out")
			return

		default:
			fmt.Println("I don't understand your entry")
		}
	}

}
