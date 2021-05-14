package engine

import (
	"github.com/woodchuckchoi/KVDB/src/engine/memtable"
	memtable "github.com/woodchuckchoi/KVDB/src/engine/memtable"

	// sstable "github.com/woodchuckchoi/KVDB/src/engine/sstable"
	myVar "github.com/woodchuckchoi/KVDB/src/engine/vars"
)

var r int = 3

type Engine struct {
	memTable Memtable
	ssTable  SStable
}

func NewEngine() *Engine {
	return &Engine {
		memTable: memtable.NewMemtable(),
		ssTable:  sstable.,
	}
}

type Memtable interface {
	Put(string, string) error
	Get(string) (string, error)
}

type SStable interface {
	Get(string) (string, error)
	Merge([]myVar.KeyValue)
}

func ImportTest() string {
	return "Engine"
}
