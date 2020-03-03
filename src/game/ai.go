package game

import (
	"math"
	"math/rand"
	"reflect"
)

type Algorithm int

const (
	AI_MINMAX Algorithm = iota
	AI_RAND
	AI_NEURAL
)

// AI_ALG is set by main which uses cmd line flags to
// set the desired ai algorithm
var AI_ALG Algorithm

// normalize takes a number with values between min and
// max and normalizes it to between 0 and 1
func normalize(n, min, max float32) float32 {
	return (n - min) / (max - min)
}

func enumStates() [][]float32 {
	return [][]float32{{0}}
}

// getState returns an interface containing health, armor, and stats
// for both players in the game.
func getState(p1, p2 *Class) []interface{} {
	var p1I interface{} = p1
	var p2I interface{} = p2

	// get values of p1 class
	vOf1 := reflect.ValueOf(p1I).Elem()
	var sOf1 []interface{}
	for i := 2; i < vOf1.NumField(); i++ {
		// ignore player and class name, convert all to float32
		v := vOf1.Field(i).Interface()
		sOf1 = append(sOf1, v)
	}

	// get values of p2 class
	vOf2 := reflect.ValueOf(p2I).Elem()
	var sOf2 []interface{}
	for i := 2; i < vOf2.NumField(); i++ {
		// ignore player and class name, convert all to float32
		v := vOf2.Field(i).Interface()
		sOf2 = append(sOf2, v)
	}

	states := make([]interface{}, len(sOf1)+len(sOf2))
	for i := 0; i < len(sOf1); i++ {
		states[i] = (sOf1[i])
	}
	for i := len(sOf1); i < len(sOf2)+len(sOf1); i++ {
		states[i] = sOf2[i-len(sOf1)]
	}

	return states
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
func getTurnRand(p, e *Class) Move {
	return Move(rand.Intn(6))
}

// aiGetTurn handles getting the next move of the AI using whatever strategy was
// selected at server launch
func AIGetTurn(p, e *Class) Move {
	var m Move

	switch AI_ALG {
	case AI_MINMAX:
		m = getTurnMinMax(p, e)
	case AI_RAND:
		m = getTurnRand(p, e)
	}

	return m
}
