package engine

import (

	// sstable "github.com/woodchuckchoi/KVDB/src/engine/sstable"
	// memtable "github.com/woodchuckchoi/KVDB/src/engine/memtable"
	myVar "github.com/woodchuckchoi/KVDB/src/engine/vars"
)

const r int = 3

type Engine struct {
	memTable Memtable
	ssTable  SStable
}

type Memtable interface {
	Put(key, value string) error
	Get(key string) (string, error)
	Flush() []myVar.KeyValue
}

type SStable interface {
	Get(key string) (string, error)
}

func ImportTest() string {
	return "Engine"
}
