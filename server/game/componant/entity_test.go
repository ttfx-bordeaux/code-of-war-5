package componant

import "testing"

func TestIdentifiableTower(t *testing.T) {
	to := Tower{id: 12}
	if to.ID() != to.id {
		t.Fail()
	}
}
