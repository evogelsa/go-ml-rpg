package game

import "testing"

// how many test classes to generate
const nClassTests int = 1000

// printClassErrs prints out error cases from class tests
func printClassErrs(
	pns, cns, hs, ss, as, strs, dexs, ints [][]interface{},
	t *testing.T,
) {

	for _, i := range pns {
		t.Errorf("Got PlayerName %s on test %d, expected Archer", i...)
	}
	for _, i := range cns {
		t.Errorf("Got ClassName %s on test %d, expected Archer", i...)
	}
	for _, i := range hs {
		t.Errorf("Health out of bounds, got %d on test %d", i...)
	}
	for _, i := range ss {
		t.Errorf("Stamina out of bounds, got %d on test %d", i...)
	}
	for _, i := range as {
		t.Errorf("Armor out of bounds, got %d on test %d", i...)
	}
	for _, i := range strs {
		t.Errorf("Strength out of bounds, got %f on test %d", i...)
	}
	for _, i := range dexs {
		t.Errorf("Dexterity out of bounds, got %f on test %d", i...)
	}
	for _, i := range ints {
		t.Errorf("Intellect out of bounds, got %f on test %d", i...)
	}
}

// TestNewKnight generates nClassTests knights and checks thats stats
// are within expected values
func TestNewKnight(t *testing.T) {
	var pns, cns [][]interface{}
	var hs, ss, as [][]interface{}
	var strs, dexs, ints [][]interface{}
	for i := 0; i < nClassTests; i++ {
		c := NewKnight("Knight")

		if c.PlayerName != "Knight" {
			pns = append(pns, []interface{}{c.PlayerName, i})
		}
		if c.ClassName != "Knight" {
			cns = append(cns, []interface{}{c.ClassName, i})
		}
		if c.Health < 80 || c.Health > 100 {
			hs = append(hs, []interface{}{c.Health, i})
		}
		if c.Stamina < 50 || c.Stamina > 70 {
			ss = append(ss, []interface{}{c.Stamina, i})
		}
		if c.Armor < 0 || c.Armor > 20 {
			as = append(as, []interface{}{c.Armor, i})
		}
		if c.Strength < .75 || c.Strength > 1 {
			strs = append(strs, []interface{}{c.Strength, i})
		}
		if c.Dexterity < .50 || c.Dexterity > .75 {
			dexs = append(dexs, []interface{}{c.Dexterity, i})
		}
		if c.Intellect < .25 || c.Intellect > .5 {
			ints = append(ints, []interface{}{c.Intellect, i})
		}
	}
	printClassErrs(pns, cns, hs, ss, as, strs, dexs, ints, t)
}

// TestNewArcher generates nClassTests archers and checks thats stats
// are within expected values
func TestNewArcher(t *testing.T) {
	var pns, cns [][]interface{}
	var hs, ss, as [][]interface{}
	var strs, dexs, ints [][]interface{}
	for i := 0; i < nClassTests; i++ {
		c := NewArcher("Archer")

		if c.PlayerName != "Archer" {
			pns = append(pns, []interface{}{c.PlayerName, i})
		}
		if c.ClassName != "Archer" {
			cns = append(cns, []interface{}{c.ClassName, i})
		}
		if c.Health < 80 || c.Health > 100 {
			hs = append(hs, []interface{}{c.Health, i})
		}
		if c.Stamina < 50 || c.Stamina > 70 {
			ss = append(ss, []interface{}{c.Stamina, i})
		}
		if c.Armor < 0 || c.Armor > 20 {
			as = append(as, []interface{}{c.Armor, i})
		}
		if c.Strength < .25 || c.Strength > .50 {
			strs = append(strs, []interface{}{c.Strength, i})
		}
		if c.Dexterity < .75 || c.Dexterity > 1 {
			dexs = append(dexs, []interface{}{c.Dexterity, i})
		}
		if c.Intellect < .50 || c.Intellect > .75 {
			ints = append(ints, []interface{}{c.Intellect, i})
		}
	}
	printClassErrs(pns, cns, hs, ss, as, strs, dexs, ints, t)
}

// TestNewWizard generates nClassTests wizards and checks thats stats
// are within expected values
func TestNewWizard(t *testing.T) {
	var pns, cns [][]interface{}
	var hs, ss, as [][]interface{}
	var strs, dexs, ints [][]interface{}
	for i := 0; i < nClassTests; i++ {
		c := NewWizard("Wizard")

		if c.PlayerName != "Wizard" {
			pns = append(pns, []interface{}{c.PlayerName, i})
		}
		if c.ClassName != "Wizard" {
			cns = append(cns, []interface{}{c.ClassName, i})
		}
		if c.Health < 80 || c.Health > 100 {
			hs = append(hs, []interface{}{c.Health, i})
		}
		if c.Stamina < 50 || c.Stamina > 70 {
			ss = append(ss, []interface{}{c.Stamina, i})
		}
		if c.Armor < 0 || c.Armor > 20 {
			as = append(as, []interface{}{c.Armor, i})
		}
		if c.Strength < .50 || c.Strength > .75 {
			strs = append(strs, []interface{}{c.Strength, i})
		}
		if c.Dexterity < .25 || c.Dexterity > .50 {
			dexs = append(dexs, []interface{}{c.Dexterity, i})
		}
		if c.Intellect < .75 || c.Intellect > 1 {
			ints = append(ints, []interface{}{c.Intellect, i})
		}
	}
	printClassErrs(pns, cns, hs, ss, as, strs, dexs, ints, t)
}
