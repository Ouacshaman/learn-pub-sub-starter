package main

import (
	"fmt"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/rabbitmq/amqp091-go"
)

func handlerWar(ch *amqp091.Channel, gs *gamelogic.GameState) func(gamelogic.RecognitionOfWar) pubsub.Acktype {
	return func(war gamelogic.RecognitionOfWar) pubsub.Acktype {
		defer fmt.Println("> ")
		out, winner, loser := gs.HandleWar(war)
		switch out {
		case gamelogic.WarOutcomeNotInvolved:
			return pubsub.NackRequeue
		case gamelogic.WarOutcomeNoUnits:
			return pubsub.NackDiscard
		case gamelogic.WarOutcomeOpponentWon:
			message := fmt.Sprintf("%s won a war against %s\n", winner, loser)
			err := publishGameLog(ch, message, gs.GetUsername())
			if err != nil {
				fmt.Printf("error: %s\n", err)
				return pubsub.NackRequeue
			}
			return pubsub.Ack
		case gamelogic.WarOutcomeYouWon:
			message := fmt.Sprintf("%s won a war against %s\n", winner, loser)
			err := publishGameLog(ch, message, gs.GetUsername())
			if err != nil {
				fmt.Printf("error: %s\n", err)
				return pubsub.NackRequeue
			}
			return pubsub.Ack
		case gamelogic.WarOutcomeDraw:
			message := fmt.Sprintf("A war between %s and %s resulted in a draw\n", winner, loser)
			err := publishGameLog(ch, message, gs.GetUsername())
			if err != nil {
				return pubsub.NackRequeue
			}
			return pubsub.Ack
		default:
			err := fmt.Errorf("Outcome Not Found")
			fmt.Println(err)
			return pubsub.NackDiscard
		}
	}
}
