package engine

import (
	"github.com/woodchuckchoi/KVDB/src/engine/memtable"
	memtable "github.com/woodchuckchoi/KVDB/src/engine/memtable"
	"github.com/woodchuckchoi/KVDB/src/engine/vars"

	// sstable "github.com/woodchuckchoi/KVDB/src/engine/sstable"
	myVar "github.com/woodchuckchoi/KVDB/src/engine/vars"
)

var r int = 3

type Engine struct {
	memTable Memtable
	ssTable  SStable
}

type Memtable interface {
	Put(string, string) error
	Get(string) (string, error)
	Flush() []vars.KeyValue
}

type SStable interface {
	Get(string) (string, error)
	L0Merge([]myVar.KeyValue) error
}

func NewEngine() *Engine {
	return &Engine {
		memTable: memtable.NewMemtable(),
		ssTable:  sstable.,
	}
}
