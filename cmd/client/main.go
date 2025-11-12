package main

import (
	"fmt"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
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

	fmt.Println(username)
}
