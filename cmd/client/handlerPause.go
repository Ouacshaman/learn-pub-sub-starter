package main

import (
	"fmt"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
)

func handlerPause(gs *gamelogic.GameState) func(routing.PlayingState) {
	state := routing.PlayingState{
		IsPaused: gs.Paused,
	}
	return func(playstate routing.PlayingState) {
		defer fmt.Println(">")
		gs.HandlePause(state)
	}
}
