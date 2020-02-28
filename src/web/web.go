package web

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.iu.edu/evogelsa/go-ml-rpg/game"

	"github.com/gorilla/mux"
)

const (
	FILE_DIR = "./web/assets/"
	SAVE_DIR = "./saves/"
	PORT     = ":8080"
)

// fileToString takes in file name and convert to string
func fileToString(fn string) (string, error) {
	fn = FILE_DIR + fn

	body, err := ioutil.ReadFile(fn)
	if err != nil {
		return "", err
	}

	text := string(body)
	return text, nil
}

func logFile(fn, logStr string) error {
	fn = FILE_DIR + fn

	var text string

	_, err := os.Stat(fn)
	if err == nil {
		body, _ := ioutil.ReadFile(fn)
		text = string(body)
	}

	f, err := os.Create(fn)
	if err != nil {
		return err
	}

	text = logStr + text

	fmt.Fprint(f, text)

	return nil
}

// divWrap takes a string in and replaces newlines with html divs
func divWrap(s string) string {
	lines := strings.SplitN(s, "\n", -1)
	for i, line := range lines {
		lines[i] = "<div>" + line + "</div>"
	}
	lines = lines[:len(lines)-1]

	var res string
	for _, line := range lines {
		res += line
	}

	return res
}

// getMoves takes in character and returns string containing
// html table formatted moveset
func getMoves(char game.Class) string {
	s, err := fileToString("moveTable.html")
	if err != nil {
		panic(err)
	}

	var moves []interface{}
	switch char.ClassName {
	case "Knight":
		moves = []interface{}{
			"Crushing Blow",
			"Quick Thrust",
			"Sword Slash",
			"Shield",
			"Counter",
			"Drink Potion",
		}
	case "Archer":
		moves = []interface{}{
			"Piercing Shot",
			"Quick Fire",
			"Long Shot",
			"Block",
			"Dagger",
			"Apply Bandaid",
		}
	case "Wizard":
		moves = []interface{}{
			"Lightning",
			"Arcane Bolt",
			"Fireball",
			"Magic Shield",
			"Counterspell",
			"Heal",
		}
	}

	moveFmt := fmt.Sprintf(s, moves...)

	return moveFmt
}

