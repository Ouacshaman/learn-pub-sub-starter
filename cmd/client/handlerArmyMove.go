package main

import (
	"fmt"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
)

func handlerArmyMove(gs *gamelogic.GameState) func(gamelogic.ArmyMove) {
	return func(army_move gamelogic.ArmyMove) {
		defer fmt.Println("> ")
		gs.HandleMove(army_move)
	}
}
