package memtable

import (
  // rbtree "github.com/woodchuckchoi/KVDB/src/engine/memtable/rbtree"
)

type Memtable struct {
  // redBlackTree rbtree.RBTree
}

type Tree interface {

}

func Put(key, value string) error {

}

func Get(key string) string, error {

}

func Flush() error {
  // after flush (save on disk), reset rbtree
}

func merge() error {

}
