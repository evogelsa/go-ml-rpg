package game

import (
	"encoding/gob"
	"fmt"
	"math"
	"math/rand"
	"os"
	"sync"
)

type Algorithm int

const (
	AI_MINMAX Algorithm = iota
	AI_RAND
	AI_REINFORCEMENT
)

const (
	phMask = 0x0300
	paMask = 0x00C0
	ehMask = 0x000C
	eaMask = 0x0003
)

var Train bool
var LearningRate float32
var Discount float32

var ExploreRate float32
var turns int
var exploreMutex sync.Mutex

var madeQT bool = false
var madeLock sync.Mutex

var QT map[uint16][]float32
var qtMutex sync.RWMutex

// AI_ALG is set by main which uses cmd line flags to
// set the desired ai algorithm
var AI_ALG Algorithm

// normalize takes a number with values between min and
// max and normalizes it to between 0 and 1
func normalize(n, min, max float32) float32 {
	return (n - min) / (max - min)
}

// initQT enumerates all states and initializes their qtable rows
// to be zeros
func initQT() {
	qtMutex.Lock()
	defer qtMutex.Unlock()

	// check if qtable already exists
	_, err := os.Stat("qtable")
	if err == nil {
		// file exists
		loadQT()
	} else {
		QT = make(map[uint16][]float32)
		for h1 := 0; h1 <= 2; h1++ {
			for a1 := 0; a1 <= 2; a1++ {
				for c1 := 0; c1 <= 2; c1++ {
					for h2 := 0; h2 <= 2; h2++ {
						for a2 := 0; a2 <= 2; a2++ {
							for c2 := 0; c2 <= 2; c2++ {
								state := uint16((h1 << 10) + (a1 << 8) + (c1 << 6) +
									(h2 << 4) + (a2 << 2) + c2)
								QT[state] = []float32{0, 0, 0, 0, 0, 0}
							}
						}
					}
				}
			}
		}
	}
}

// loadQT loads qTable from file. Only called from initQT since
// there is no mutex enforcement within
func loadQT() {
	// only call from initQT
	f, err := os.Open("qtable")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	dec := gob.NewDecoder(f)
	err = dec.Decode(&QT)
	if err != nil {
		panic(err)
	}
}

// saveQT saves qtable to disk
func saveQT() {
	qtMutex.RLock()
	defer qtMutex.RUnlock()

	f, err := os.Create("qtable")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	enc := gob.NewEncoder(f)
	err = enc.Encode(QT)
	if err != nil {
		panic(err)
	}
}

func getState(p, e *Class) uint16 {
	var state uint16

	switch p.ClassName {
	case "Knight":
		state += 0 << 10
	case "Archer":
		state += 1 << 10
	case "Wizard":
		state += 2 << 10
	}

	if p.Health < 25 {
		state += 0 << 8
	} else if p.Health < 50 {
		state += 1 << 8
	} else {
		state += 2 << 8
	}

	if p.Armor < 5 {
		state += 0 << 6
	} else if p.Armor < 10 {
		state += 1 << 6
	} else {
		state += 2 << 6
	}

	switch e.ClassName {
	case "Knight":
		state += 0 << 4
	case "Archer":
		state += 1 << 4
	case "Wizard":
		state += 2 << 4
	}

	if e.Health < 25 {
		state += 0 << 2
	} else if e.Health < 50 {
		state += 1 << 2
	} else {
		state += 2 << 2
	}

	if e.Armor < 5 {
		state += 0
	} else if e.Armor < 10 {
		state += 1
	} else {
		state += 2
	}

	return state
}

func updateQT(state, nextState uint16, action Move) {
	if Train {
		// lock QT for reading
		qtMutex.RLock()
		// get current q value
		qv := QT[state][action]
		// get max q future
		qf := QT[nextState]
		// unlock QT
		qtMutex.RUnlock()

		var max float32 = -math.MaxFloat32
		for _, v := range qf {
			if max < v {
				max = v
			}
		}
		// get reward (check if state values change)
		var reward float32
		// extract p health, if decrease + reward
		ph := (state & phMask) >> 8
		phNext := (nextState & phMask) >> 8
		if phNext < ph {
			reward += 1.5
		}
		// extract p armor, if decrease + reward
		pa := (state & paMask) >> 6
		paNext := (nextState & paMask) >> 6
		if paNext < pa {
			reward += 1.5
		}
		// extract e health, if increase + reward, if dec small penalty
		eh := (state & ehMask) >> 2
		ehNext := (nextState & ehMask) >> 2
		if ehNext > eh {
			reward += 1
		} else if ehNext < eh {
			reward -= .5
		}
		// extract e armor, if increase + reward, if dec small penalty
		ea := (state & eaMask)
		eaNext := (nextState & eaMask)
		if eaNext > ea {
			reward += 1
		} else if eaNext < ea {
			reward -= .5
		}

		// lock qt for writing
		qtMutex.Lock()
		QT[state][action] = qv + LearningRate*(reward+Discount*max-qv)
		fmt.Printf("%f = ", QT[state][action])
		fmt.Printf("%f + %f*(%f+%f*%f-%f)\n",
			qv, LearningRate, reward, Discount, max, qv)
		qtMutex.Unlock()
	}
}

// minMaxDamage returns array of avg outcome for each move wrt
// damage dealt to player
func minMaxDamage(p, e *Class) []float32 {
	// actions: heavy, quick, standard, block, parry, evade
	// calculate average damage to player of each attack
	// chance of success * avg outcome = avg val
	avgH := (1 - p.Intellect) * float32(int(10*e.Strength+1.5))
	avgQ := (1 - p.Strength) * float32(int(10*e.Dexterity+1.5))
	avgS := (1 - p.Dexterity) * float32(int(10*e.Intellect+1.5))
	avgB := float32(0)
	avgP := (e.Dexterity) * float32(int(10*e.Dexterity+.5))
	avgE := float32(0)

	return []float32{avgH, avgQ, avgS, avgB, avgP, avgE}
}

