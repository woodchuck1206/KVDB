package memtable

import (
	rbtree "github.com/woodchuckchoi/KVDB/src/engine/memtable/rbtree"
	"github.com/woodchuckchoi/KVDB/src/engine/vars"
)

type Memtable struct {
	tree      Tree
	sstb      target
	size      int
	threshold int
}

type target interface {
	Merge([]vars.KeyValue) error
}

type Tree interface {
	Get(key string) (string, error)
	Put(key, value string) error
	Flush() []vars.KeyValue
}

func NewMemtable(t target) *Memtable {
	return &Memtable{
		tree: rbtree.NewTree(),
		size: 0,
		sstb: t,
	}
}

func (memtable *Memtable) Put(key, value string) error {
	// defer func() error {
	// 	var ret error = nil
	// 	if err := recover(); err != nil {
	// 		ret = vars.PUT_FAIL_ERROR
	// 		ret = errors.New("PUT FAILED")
	// 	}
	// 	return ret
	// }()

	if err := memtable.tree.Put(key, value); err != nil {
		return err
	}
	memtable.size += varToSize(key, value)
	if memtable.size >= memtable.threshold {
		memtable.sstb.Merge(memtable.flush())
	}
	return nil
}

func (memtable *Memtable) Get(key string) (string, error) {
	return memtable.tree.Get(key)
}

func (memtable *Memtable) flush() []vars.KeyValue {
	toFlush := memtable.tree.Flush()
	memtable.tree = rbtree.NewTree()
	return toFlush
}

func varToSize(vars ...string) int {
	ret := 0
	for _, v := range vars {
		ret += len(v)
	}
	return ret
}
