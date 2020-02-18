package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.iu.edu/evogelsa/go-ml-rpg/game"
	"github.iu.edu/evogelsa/go-ml-rpg/ui"

	"github.com/gorilla/mux"
)

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `<h1>Hello, World!</h1>`)
}

func newRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", home)
	r.HandleFunc("/uitest", ui.Load)

	return r
}

func getIntInput() int {
	reader := bufio.NewReader(os.Stdin)
	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		choice, err := strconv.Atoi(input[:len(input)-1])
		if err == nil && choice >= 0 && choice <= 5 {
			return choice
		} else {
			fmt.Printf("Please input a valid integer selection\n")
		}
	}
}

func selectMove(p game.Class) game.Move {
	fmt.Printf("%s select move:\n", p.PlayerName)
	fmt.Printf("\t(0)Heavy attack\n\t(1)Quick attack\n\t(2)Standard attack\n")
	fmt.Printf("\t(3)Block\n\t(4)Parry\n\t(5)Evade\n")
	return game.Move(getIntInput())
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	r := newRouter()

	k := game.NewKnight("Ethan")
	w := game.NewWizard("Renan")

	game.PrintCharacter(k)
	game.PrintCharacter(w)

	for {
		m1 := selectMove(k)
		m2 := selectMove(w)
		game.Turn(&k, &w, m1, m2)
	}

	log.Fatal(http.ListenAndServe(":8080", r))
}
