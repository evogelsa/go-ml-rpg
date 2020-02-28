package game

import (
	"fmt"
	"math/rand"
	"reflect"
)

//Q(s,a) = reward + discount*(max future reward)

// normalize takes a number with values between min and
// max and normalizes it to between 0 and 1
func normalize(n, min, max float32) float32 {
	return (n - min) / (max - min)
}

func getStates(p1, p2 *Class) [][]float32 {
	var p1I interface{} = p1
	var p2I interface{} = p2

	// get values of p1 class
	vOf1 := reflect.ValueOf(p1I).Elem()
	var sOf1 []float32
	for i := 2; i < vOf1.NumField(); i++ {
		// ignore player and class name, convert all to float32
		v := vOf1.Field(i).Interface().(float32)
		// normalize health, stamina, armor
		if i < 4 {
			v = normalize(v, 0, 100)
		} else if i == 4 {
			v = normalize(v, 0, 20)
		}

		sOf1 = append(sOf1, v)
	}
	fmt.Println(sOf1)

	// get values of p2 class
	vOf2 := reflect.ValueOf(p2I).Elem()
	var sOf2 []float32
	for i := 2; i < vOf2.NumField(); i++ {
		// ignore player and class name, convert all to float32
		v := vOf2.Field(i).Interface().(float32)
		// normalize health, stamina, and armor
		if i < 4 {
			v = normalize(v, 0, 100)
		} else if i == 4 {
			v = normalize(v, 0, 20)
		}

		sOf2 = append(sOf2, v)
	}
	fmt.Println(sOf2)

	var tab = make([][]float32, len(sOf2))
	for i := 0; i < len(tab); i++ {
		tab[i] = make([]float32, len(sOf1))
	}

	for i := 0; i < len(sOf2); i++ {
		for j := 0; j < len(sOf1); j++ {
			tab[i][j] = sOf2[i] * sOf1[j]
		}
	}

	return tab
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
	// enemy block prob * avg heal - enemy fail block prob * avg damage
	avgB := (((e.Strength) * float32(int(10*e.Strength+.5))) -
		((1 - e.Strength) * avgPlayer))
	// enemy fail parry prob * extra dmg + avg damage
	avgP := -(((1 - e.Dexterity) * float32(int(10*e.Dexterity+.5))) + avgPlayer)
	// enemy fail evade prob * avg damage
	avgE := -((1 - e.Intellect) * avgPlayer)

	return []float32{avgH, avgQ, avgS, avgB, avgP, avgE}
}

func minMaxArmor(p, e *Class) []float32 {

	return []float32{0}
}

func aiGetTurn(p1, p2 *Class) Move {
	return Move(rand.Intn(6))
}