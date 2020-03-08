package web

import (
	"bufio"
	"encoding/gob"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.iu.edu/evogelsa/go-ml-rpg/game"

	"github.com/gorilla/mux"
)

const (
	FILE_DIR = "./web/assets/"
	SAVE_DIR = "./saves/"
	IMG_DIR  = SAVE_DIR + "imgs/"
	CHAR_DIR = SAVE_DIR + "characters/"
	LOG_DIR  = SAVE_DIR + "logs/"
)

var enemyMap map[string]string
var enemyMapLock sync.RWMutex
var enemyMapLoaded bool

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

func loadLog(fn string) (string, error) {
	fn = LOG_DIR + fn

	body, err := ioutil.ReadFile(fn)
	if err != nil {
		return "", err
	}

	text := string(body)
	return text, nil
}

// logFile adds logStr to file specified by fn inside FILE_DIR
func logFile(fn, logStr string) error {
	fn = LOG_DIR + fn

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

	fn = CHAR_DIR + fn
	f, err := os.Open(fn)
	if err != nil {
		return char, err
	}
	defer f.Close()

	dec := gob.NewDecoder(f)
	err = dec.Decode(&char)
	if err != nil {
		return char, err
	}

	return char, nil
}

// writeCharToFile takes in character class and writes character info
// to a file called CharName.CharClass
func writeCharToFile(c game.Class) error {
	//get file from fn (overwrites if file exists)
	fn := CHAR_DIR + c.PlayerName + "." + c.ClassName
	f, err := os.Create(fn)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := gob.NewEncoder(f)
	err = enc.Encode(c)
	if err != nil {
		return err
	}

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
	enemyMove := game.AIGetTurn(&c1, &c2)
	outStr, end := game.Turn(&c1, &c2, move, enemyMove)

	// write turns to file
	err = setImages(c1, c2, move, enemyMove)
	if err != nil {
		panic(err)
	}

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

	var redirect string
	if end {
		redirect = "/end/" + char1Name + "/" + char2Name
	} else {
		redirect = "/game/" + char1Name + "/" + char2Name
	}
	http.Redirect(w, r, redirect, http.StatusFound)
}

