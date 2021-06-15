package engine

import (
	"github.com/woodchuckchoi/KVDB/src/engine/memtable"
	"github.com/woodchuckchoi/KVDB/src/engine/sstable"
	"github.com/woodchuckchoi/KVDB/src/engine/vars"
)

var r int = 3

type Engine struct {
	memTable  Memtable
	ssTable   SStable
	compactor sstable.Compactor
}

type Memtable interface {
	Put(string, string) error
	Get(string) (string, error)
	Flush() []vars.KeyValue
}

type SStable interface {
	Get(string) (string, error)
	L0Merge([]vars.KeyValue) (int, error)
}

func NewEngine(memTableThresholdSize, r int) *Engine {
	return &Engine{
		memTable:  memtable.NewMemtable(memTableThresholdSize),
		ssTable:   sstable.NewSsTable(r),
		compactor: sstable.NewCompactor(0),
	}
}

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
	targetLevel := -1
	err := this.memTable.Put(key, value)

	if err == vars.MEM_TBL_FULL_ERROR {
		flushedMemtable := this.memTable.Flush()
		targetLevel, err = this.ssTable.L0Merge(flushedMemtable) // should keep the old memtable till the job finishes
	}

	if err == vars.SS_TBL_LVL_FULL_ERROR {
		mergeSignal := sstable.MergeSignal {
			Level: targetLevel,
			// compaction should occur recursively too
			
		}
		this.compactor.Receiver <- 
	}

	return err
}