// readCharFromFile takes in filename and turns character save
// into character class struct
func readCharFromFile(fn string) (game.Class, error) {
	char := game.Class{}

	fn = SAVE_DIR + fn
	f, err := os.Open(fn)
	if err != nil {
		return char, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	scanner.Scan()
	char.PlayerName = scanner.Text()

	scanner.Scan()
	char.ClassName = scanner.Text()

	scanner.Scan()
	char.Health, err = strconv.Atoi(scanner.Text())
	if err != nil {
		return char, err
	}

	scanner.Scan()
	char.Stamina, err = strconv.Atoi(scanner.Text())
	if err != nil {
		return char, err
	}

	scanner.Scan()
	char.Armor, err = strconv.Atoi(scanner.Text())
	if err != nil {
		return char, err
	}

	scanner.Scan()
	strength, err := strconv.ParseFloat(scanner.Text(), 32)
	char.Strength = float32(strength)
	if err != nil {
		return char, err
	}

	scanner.Scan()
	dexterity, err := strconv.ParseFloat(scanner.Text(), 32)
	char.Dexterity = float32(dexterity)
	if err != nil {
		return char, err
	}

	scanner.Scan()
	intellect, err := strconv.ParseFloat(scanner.Text(), 32)
	char.Intellect = float32(intellect)
	if err != nil {
		return char, err
	}

	if err := scanner.Err(); err != nil {
		return char, err
	}

	return char, nil
}

// writeCharToFile takes in character class and writes character info
// to a file called CharName.CharClass
func writeCharToFile(c game.Class) error {
	//get file from fn (overwrites if file exists)
	fn := SAVE_DIR + c.PlayerName + "." + c.ClassName
	f, err := os.Create(fn)
	if err != nil {
		return err
	}

	fmt.Fprint(
		f,
		c.PlayerName+"\n",
		c.ClassName+"\n",
		fmt.Sprintf("%d", c.Health)+"\n",
		fmt.Sprintf("%d", c.Stamina)+"\n",
		fmt.Sprintf("%d", c.Armor)+"\n",
		fmt.Sprintf("%f", c.Strength)+"\n",
		fmt.Sprintf("%f", c.Dexterity)+"\n",
		fmt.Sprintf("%f", c.Intellect)+"\n",
	)

	return nil
}

// charToHTML takes in character and reads its info in and returns
// an html table formatted string
func charToHTML(c game.Class) (string, error) {
	html, err := fileToString("charTable.html")
	if err != nil {
		return "", err
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

	return s, nil
}

// generateChar takes in a class and name and calls game to
// generate the char and writes to file
func generateChar(class, name string) error {
	var char game.Class
	switch class {
	case "Knight":
		char = game.NewKnight(name)
	case "Archer":
		char = game.NewArcher(name)
	case "Wizard":
		char = game.NewWizard(name)
	default:
		return errors.New("Could not parse class type")
	}

	err := writeCharToFile(char)
	if err != nil {
		return err
	}

	return nil
}

// parseMoveForm processes which move to execute and calls
// backend in game
func parseMoveForm(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	char1Name := vars["char1"]
	char2Name := vars["char2"]
	moveName := vars["move"]

	c1, err := readCharFromFile(char1Name)
	if err != nil {
		panic(err)
	}
	c2, err := readCharFromFile(char2Name)
	if err != nil {
		panic(err)
	}

	var move game.Move
	switch moveName {
	case "Heavy":
		move = game.HEAVY
	case "Quick":
		move = game.QUICK
	case "Standard":
		move = game.STANDARD
	case "Block":
		move = game.BLOCK
	case "Parry":
		move = game.PARRY
	case "Evade":
		move = game.EVADE
	}

	// process turn and get result
	outStr := game.Turn(&c1, &c2, move)
	// divwrap result
	res := divWrap(outStr)
	// get log file name
	logName := fmt.Sprintf("%s%s%s%s.log",
		c1.PlayerName, c1.ClassName, c2.PlayerName, c2.ClassName)
	// add result to log file
	err = logFile(logName, res)
	if err != nil {
		panic(err)
	}

	err = writeCharToFile(c1)
	if err != nil {
		panic(err)
	}
	err = writeCharToFile(c2)
	if err != nil {
		panic(err)
	}

	redirect := "/game/" + char1Name + "/" + char2Name
	http.Redirect(w, r, redirect, http.StatusFound)
}

func parseNewCharForm(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	class := r.Form.Get("class")
	name := r.Form.Get("name")

	err = generateChar(class, name)
	if err != nil {
		panic(err)
	}

	http.Redirect(w, r, "/selectChar", http.StatusFound)
}

// newCharacterScreen displays a screen to create a new character
func newCharacterScreen(w http.ResponseWriter, r *http.Request) {
	s, err := fileToString("newCharacterScreen.html")
	if err != nil {
		panic(err)
	}

	style, err := fileToString("styleHead.html")
	if err != nil {
		panic(err)
	}

	s = style + s

	fmt.Fprint(w, s)
}

// characterSelectScreen displays all characters in save dir with
// options to select each char
func characterSelectScreen(w http.ResponseWriter, r *http.Request) {
	charFiles, err := ioutil.ReadDir(SAVE_DIR)
	if err != nil {
		if os.IsNotExist(err) {
			err := os.Mkdir(SAVE_DIR, os.ModeDir)
			if err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}

	var chars [][]string
	for _, charFile := range charFiles {
		chars = append(chars, strings.SplitN(charFile.Name(), ".", -1))
	}

	var players [][]string
	var enemies [][]string
	for _, char := range chars {
		if strings.Contains(char[0], "_enemy") {
			enemies = append(enemies, char)
		} else {
			players = append(players, char)
		}
	}

	s, err := fileToString("charSelectScreen.html")
	if err != nil {
		panic(err)
	}

	style, err := fileToString("styleHead.html")
	if err != nil {
		panic(err)
	}

	s = style + s

	fmt.Fprint(w, s)

	nameStr, err := fileToString("names.txt")
	if err != nil {
		panic(err)
	}
	names := strings.SplitN(nameStr, ",", -1)

	for i, player := range players {
		name := player[0]
		class := player[1]

		var opponent string
		if len(enemies) >= len(players) {
			opponent = enemies[i][0] + "." + enemies[i][1]
		} else {
			opClass := []string{"Knight", "Archer", "Wizard"}[rand.Intn(3)]
			opName := fmt.Sprintf("%s_enemy", names[rand.Intn(len(names))])
			err := generateChar(opClass, opName)
			if err != nil {
				panic(err)
			}
			opponent = opName + "." + opClass
		}

		fmt.Fprintf(
			w,
			`
			<tr>
				<td>%s</td>
				<td>%s</td>
				<td>
					<form action="/game/%s/%s">
						<input type="submit" name="character" value="Select">
					</form>
				</td>
			</tr>
			`,
			name, class,
			player[0]+"."+player[1], opponent,
		)
	}
	fmt.Fprint(w, `</table></body>`)
}

// gameScreen shows character stats and moves, main game screen
func gameScreen(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cn1 := vars["char1"]
	cn2 := vars["char2"]

	c1, err := readCharFromFile(cn1)
	if err != nil {
		panic(err)
	}
	c2, err := readCharFromFile(cn2)
	if err != nil {
		panic(err)
	}

	c1HTML, err := charToHTML(c1)
	if err != nil {
		panic(err)
	}
	c2HTML, err := charToHTML(c2)
	if err != nil {
		panic(err)
	}

	screen, err := fileToString("fightScreen.html")
	if err != nil {
		panic(err)
	}

	style, err := fileToString("styleHead.html")
	if err != nil {
		panic(err)
	}

	screen = style + screen

	// get log file name
	logName := fmt.Sprintf("%s%s%s%s.log",
		c1.PlayerName, c1.ClassName, c2.PlayerName, c2.ClassName)
	gameLog, err := fileToString(logName)
	if err != nil {
		gameLog = ""
	}

	screen += "<br>" + gameLog

	c1Moves := getMoves(c1)

	info := "Heavy attacks effective against low int (str damage)\n<br>" +
		"Quick attacks effective against low str (dex damage)\n<br>" +
		"Standard attacks effective against low dex (int damage)\n<br>" +
		"Blocks effective with high str, heals armor\n<br>" +
		"Parry effective with high dex, counters but can backfire\n<br>" +
		"Evade effective with high int, heals HP\n<br>"
	info = divWrap(info)

	fmt.Fprintf(w, screen, c1HTML, c2HTML, c1Moves, info)
}

func home(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/selectChar", http.StatusFound)
}

// newRouter returns a router and endpoints
func newRouter() *mux.Router {
	r := mux.NewRouter()

	fileServer := http.FileServer(http.Dir(FILE_DIR))
	handler := http.StripPrefix("/assets/", fileServer)
	r.PathPrefix("/assets/").Handler(handler)

	r.HandleFunc("/", home)
	r.HandleFunc("/newChar", newCharacterScreen).Methods("GET")
	r.HandleFunc("/newChar", parseNewCharForm).Methods("POST")
	r.HandleFunc("/selectChar", characterSelectScreen).Methods("GET")
	r.HandleFunc("/turn/{char1}/{char2}/{move}", parseMoveForm)
	r.HandleFunc("/game/{char1}/{char2}", gameScreen)

	return r
}

// Server starts server using newRouter
func Server() {
	r := newRouter()

	log.Fatal(http.ListenAndServe(PORT, r))
}
