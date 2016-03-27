package game

import "testing"

func TestNewMap(t *testing.T) {
	m := NewMap(10, 2)
	m[0][0] = example{"1"}
	m[5][1] = other{"plop"}

	ex1 := m[0][0].(example)
	ex2 := m[5][1].(other)
	if ex1.name != "1" || ex2.plop != "plop" {
		t.Fail()
	}

}

type example struct {
	name string
}

type other struct {
	plop string
}
