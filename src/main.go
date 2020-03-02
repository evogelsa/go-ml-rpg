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
	// define ai flag
	aiAlg := flag.String("ai", "minmax",
		"Specifies algorithm AI will use. Options:"+
			"\n\trand\n\tminmax\n\tneural\n")

	// parse flags
	flag.Parse()

	switch *aiAlg {
	case "rand":
		fmt.Println("Using AI Rand")
		game.AI_ALG = game.AI_RAND
	case "minmax":
		fmt.Println("Using AI MinMax")
		game.AI_ALG = game.AI_MINMAX
	case "neural":
		fmt.Println("Using AI Neural")
		game.AI_ALG = game.AI_NEURAL
	default:
		fmt.Println("AI Unrecognized, run with flag -h for help")
		fmt.Println("Defaulting to AI MinMax")
		game.AI_ALG = game.AI_MINMAX
	}

	rand.Seed(time.Now().Unix())
	web.Server()
}
