package main

import (
	"fmt"
	"time"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
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
			return publishGameLog(ch, message, gs.GetUsername())
		case gamelogic.WarOutcomeYouWon:
			message := fmt.Sprintf("%s won a war against %s\n", winner, loser)
			return publishGameLog(ch, message, gs.GetUsername())
		case gamelogic.WarOutcomeDraw:
			message := fmt.Sprintf("A war between %s and %s resulted in a draw\n", winner, loser)
			return publishGameLog(ch, message, gs.GetUsername())

		default:
			err := fmt.Errorf("Outcome Not Found")
			fmt.Println(err)
			return pubsub.NackDiscard
		}
	}
}

func publishGameLog(ch *amqp091.Channel, message, username string) pubsub.Acktype {
	gs := routing.GameLog{
		CurrentTime: time.Now(),
		Message:     message,
		Username:    username,
	}

	err := pubsub.PublishGob(ch, routing.ExchangePerilTopic,
		routing.GameLogSlug+"."+username,
		gs)
	if err != nil {
		fmt.Println(err)
		return pubsub.NackRequeue
	}
	return pubsub.Ack
}
