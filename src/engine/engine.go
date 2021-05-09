package engine

import (
  // memtable "github.com/woodchuckchoi/KVDB/src/engine/memtable"
)

type Engine struct {
  mt  Memtable
  ss
}

type Memtable interface {
  Put(key, value string) error
  Get(key string) (string, error)
  Flush() error
}

type SStable interface [
  Get(key string) string, error
}

func ImportTest() string {
  return "Engine"
}
