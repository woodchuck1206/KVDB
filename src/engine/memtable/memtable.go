package memtable

import (
	"github.com/woodchuckchoi/KVDB/src/engine/memtable/rbtree"
	"github.com/woodchuckchoi/KVDB/src/engine/vars"
)

type Memtable struct {
	tree      Tree
	size      int
	threshold int
}

type Tree interface {
	Get(key string) (string, error)
	Put(key, value string) (error, int)
	Flush() []vars.KeyValue
}

func NewMemtable(threshold int) *Memtable {
	return &Memtable{
		tree:      rbtree.NewTree(), // temporary place to keep the original tree while flushing needed?
		size:      0,
		threshold: threshold,
	}
}

func (this *Memtable) Put(key, value string) error {
	err, sizeChange := this.tree.Put(key, value)
	if err != nil {
		return err
	}

	this.size += sizeChange
	if this.size >= this.threshold { //
		return vars.MEM_TBL_FULL_ERROR
		// memtable.sstb.Merge(memtable.flush())
	}
	return nil
}

func (this *Memtable) Get(key string) (string, error) {
	return this.tree.Get(key)
}

func (this *Memtable) Flush() []vars.KeyValue {
	toFlush := this.tree.Flush()
	this.reborn()
	return toFlush
}

func (this *Memtable) Show() []vars.KeyValue {
	return this.tree.Flush()
}

func (this *Memtable) reborn() {
	this.tree = rbtree.NewTree()
	this.size = 0
}