// minMaxHealth returns array of avg outcome for each move wrt
// change in enemy health
func minMaxHealth(p, e *Class) []float32 {
	avgsFromPlayer := minMaxDamage(e, p)
	var avgPlayer float32
	for _, v := range avgsFromPlayer {
		avgPlayer += v
	}
	avgPlayer /= 6

	avgH := -avgPlayer
	avgQ := -avgPlayer
	avgS := -avgPlayer
	// enemy fail block prob * avg damage
	avgB := ((1 - e.Strength) * avgPlayer)
	// enemy fail parry prob * extra dmg + avg damage
	avgP := -(((1 - e.Dexterity) * float32(int(10*e.Dexterity+.5))) + avgPlayer)
	// enemy evade prob * heal - enemy fail evade prob * avg damage
	avgE := ((e.Intellect) * float32(int(10*e.Intellect+.5))) -
		((1 - e.Intellect) * avgPlayer)

	return []float32{avgH, avgQ, avgS, avgB, avgP, avgE}
}

// min/maxArmor gets avg outcome for each move wrt to change in
// enemy armor
func minMaxArmor(p, e *Class) []float32 {
	avgsFromPlayer := minMaxDamage(e, p)
	var avgPlayer float32
	for _, v := range avgsFromPlayer {
		avgPlayer += v
	}
	avgPlayer /= 6

	avgH := -avgPlayer
	avgQ := -avgPlayer
	avgS := -avgPlayer
	// enemy block prob * avg heal - enemy fail block prob * avg damage
	avgB := (((e.Strength) * float32(int(10*e.Strength+.5))) -
		((1 - e.Strength) * avgPlayer))
	// enemy fail parry prob * extra dmg + avg damage
	avgP := -(((1 - e.Dexterity) * float32(int(10*e.Dexterity+.5))) + avgPlayer)
	// enemy fail evade prob * avg damage
	avgE := -((1 - e.Intellect) * avgPlayer)

	return []float32{avgH, avgQ, avgS, avgB, avgP, avgE}
}

// getMinMaxAll calls the three minmax funcs and returns array containing
// all of the results.
func getMinMaxAll(p, e *Class) [][]float32 {
	return [][]float32{
		minMaxDamage(p, e),
		minMaxHealth(p, e),
		minMaxArmor(p, e),
	}
}

// normalizedMinMaxes returns a normalized slice containing weights
// for each move, sums to 1
func normalizedMinMaxes(p, e *Class) []float32 {
	minMaxes := getMinMaxAll(p, e)

	var max float32 = -math.MaxFloat32
	for _, mm := range minMaxes {
		for _, v := range mm {
			if v > max {
				max = v
			}
		}
	}

	var min float32 = math.MaxFloat32
	for _, mm := range minMaxes {
		for _, v := range mm {
			if v < min {
				min = v
			}
		}
	}

	for i, mm := range minMaxes {
		for j := range mm {
			minMaxes[i][j] = normalize(minMaxes[i][j], min, max)
		}
	}

	var sum float32
	for _, mm := range minMaxes {
		for _, v := range mm {
			sum += v
		}
	}

	for i, mm := range minMaxes {
		for j := range mm {
			minMaxes[i][j] = normalize(minMaxes[i][j], 0, sum)
		}
	}

	var vals []float32
	for _, r := range minMaxes {
		vals = append(vals, r...)
	}

	for i := 6; i < len(vals); i++ {
		vals[i%6] += vals[i]
	}

	var ret []float32 = vals[:6]

	return ret
}

// getTurnMinMax uses the minmax strategy to determine weighted probabilities
// for each move and pseudorandomly selects the move to use based on weights
func getTurnMinMax(p, e *Class) Move {
	minMaxes := normalizedMinMaxes(p, e)

	r := rand.Float32()
	var m int
	for i, v := range minMaxes {
		r -= v
		if r <= 0 {
			m = i
			break
		}
	}
	return Move(m)
}

// getTurnRand randomly selects a move to use
func getTurnRand() Move {
	return Move(rand.Intn(6))
}

func getTurnReinforcement(p, e *Class) Move {
	// lock mutex here to prevent two possible inits
	madeLock.Lock()
	if !madeQT {
		initQT()
		madeQT = true
		fmt.Println("QT initialized")
	} else {
		saveQT()
	}
	madeLock.Unlock()

	// get state
	state := getState(p, e)
	// select next action, check explore or exploit
	exploreMutex.Lock()
	var action Move
	if rand.Float32() < ExploreRate {
		turns++
		action = getTurnRand()
		if ExploreRate > .25 {
			ExploreRate -= (float32(turns) * .001)
		}
	} else {
		var max float32 = -math.MaxFloat32
		for i, v := range QT[state] {
			if max < v {
				max = v
				action = Move(i)
			}
		}
	}
	exploreMutex.Unlock()
	return action
}

// aiGetTurn handles getting the next move of the AI using whatever strategy was
// selected at server launch
func AIGetTurn(p, e *Class) Move {
	var m Move

	switch AI_ALG {
	case AI_MINMAX:
		m = getTurnMinMax(p, e)
	case AI_RAND:
		m = getTurnRand()
	case AI_REINFORCEMENT:
		m = getTurnReinforcement(p, e)
	}

	return m
}
