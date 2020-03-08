package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.iu.edu/evogelsa/go-ml-rpg/game"
	"github.iu.edu/evogelsa/go-ml-rpg/web"
)

const PORT = ":42069"

func main() {
	// define ai flag with default minmax
	aiAlg := flag.String("ai", "reinforcement",
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
		fmt.Println("Defaulting to AI Reinforcement")
		game.AI_ALG = game.AI_MINMAX
	}

	// check if directory structure exists and make where absent
	_, err := os.Stat(web.SAVE_DIR)
	if err != nil {
		if os.IsNotExist(err) {
			err := os.Mkdir(web.SAVE_DIR, 0700)
			if err != nil {
				panic(err)
			}
		}
	}
	_, err = os.Stat(web.LOG_DIR)
	if err != nil {
		if os.IsNotExist(err) {
			err := os.Mkdir(web.LOG_DIR, 0700)
			if err != nil {
				panic(err)
			}
		}
	}
	_, err = os.Stat(web.CHAR_DIR)
	if err != nil {
		if os.IsNotExist(err) {
			err := os.Mkdir(web.CHAR_DIR, 0700)
			if err != nil {
				panic(err)
			}
		}
	}
	_, err = os.Stat(web.IMG_DIR)
	if err != nil {
		if os.IsNotExist(err) {
			err := os.Mkdir(web.IMG_DIR, 0700)
			if err != nil {
				panic(err)
			}
		}
	}

	// seed random function with time
	rand.Seed(time.Now().Unix())
	// start webserver
	web.Server(PORT)
}