// parseNewCharForm handles extracting the name and class from character
// creation screen and generating a new character with that info
func parseNewCharForm(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println(fmt.Errorf("Error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	class := r.Form.Get("class")
	name := r.Form.Get("name")

	if name == "" {
		style, err := fileToString("styleHead.html")
		if err != nil {
			panic(err)
		}
		body, err := fileToString("newCharacterScreen.html")
		if err != nil {
			panic(err)
		}
		fmt.Fprint(w, style+body)
		fmt.Fprint(w, `<h3 style="color:red">Character name cannot be empty!</h3>`)
	} else {
		err = generateChar(class, name)
		if err != nil {
			panic(err)
		}

		http.Redirect(w, r, "/selectChar", http.StatusFound)
	}
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

func initEnemyMap() {
	enemyMapLock.Lock()
	defer enemyMapLock.Unlock()

	// check if qtable already exists
	_, err := os.Stat(SAVE_DIR + "enemy_map")
	if err == nil {
		loadEnemyMap()
	} else {
		enemyMap = make(map[string]string)
	}
	enemyMapLoaded = true
}

func loadEnemyMap() {
	// only call from initQT
	f, err := os.Open(SAVE_DIR + "enemy_map")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	dec := gob.NewDecoder(f)
	err = dec.Decode(&enemyMap)
	if err != nil {
		panic(err)
	}
}

func saveEnemyMap() {
	enemyMapLock.RLock()
	defer enemyMapLock.RUnlock()

	f, err := os.Create(SAVE_DIR + "enemy_map")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	enc := gob.NewEncoder(f)
	err = enc.Encode(enemyMap)
	if err != nil {
		panic(err)
	}
}

func getOpponent(character string) string {
	enemyMapLock.RLock()

	opponent, ok := enemyMap[character]
	enemyMapLock.RUnlock()
	if !ok {
		nameStr, err := fileToString("names.txt")
		if err != nil {
			panic(err)
		}
		names := strings.SplitN(nameStr, ",", -1)

		opClass := []string{"Knight", "Archer", "Wizard"}[rand.Intn(3)]
		opName := fmt.Sprintf("%s_enemy", names[rand.Intn(len(names))])
		err = generateChar(opClass, opName)
		if err != nil {
			panic(err)
		}
		opponent = opName + "." + opClass

		enemyMapLock.Lock()
		enemyMap[character] = opponent
		enemyMapLock.Unlock()
		saveEnemyMap()
	}
	return opponent
}

// characterSelectScreen displays all characters in save dir with
// options to select each char
func characterSelectScreen(w http.ResponseWriter, r *http.Request) {
	if !enemyMapLoaded {
		initEnemyMap()
	}

	charFiles, err := ioutil.ReadDir(CHAR_DIR)
	if err != nil {
		panic(err)
	}

	var chars [][]string
	for _, charFile := range charFiles {
		chars = append(chars, strings.SplitN(charFile.Name(), ".", -1))
	}

	var players [][]string
	for _, char := range chars {
		if !strings.Contains(char[0], "_enemy") {
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

	for _, player := range players {
		name := player[0]
		class := player[1]

		character := name + "." + class
		opponent := getOpponent(character)

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
				<td>
					<form action="/deleteChar/%s/%s">
						<input type="submit" name="character" value="Delete">
					</form>
				</td>
			</tr>
			`,
			name, class,
			character, opponent,
			character, opponent,
		)
	}
	fmt.Fprint(w, `</table></body>`)
}

func setImages(c1, c2 game.Class, m1, m2 game.Move) error {
	filePrefix := c1.PlayerName + c1.ClassName + c2.PlayerName + c2.ClassName

	fn := IMG_DIR + filePrefix + ".images"
	f, err := os.Create(fn)
	if err != nil {
		return err
	}
	defer f.Close()

	// parse moves to string
	var m1Str string
	switch m1 {
	case game.HEAVY:
		m1Str = "HEAVY"
	case game.QUICK:
		m1Str = "QUICK"
	case game.STANDARD:
		m1Str = "STANDARD"
	case game.BLOCK:
		m1Str = "BLOCK"
	case game.PARRY:
		m1Str = "PARRY"
	case game.EVADE:
		m1Str = "EVADE"
	}

	var m2Str string
	switch m2 {
	case game.HEAVY:
		m2Str = "HEAVY"
	case game.QUICK:
		m2Str = "QUICK"
	case game.STANDARD:
		m2Str = "STANDARD"
	case game.BLOCK:
		m2Str = "BLOCK"
	case game.PARRY:
		m2Str = "PARRY"
	case game.EVADE:
		m2Str = "EVADE"
	}

	fmt.Fprint(
		f,
		c1.ClassName+"-"+m1Str+".png"+"\n",
		c2.ClassName+"-"+m2Str+".png"+"\n",
	)

	return nil
}

func getImageStrings(c1, c2 game.Class) ([]string, error) {
	var images []string

	filePrefix := c1.PlayerName + c1.ClassName + c2.PlayerName + c2.ClassName

	fn := IMG_DIR + filePrefix + ".images"
	f, err := os.Open(fn)
	if err != nil {
		if os.IsNotExist(err) {
			f, err = os.Create(fn)
			if err != nil {
				return images, err
			}
			fmt.Fprint(
				f,
				c1.ClassName+"-"+"IDLE"+".png"+"\n",
				c2.ClassName+"-"+"IDLE"+".png"+"\n",
			)
			f.Close()
			f, err = os.Open(fn)
			if err != nil {
				return images, err
			}
		} else {
			return images, err
		}
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	scanner.Scan()
	img1 := "../../assets/imgs/" + scanner.Text()

	scanner.Scan()
	img2 := "../../assets/imgs/" + scanner.Text()

	images = []string{img1, img2}

	//open file, parse two lines
	//return two imagefile name corresponding to both moves
	return images, nil
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

	// check for game over for when select dead char
	if c1.Health <= 0 || c2.Health <= 0 {
		redirect := "/end/" + cn1 + "/" + cn2
		http.Redirect(w, r, redirect, http.StatusMovedPermanently)
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
	gameLog, err := loadLog(logName)
	if err != nil {
		gameLog = ""
	}

	rthButton, err := fileToString("returnToHomeButton.html")
	if err != nil {
		panic(err)
	}

	// screen += "<br>" + rthButton + "<br>" + gameLog
	screen += "<br>" + gameLog

	c1Moves := getMoves(c1)

	info := "Heavy attacks effective against low int (str damage)\n<br>" +
		"Quick attacks effective against low str (dex damage)\n<br>" +
		"Standard attacks effective against low dex (int damage)\n<br>" +
		"Blocks effective with high str, heals armor\n<br>" +
		"Parry effective with high dex, counters but can backfire\n<br>" +
		"Evade effective with high int, heals HP\n<br>"
	info = divWrap(info)

	info += "<br><br>" + rthButton

	images, err := getImageStrings(c1, c2)
	if err != nil {
		panic(err)
	}

	fmt.Fprintf(w, screen, c1HTML, c2HTML, c1Moves, info, images[0], images[1])
}

func gameEnd(w http.ResponseWriter, r *http.Request) {
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

	style, err := fileToString("styleHead.html")
	if err != nil {
		panic(err)
	}

	// get log file name
	logName := fmt.Sprintf("%s%s%s%s.log",
		c1.PlayerName, c1.ClassName, c2.PlayerName, c2.ClassName)
	gameLog, err := loadLog(logName)
	if err != nil {
		gameLog = ""
	}

	//get button to rth
	button, err := fileToString("returnToHomeButton.html")
	if err != nil {
		panic(err)
	}

	screen := style + button + "<br>" + gameLog

	fmt.Fprint(w, screen)
}

func deleteChar(w http.ResponseWriter, r *http.Request) {
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

	err = os.Remove(CHAR_DIR + cn1)
	if err != nil {
		panic(err)
	}

	err = os.Remove(CHAR_DIR + cn2)
	if err != nil {
		panic(err)
	}

	logName := fmt.Sprintf("%s%s%s%s.log",
		c1.PlayerName, c1.ClassName, c2.PlayerName, c2.ClassName)

	os.Remove(LOG_DIR + logName)

	imagesName := fmt.Sprintf("%s%s%s%s.images",
		c1.PlayerName, c1.ClassName, c2.PlayerName, c2.ClassName)

	os.Remove(IMG_DIR + imagesName)

	http.Redirect(w, r, "/selectChar", http.StatusFound)
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
	r.HandleFunc("/deleteChar/{char1}/{char2}", deleteChar)
	r.HandleFunc("/turn/{char1}/{char2}/{move}", parseMoveForm)
	r.HandleFunc("/game/{char1}/{char2}", gameScreen)
	r.HandleFunc("/end/{char1}/{char2}", gameEnd)

	return r
}

// Server starts server using newRouter
func Server(port string) {
	r := newRouter()

	log.Fatal(http.ListenAndServe(port, r))
}
