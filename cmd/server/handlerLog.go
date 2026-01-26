package main

import (
	"fmt"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
)

func handlerLog() func(gamelog routing.GameLog) pubsub.Acktype {
	return func(gamelog routing.GameLog) pubsub.Acktype {
		defer fmt.Println("> ")
		err := gamelogic.WriteLog(gamelog)
		if err != nil {
			fmt.Println("WriteLog error:", err)
			return pubsub.NackRequeue
		}
		fmt.Println("WriteLog OK")
		return pubsub.Ack
	}
}
