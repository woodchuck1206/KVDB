package engine

import (
	"fmt"

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
	Show() []vars.KeyValue
}

type SStable interface {
	Get(string) (string, error)
	// L0Merge([]vars.KeyValue) (int, error)
	Merge(int, []vars.KeyValue) (int, error)
	MergeBlock(int, sstable.Block) (int, error)
	GetSelf() *sstable.SSTable
}

func NewEngine(memTableThresholdSize, r int) *Engine {
	return &Engine{
		memTable:  memtable.NewMemtable(memTableThresholdSize),
		ssTable:   sstable.NewSsTable(r),
		compactor: sstable.NewCompactor(0),
	}
}

func (this *Engine) Compact(level int) sstable.Block {
	mergedBlockReceiver := make(chan sstable.Block)
	mergeSignal := sstable.MergeSignal{
		Level:    level,
		LevelRef: this.ssTable.GetSelf(),
		Returner: mergedBlockReceiver,
	}

	this.compactor.Receive(mergeSignal)

	var mergedBlock sstable.Block

	select {
	case received := <-mergedBlockReceiver:
		mergedBlock = received
	}
	return mergedBlock
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

	if err == nil && value == vars.TOMBSTONE {
		err = vars.GET_FAIL_ERROR
	}

	return value, err
}

func (this *Engine) Put(key, value string) error {
	targetLevel := -1
	err := this.memTable.Put(key, value)

	if err == vars.MEM_TBL_FULL_ERROR {
		fmt.Println("MEM TABLE FULL!")
		flushedMemtable := this.memTable.Flush()
		// targetLevel, err = this.ssTable.L0Merge(flushedMemtable) // should keep the old memtable till the job finishes
		targetLevel, err = this.ssTable.Merge(0, flushedMemtable)
	}

	for err == vars.SS_TBL_LVL_FULL_ERROR {
		fmt.Println("SS TABLE FULL!")
		mergedBlock := this.Compact(targetLevel) // do it sequentially, refactor it to run concurrent later
		targetLevel++
		targetLevel, err = this.ssTable.MergeBlock(targetLevel, mergedBlock) // does it recursively
		// targetLevel, err = this.ssTable.Merge(targetLevel, )
	}

	return err
}

func (this *Engine) Delete(key string) error {
	return this.Put(key, vars.TOMBSTONE)
}

func (this *Engine) Status() {
	fmt.Println(this.memTable.Show())
	this.ssTable.GetSelf().Status()
}
