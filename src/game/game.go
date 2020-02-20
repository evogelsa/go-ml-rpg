package game

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
)

type Move int

const (
	HEAVY Move = iota
	QUICK
	STANDARD
	BLOCK
	PARRY
	EVADE
)

// Class is the base player struct storing stats and data
type Class struct {
	PlayerName string
	ClassName  string
	Health     int     // capped 100
	Stamina    int     // capped 100
	Armor      int     // capped 20
	Strength   float32 // normalized
	Dexterity  float32 // normalized
	Intellect  float32 // normalized
}

// NewKnight generates a new knight class with initial values
func NewKnight(playerName string) Class {
	return Class{
		PlayerName: playerName,
		ClassName:  "Knight",
		Health:     rand.Intn(20+1) + 80,
		Stamina:    rand.Intn(20+1) + 50,
		Armor:      rand.Intn(20 + 1),
		Strength:   float32(rand.Intn(5+1)+15) / 20,
		Dexterity:  float32(rand.Intn(5+1)+10) / 20,
		Intellect:  float32(rand.Intn(5+1)+5) / 20,
	}
}

// NewArcher generates a new arhcer lcass with initial values
func NewArcher(playerName string) Class {
	return Class{
		PlayerName: playerName,
		ClassName:  "Archer",
		Health:     rand.Intn(20+1) + 80,
		Stamina:    rand.Intn(20+1) + 50,
		Armor:      rand.Intn(20 + 1),
		Strength:   float32(rand.Intn(5+1)+5) / 20,
		Dexterity:  float32(rand.Intn(5+1)+15) / 20,
		Intellect:  float32(rand.Intn(5+1)+10) / 20,
	}
}

// NewWizard generates a new wizard class with initial values
func NewWizard(playerName string) Class {
	return Class{
		PlayerName: playerName,
		ClassName:  "Wizard",
		Health:     rand.Intn(20+1) + 80,
		Stamina:    rand.Intn(20+1) + 50,
		Armor:      rand.Intn(20 + 1),
		Strength:   float32(rand.Intn(5+1)+10) / 20,
		Dexterity:  float32(rand.Intn(5+1)+5) / 20,
		Intellect:  float32(rand.Intn(5+1)+15) / 20,
	}
}

// heavyAttack deals damaged based off of attackers strength but
// success probability is determined based off defenders intellect
func (c *Class) heavyAttack(e *Class) int {
	// higher enemy intellect -> lower chance to hit
	if rand.Float32() > e.Intellect {
		// higher attacker strength -> more damage
		return int(float32(rand.Intn(20+1))*c.Strength + 1.5)
	} else {
		return 0
	}
}

// quickAttack deals damaged based off of attackers dexterity but
// success probability is determined based off defenders strength
func (c *Class) quickAttack(e *Class) int {
	// higher enemy strength -> lower chance to hit
	if rand.Float32() > e.Strength {
		// higher attacker dexterity -> more damage
		return int(float32(rand.Intn(20+1))*c.Dexterity + 1.5)
	} else {
		return 0
	}
}

// standardAttack deals damaged based off of attacker intellect but
// success probability is determined based off defenders dexterity
func (c *Class) standardAttack(e *Class) int {
	// higer enemy dexterity -> lower chance to hit
	if rand.Float32() > e.Dexterity {
		// higher attacker intellect -> more damage
		return int(float32(rand.Intn(20+1))*c.Intellect + 1.5)
	} else {
		return 0
	}
}

// block attempts to block an attack using strength to determine
// success. on success blocker takes no damage and repairs armor
func (c *Class) blockDefense(e *Class) int {
	if rand.Float32() < c.Strength {
		return rand.Intn(int(c.Strength*10 + .5))
	} else {
		return 0
	}
}

// parry attempts to dodge an attack and counter, success based on
// dexterity of person parrying, failure results in damaging self
func (c *Class) parryDefense(e *Class) int {
	if rand.Float32() < c.Dexterity {
		return rand.Intn(int(c.Dexterity*10 + .5))
	} else {
		return -rand.Intn(int(c.Dexterity*10 + .5))
	}
}

// evade attempts to evade an attack using intelligence to
// determine success. on success evader takes no damage and heals
func (c *Class) evadeDefense(e *Class) int {
	if rand.Float32() < c.Intellect {
		return rand.Intn(int(c.Intellect*10 + .5))
	} else {
		return 0
	}
}

