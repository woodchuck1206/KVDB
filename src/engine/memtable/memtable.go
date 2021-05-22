package memtable

import (
	rbtree "github.com/woodchuckchoi/KVDB/src/engine/memtable/rbtree"
	"github.com/woodchuckchoi/KVDB/src/engine/vars"
)

const BYTE_FORMAT_AUX int = 2

type Memtable struct {
	tree      Tree
	size      int
	threshold int
}

type Tree interface {
	Get(key string) (string, error)
	Put(key, value string) error
	Flush() []vars.KeyValue
}

func NewMemtable(threshold int) *Memtable {
	return &Memtable{
		tree:      rbtree.NewTree(), // temporary place to keep the original tree while flushing needed?
		size:      0,
		threshold: threshold,
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
	if memtable.size >= memtable.threshold { //
		return vars.MEM_TBL_FULL_ERROR
		// memtable.sstb.Merge(memtable.flush())
	}
	return nil
}

func (memtable *Memtable) Get(key string) (string, error) {
	return memtable.tree.Get(key)
}

func (memtable *Memtable) Flush() []vars.KeyValue {
	toFlush := memtable.tree.Flush()
	memtable.tree = rbtree.NewTree()
	return toFlush
}

// func (memtable *Memtable) FlushV2() []byte {
// 	toFlush := memtable.tree.Flush()
// 	ret := make([]byte, memtable.threshold+BYTE_FORMAT_AUX*len(toFlush))

// 	idx := 0
// 	for _, keyVal := range toFlush {
// 		for _, letter := range []byte(keyVal.Key) {
// 			ret[idx] = letter
// 			idx++
// 		}
// 		ret[idx] = '\a'
// 		idx++
// 		for _, letter := range []byte(keyVal.Value) {
// 			ret[idx] = letter
// 			idx++
// 		}
// 		ret[idx] = '\n'
// 		idx++
// 	}
// 	return ret[:idx]
// }

func varToSize(vars ...string) int {
	ret := 0
	for _, v := range vars {
		ret += len(v)
	}
	return ret
}
