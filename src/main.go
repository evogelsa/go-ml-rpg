package main

import (
	"math/rand"
	"time"

	"github.iu.edu/evogelsa/go-ml-rpg/game"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	k := game.NewKnight("Ethan")
	w := game.NewWizard("Renan")

	game.PrintCharacter(k)
	game.PrintCharacter(w)

	for {
		m1 := game.SelectMove(k)
		m2 := game.SelectMove(w)
		game.Turn(&k, &w, m1, m2)
	}
}