// PrintCharacter prints out a character information box with name
// and character stats
func PrintCharacter(p Class) {
	fmt.Printf("********************\n* %-16s *\n********************\n",
		p.PlayerName)
	fmt.Printf("* Health: %-8d *\n", p.Health)
	fmt.Printf("* Armor: %-9d *\n", p.Armor)
	fmt.Printf("* Strength: %-6d *\n", int(p.Strength*100))
	fmt.Printf("* Dexterity: %-5d *\n", int(p.Dexterity*100))
	fmt.Printf("* Intellect: %-5d *\n", int(p.Intellect*100))
	fmt.Printf("********************\n")
}

// parseMove takes in two players and parses a move m for the first
// player p1
func parseMove(p1, p2 *Class, m Move) (bool, int) {
	var def bool
	var r int
	switch m {
	case HEAVY:
		def = false
		r = p1.heavyAttack(p2)
	case QUICK:
		def = false
		r = p1.quickAttack(p2)
	case STANDARD:
		def = false
		r = p1.standardAttack(p2)
	case BLOCK:
		def = true
		r = p1.blockDefense(p2)
	case PARRY:
		def = true
		r = p1.parryDefense(p2)
	case EVADE:
		def = true
		r = p1.evadeDefense(p2)
	}

	return def, r
}

// handleDamage parses damage to determine armor and health effects
func handleDamage(p *Class, d int) {
	// armor absorbs all damage as long as armor health >0
	if p.Armor > 0 {
		p.Armor -= d
		if p.Armor < 0 {
			p.Armor = 0
		}
	} else {
		p.Health -= d
	}
}

// printTurn prints out text regarding outcome of turn
func printTurn(p1, p2 *Class, m1, m2 Move, s1, s2 bool) string {
	var res string
	var moveNames = map[Move]string{
		HEAVY:    "heavy attack",
		QUICK:    "quick attack",
		STANDARD: "standard attack",
		BLOCK:    "block",
		PARRY:    "parry",
		EVADE:    "evade",
	}
	var result = map[bool]string{
		false: "fails",
		true:  "succeeds",
	}

	res += fmt.Sprintf("%s attempts to %s %s and %s\n",
		p1.PlayerName, moveNames[m1], p2.PlayerName, result[s1])
	res += fmt.Sprintf("%s attempts to %s %s and %s\n",
		p2.PlayerName, moveNames[m2], p1.PlayerName, result[s2])

	return res
}

// printStatus takes in a player and prints out the player name
// and their current health and armor
func printStatus(p *Class) string {
	var res string
	res += fmt.Sprintf("%s\n", p.PlayerName)
	res += fmt.Sprintf("\tHealth: %d\n", p.Health)
	res += fmt.Sprintf("\tArmor: %d\n", p.Armor)
	return res
}

// getIntInput gets an error checked integer input from os.stdin
func getIntInput(w io.Writer) int {
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
			fmt.Fprintf(w, "Please input a valid integer selection\n")
		}
	}
}

// SelectMove prints move selection to stdout and calls getIntInput
func SelectMove(w io.Writer, p Class) Move {
	fmt.Fprintf(w, "%s select move:\n", p.PlayerName)
	fmt.Fprintf(w, "\t(0)Heavy attack\n\t(1)Quick attack\n\t(2)Standard attack\n")
	fmt.Fprintf(w, "\t(3)Block\n\t(4)Parry\n\t(5)Evade\n")
	return Move(getIntInput(w))
}

