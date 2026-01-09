package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril client...")

	conn_str := "amqp://guest:guest@localhost:5672/"

	conn, err := amqp.Dial(conn_str)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer conn.Close()

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		fmt.Println(err)
		return
	}

	n_username := "pause." + username

	_, _, err = pubsub.DeclareAndBind(conn, routing.ExchangePerilDirect, n_username, routing.PauseKey, pubsub.Transient)
	if err != nil {
		fmt.Println(err)
		return
	}

	gamestate := gamelogic.NewGameState(username)

	err = pubsub.SubscribeJSON(conn, routing.ExchangePerilDirect, n_username, routing.PauseKey, pubsub.Transient, handlerPause(gamestate))
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
		case "spawn":
			err := gamestate.CommandSpawn(input)
			if err != nil {
				fmt.Println(err)
				continue
			}
		case "move":
			_, err := gamestate.CommandMove(input)
			if err != nil {
				fmt.Println(err)
				continue
			}
		case "status":
			gamestate.CommandStatus()
		case "help":
			gamelogic.PrintClientHelp()
		case "spam":
			fmt.Println("Spamming not allowed yet!")
		case "quit":
			gamelogic.PrintQuit()
			return

		default:
			fmt.Println("Invalid Command and I don't understand your entry")
			continue
		}
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan

	fmt.Println("Connection Closing")

}
