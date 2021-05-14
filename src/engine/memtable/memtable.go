package memtable

import (
	rbtree "github.com/woodchuckchoi/KVDB/src/engine/memtable/rbtree"
	myVar "github.com/woodchuckchoi/KVDB/src/engine/vars"
)

type Memtable struct {
	tree Tree
	size int
}

type Tree interface {
	Get(key string) (string, error)
	Put(key, value string) error
	Flush() []myVar.KeyValue
}

func NewMemtable() *Memtable {
	return &Memtable{
		tree: rbtree.NewTree(),
		size: 0,
	}
}

func (memtable *Memtable) Put(key, value string) error {
	if err := memtable.tree.Put(key, value); err != nil {
		return err
	}
	memtable.size += varToSize(key, value)
	return nil
}

func (memtable *Memtable) Get(key string) (string, error) {

}

func (memtable *Memtable) Flush() []myVar.KeyValue {
	// after flush (save on disk), reset rbtree
}

func merge() error {

}

func varToSize(vars ...string) int {
	ret := 0
	for _, v := range vars {
		ret += len(v)
	}
	return ret
}