// Turn takes in two players and two moves and and handles the events
// which occur from player p1 executing move m1 and p2 executing m2
func Turn(p1, p2 *Class, m1 Move) string {
	var res string
	m2 := Move(rand.Intn(6))
	def1, a1 := parseMove(p1, p2, m1)
	def2, a2 := parseMove(p2, p1, m2)
	if !def1 && !def2 {
		// both players attack, handle damage
		var s1, s2 bool
		if a1 > 0 {
			handleDamage(p2, a1)
			s1 = true
		}
		if a2 > 0 {
			handleDamage(p1, a2)
			s2 = true
		}
		printTurn(p1, p2, m1, m2, s1, s2)
		if s1 {
			res += fmt.Sprintf("%s deals %d damage\n", p1.PlayerName, a1)
		}
		if s2 {
			res += fmt.Sprintf("%s deals %d damage\n", p2.PlayerName, a2)
		}
	} else if !def1 && def2 {
		// if player2 defends then use attack of player1 to
		// determine the outcome
		switch m2 {
		case BLOCK:
			// if p2 success, repair armor of p2, else p2
			// takes damage of p1
			if a2 > 0 {
				p2.Armor += a2
				res += printTurn(p1, p2, m1, m2, false, true)
				res += fmt.Sprintf("%s repairs armor for %d points\n",
					p2.PlayerName, a2)
			} else if a1 > 0 {
				handleDamage(p2, a1)
				res += printTurn(p1, p2, m1, m2, true, false)
				res += fmt.Sprintf("%s deals %d damage\n", p1.PlayerName, a1)
			} else {
				res += printTurn(p1, p2, m1, m2, false, false)
			}
		case PARRY:
			// if p2 success, p1 takes damage of p1 attack +
			// damage of p2 counter, vise versa if fail
			if a2 > 0 {
				handleDamage(p1, a1+a2)
				res += printTurn(p1, p2, m1, m2, false, true)
				res += fmt.Sprintf("%s deals %d damage\n", p2.PlayerName, a1+a2)
			} else if a1-a2 > 0 {
				handleDamage(p2, a1-a2)
				res += printTurn(p1, p2, m1, m2, true, false)
				res += fmt.Sprintf("%s deals %d damage\n", p1.PlayerName, a1-a2)
			} else {
				res += printTurn(p1, p2, m1, m2, false, false)
			}
		case EVADE:
			// if p2 success p2 heals a little, on fail p2
			// takes damage of p1 attack
			if a2 > 0 {
				p2.Health += a2
				res += printTurn(p1, p2, m1, m2, false, true)
				res += fmt.Sprintf("%s heals %d damage\n", p2.PlayerName, a2)
			} else if a1 > 0 {
				handleDamage(p2, a1)
				res += printTurn(p1, p2, m1, m2, true, false)
				res += fmt.Sprintf("%s deals %d damage\n", p1.PlayerName, a1)
			} else {
				res += printTurn(p1, p2, m1, m2, false, false)
			}
		}
	} else if def1 && !def2 {
		// if player1 defends then use attack of player2 to
		// determine the outcome
		switch m1 {
		case BLOCK:
			// if p1 success, repair armor of p1, else p1
			// takes damage of p2
			if a1 > 0 {
				p1.Armor += a1
				res += printTurn(p1, p2, m1, m2, true, false)
				res += fmt.Sprintf("%s repairs armor for %d points\n",
					p1.PlayerName, a1)
			} else if a2 > 0 {
				handleDamage(p1, a2)
				res += printTurn(p1, p2, m1, m2, false, true)
				res += fmt.Sprintf("%s deals %d damage\n", p2.PlayerName, a2)
			} else {
				res += printTurn(p1, p2, m1, m2, false, false)
			}
		case PARRY:
			// if p1 success, p2 takes damage of p2 attack +
			// damage of p1 counter, vise versa if fail
			if a1 > 0 {
				handleDamage(p2, a2+a1)
				res += printTurn(p1, p2, m1, m2, true, false)
				res += fmt.Sprintf("%s deals %d damage\n", p1.PlayerName, a2+a1)
			} else if a2-a1 > 0 {
				handleDamage(p1, a2-a1)
				res += printTurn(p1, p2, m1, m2, false, true)
				res += fmt.Sprintf("%s deals %d damage\n", p2.PlayerName, a2-a1)
			} else {
				res += printTurn(p1, p2, m1, m2, false, false)
			}
		case EVADE:
			// if p1 success p1 heals a little, on fail p1
			// takes damage of p2 attack
			if a1 > 0 {
				p1.Health += a1
				res += printTurn(p1, p2, m1, m2, true, false)
				res += fmt.Sprintf("%s heals %d damage\n", p1.PlayerName, a1)
			} else if a2 > 0 {
				handleDamage(p1, a2)
				res += printTurn(p1, p2, m1, m2, false, true)
				res += fmt.Sprintf("%s deals %d damage\n", p2.PlayerName, a2)
			} else {
				res += printTurn(p1, p2, m1, m2, false, false)
			}
		}
	} else if def1 && def2 {
		// nothing happens both players tried to defend
		res += printTurn(p1, p2, m1, m2, false, false)
		res += fmt.Sprintf("Nothing happens!\n")
	}
	// res += printStatus(p1)
	// res += printStatus(p2)
	if res == "" {
		res += "Both attacks miss!\n"
	}
	res += "----------------------------------------\n"
	return res
}
