package tests

import (
	"testing"

	"github.com/woodchuckchoi/KVDB/src/engine/sstable"
)

func TestBlock(t *testing.T) {
	fileName := ""
	// index := []vars.SparseIndex{{"abo", 25}, {"jds", 1059}, {"rpa", 2098}, {"xpm", 3136}, {"zxz", 3421}}
	b := sstable.Block{
		FileName: fileName,
		Index:    index,
		Size:     3437,
	}

	val, err := b.Get("jds")
	t.Logf("VAL %v ERR %v", val, err)
}
