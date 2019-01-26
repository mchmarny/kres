package util

import (
	"testing"
)

func TestID(t *testing.T) {

	id1 := MakeUUID()
	id2 := MakeUUID()

	if id1 == id2 {
		t.Errorf("Non unique IDs")
	}

	if len(id1) != len(id1) {
		t.Errorf("Inconsistent ID length")
	}

}
