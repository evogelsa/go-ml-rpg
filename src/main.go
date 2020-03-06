package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"

	"github.iu.edu/evogelsa/go-ml-rpg/game"
	"github.iu.edu/evogelsa/go-ml-rpg/web"
)

func main() {
	// define ai flag with default minmax
	aiAlg := flag.String("ai", "minmax",
		"Specifies algorithm AI will use. Options:"+
			"\n\trand\n\tminmax\n\treinforcement\n")

	// parse flags
	flag.Parse()

	// set algorithm accordingly to commandline flag
	switch *aiAlg {
	case "minmax":
		fmt.Println("Using AI MinMax")
		game.AI_ALG = game.AI_MINMAX
	case "rand":
		fmt.Println("Using AI Rand")
		game.AI_ALG = game.AI_RAND
	case "reinforcement":
		fmt.Println("Using AI Reinforcement")
		game.AI_ALG = game.AI_REINFORCEMENT
	default:
		fmt.Println("AI Unrecognized, run with flag -h for help")
		fmt.Println("Defaulting to AI MinMax")
		game.AI_ALG = game.AI_MINMAX
	}

	// seed random function with time
	rand.Seed(time.Now().Unix())
	// start webserver
	web.Server()
}
