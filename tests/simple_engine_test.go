package tests

import (
	"testing"

	"github.com/woodchuckchoi/KVDB/src/engine"
)

func TestSimple(t *testing.T) {
	memtableSize := 1024
	r := 3
	m := engine.NewEngineWithValues(memtableSize, r)

	err := m.Put("cya", "rzqfwkrzfpembscyugcx")
	if err != nil {
		t.Error("FAIL INIT")
	}

	value, err := m.Get("cya")
	if err != nil {
		t.Error("FAIL PUTTING")
	}
	if value != "rzqfwkrzfpembscyugcx" {
		t.Error("FAIL FORMATTING")
	}

	m.Delete("gbr")

	err = m.Put("kxb", "swwwgblypkouutrlwsjr")
	if err != nil {
		t.Error("FAIL INIT")
	}

	value, err = m.Get("kxb")
	if err != nil {
		t.Error("FAIL PUTTING")
	}
	if value != "swwwgblypkouutrlwsjr" {
		t.Error("FAIL FORMATTING")
	}
	m.Status()
}
