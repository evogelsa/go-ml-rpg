package web

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.iu.edu/evogelsa/go-ml-rpg/game"

	"github.com/gorilla/mux"
)

const (
	FILE_DIR = "./web/assets/"
	SAVE_DIR = "./web/assets/saves/"
	PORT     = ":8080"
)

func generateChar(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	class := vars["class"]
	name := vars["name"]

	var char game.Class
	switch class {
	case "knight":
		char = game.NewKnight(name)
	case "archer":
		char = game.NewArcher(name)
	case "wizard":
		char = game.NewWizard(name)
	default:
		panic("err creating char")
	}

	err := writeCharToFile(char)
	if err != nil {
		panic(err)
	}
}

func fileToString(fn string) (string, error) {
	fn = FILE_DIR + fn

	body, err := ioutil.ReadFile(fn)
	if err != nil {
		return "", err
	}

	text := string(body)
	return text, nil
}

func readCharFromFile(fn string) (string, error) {
	fn = SAVE_DIR + fn

	body, err := ioutil.ReadFile(fn)
	if err != nil {
		return "", err
	}

	text := string(body)
	return text, nil
}

func writeCharToFile(c game.Class) error {
	//get file from fn (overwrites if file exists)
	fn := SAVE_DIR + c.PlayerName + "." + c.ClassName
	f, err := os.Create(fn)
	if err != nil {
		return err
	}

	html, err := fileToString("charTable.html")
	if err != nil {
		return err
	}

	s := fmt.Sprintf(
		html,
		c.PlayerName, c.ClassName,
		"Health", c.Health,
		"Stamina", c.Stamina,
		"Armor", c.Armor,
		"Strength", c.Strength,
		"Dexterity", c.Dexterity,
		"Intellect", c.Intellect,
	)

	fmt.Fprint(f, s)

	err = f.Close()
	if err != nil {
		return err
	}

	return nil
}

func fightScreen(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cn1 := vars["char1"]
	cn2 := vars["char2"]

	screen, err := fileToString("fightScreen.html")
	if err != nil {
		panic(err)
	}

	char1, err := readCharFromFile(cn1)
	if err != nil {
		panic(err)
	}
	char2, err := readCharFromFile(cn2)
	if err != nil {
		panic(err)
	}

	fmt.Fprintf(w, screen, char1, char2, "holder", "holder")
}

func newRouter() *mux.Router {
	r := mux.NewRouter()

	fileServer := http.FileServer(http.Dir(FILE_DIR))
	handler := http.StripPrefix("/assets/", fileServer)
	r.PathPrefix("/assets/").Handler(handler)

	r.HandleFunc("/newChar/{class}/{name}", generateChar)
	r.HandleFunc("/game/{char1}/{char2}", fightScreen)

	return r
}

func Server() {
	r := newRouter()

	log.Fatal(http.ListenAndServe(PORT, r))
}
