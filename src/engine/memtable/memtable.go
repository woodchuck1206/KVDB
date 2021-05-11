package memtable

import (
  
)

type Memtable struct {
  // redBlackTree Tree
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
