package engine

import (
	"github.com/woodchuckchoi/KVDB/src/engine/memtable"
	"github.com/woodchuckchoi/KVDB/src/engine/sstable"
	"github.com/woodchuckchoi/KVDB/src/engine/vars"
	// sstable "github.com/woodchuckchoi/KVDB/src/engine/sstable"
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
	L0Merge([]vars.KeyValue) error
}

func NewEngine(memTableThresholdSize, r int) *Engine {
	return &Engine{
		memTable: memtable.NewMemtable(memTableThresholdSize),
		ssTable:  sstable.NewSsTable(r),
	}
}

// all funcs should operate concurrently

func (this *Engine) Get(key string) (string, error) {
	var (
		value string
		err   error
	)
	value, err = this.memTable.Get(key)
	if err != nil {
		value, err = this.ssTable.Get(key)
	}
	return value, err
}

func (this *Engine) Put(key, value string) error {
	err := this.memTable.Put(key, value)

	if err == vars.MEM_TBL_FULL_ERROR {
		err = this.ssTable.L0Merge(this.memTable.Flush()) // should keep the old memtable till the job finishes
	}

	return err
}
