package sstable

import (
	"os"

	"github.com/woodchuckchoi/KVDB/src/engine/util"
	"github.com/woodchuckchoi/KVDB/src/engine/vars"
)

type Compactor struct {
	Receiver <-chan MergeSignal
}

type MergeSignal struct {
	Level    int
	LevelRef *SSTable
}

type MergeUnit struct { // per Block
	KeyValue []vars.KeyValue

	SparseIndex []vars.SparseIndex
	SparseIdx   int // sparse index's index

	Finish bool
	File   *os.File
}

func (this *MergeUnit) Load() {
	// load size should be managed // each key value pairs sizes differ
	var curBufferSize int

	switch {
	case this.SparseIdx+1 == len(this.SparseIndex):
		fileInfo, _ := this.File.Stat()
		curBufferSize = int(fileInfo.Size()) - this.SparseIndex[this.SparseIdx].Offset

	case this.SparseIdx >= len(this.SparseIndex):
		this.Finish = true
		return

	default: // this.SparseIdx < len(this.SparseIndex) - 1
		curBufferSize = this.SparseIndex[this.SparseIdx+1].Offset - this.SparseIndex[this.SparseIdx].Offset
	}

	this.SparseIdx++

	buffer := make([]byte, curBufferSize)

	this.File.Read(buffer)
	this.KeyValue = util.ByteSliceToKeyValue(buffer)
}

func (this *MergeUnit) Get() (vars.KeyValue, error) {
	var ret vars.KeyValue

	if len(this.KeyValue) == 0 {
		this.Load()
	}
	if this.Finish {
		return ret, vars.FILE_EOF_ERROR
	}

	ret = this.KeyValue[0]
	return ret, nil
}

func (this *MergeUnit) Pop() (vars.KeyValue, error) {
	ret, err := this.Get()
	if err == nil {
		this.KeyValue = this.KeyValue[1:]
	}
	return ret, err
}

func NewCompactor(chanBuffer int) (*Compactor, chan<- MergeSignal) {
	if chanBuffer <= 0 {
		chanBuffer = 42 // arbitrary default number for MergeSignal buffer size
	}
	channel := make(chan MergeSignal, chanBuffer)
	return &Compactor{Receiver: channel}, channel
}

func MultiMerge(level *Level, l int) Block {
	mergeUnits := []MergeUnit{}
	mergedKeyValues := []vars.KeyValue{}

	for _, block := range level.Blocks {
		file, err := os.Open(block.FileName)
		if err == nil {
			mergeUnits = append(mergeUnits, MergeUnit{
				KeyValue:    []vars.KeyValue{},
				SparseIndex: block.Index,
				SparseIdx:   0,
				Finish:      false,
				File:        file,
			})
		}
	}

	// multi-merge
	for {
		var unitWithSmallestKeyValue *MergeUnit
		nextMergeUnits := []MergeUnit{}
		for i := 0; i < len(mergeUnits); i++ {
			keyValue, err := mergeUnits[i].Get()
			if err != nil { // EOF
				continue
			}
			nextMergeUnits = append(nextMergeUnits, mergeUnits[i])
			if unitWithSmallestKeyValue != nil {
				curSmallestKeyValue, _ := unitWithSmallestKeyValue.Get()
				if curSmallestKeyValue.Key > keyValue.Key {
					unitWithSmallestKeyValue = &(mergeUnits[i])
				}
			}
		}

		if unitWithSmallestKeyValue == nil {
			break
		}
		toAdd, _ := unitWithSmallestKeyValue.Pop()
		mergedKeyValues = append(mergedKeyValues, toAdd)
		mergeUnits = nextMergeUnits
	}

	// the mergedKeyValues should not keep all records obviously, it should write as it reads
	// util.KeyValueSliceToByteSliceAndSparseIndex()
	// util.WriteKeyValuePairs()
}

// func (this *Compactor) Run() {
// 	for {
// 		select {
// 		case sig := <-this.Receiver:
// 			// merge using sig and update the reference
// 			levelToMerge := sig.Level
// 			sstable := sig.LevelRef

// 			sstable.levels[0]
// 		}
// 	}
// }

// func mergeLevel(level *Level) *Block {

// }

// func multiMerge()
