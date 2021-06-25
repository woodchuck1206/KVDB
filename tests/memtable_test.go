package tests

import (
	"testing"

	"github.com/woodchuckchoi/KVDB/src/engine/memtable"
)

func TestMemTable(t *testing.T) {
	m := memtable.NewMemtable(1024)

	err := m.Put("a", "a string")
	if err != nil {
		t.Error("FAIL INIT")
	}
	value, err := m.Get("a")
	if err != nil {
		t.Error("FAIL PUTTING")
	}
	if value != "a string" {
		t.Error("FAIL FORMATTING")
	}

	_, err = m.Get("b")
	if err == nil {
		t.Error("FAIL GETTING")
	}

	err = m.Put("a", "new a string")
	if err != nil {
		t.Error("FAIL PUTTING")
	}

	value, err = m.Get("a")
	if err != nil {
		t.Error("FAIL GETTING")
	}
	if value != "new a string" {
		t.Error("FAIL PUTTING")
	}
}
