package sstable

import (
	"fmt"
	"os"

	"github.com/woodchuckchoi/KVDB/src/engine/util"
	"github.com/woodchuckchoi/KVDB/src/engine/vars"
)

type Compactor struct {
	Receiver chan MergeSignal
}

type MergeSignal struct {
	Level    int
	LevelRef *SSTable
	Returner chan<- Block
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

func NewCompactor(chanBuffer int) Compactor {
	if chanBuffer <= 0 {
		chanBuffer = 42 // arbitrary default number for MergeSignal buffer size
	}
	channel := make(chan MergeSignal, chanBuffer)
	compactor := Compactor{Receiver: channel}
	compactor.Run()
	return compactor
}

func (this Compactor) Run() {
	go func() {
		for {
			select {
			case mergeSignal := <-this.Receiver:
				mergedBlock := MultiMerge(mergeSignal.LevelRef.levels[mergeSignal.Level], mergeSignal.Level)
				mergeSignal.Returner <- mergedBlock
				close(mergeSignal.Returner)
			}
		}
	}()
}

func (this Compactor) Receive(mergeSignal MergeSignal) {
	this.Receiver <- mergeSignal
}

func MultiMerge(level *Level, l int) Block {
	mergeUnits := []MergeUnit{}
	mergeSparseIndex := []vars.SparseIndex{}
	offsetBefore := -1
	mergeSize := 0
	indexTerm := 1024

	for _, block := range level.Blocks {
		file, err := os.Open(block.FileName)
		defer file.Close()

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

	fileName := util.GenerateFileName(l)
	fullPath := util.GetFullPathOf(l, fileName)
	writeFD, err := os.Create(fullPath)
	if err != nil {
		fmt.Println("NOT ABLE TO MAKE A FILE!", err)
		// error handling // create a new filename?
	}
	defer writeFD.Close()

	var kvToAdd vars.KeyValue
	// multi-merge
	for {
		var unitWithSmallestKeyValue *MergeUnit = nil
		idx := 0
		for {
			if idx >= len(mergeUnits) {
				break
			}

			keyValue, err := mergeUnits[idx].Get()
			if err != nil { // empty mergeUnit
				mergeUnits = append(mergeUnits[:idx], mergeUnits[idx+1:]...)
				continue
			}

			// earlier blocks are more recent
			// compaction overwrite
			// tombstone delete
			// should be implemented

			if unitWithSmallestKeyValue != nil {
				curSmallestKeyValue, _ := unitWithSmallestKeyValue.Get()
				if curSmallestKeyValue.Key > keyValue.Key {
					unitWithSmallestKeyValue = &(mergeUnits[idx])
				}
			} else { // unitWithSmallestKeyValue == nil
				unitWithSmallestKeyValue = &mergeUnits[idx]
			}
			idx++
		}

		if unitWithSmallestKeyValue == nil { // no more merge units to process
			break
		}

		kvToAdd, _ = unitWithSmallestKeyValue.Pop()
		byteKV := util.KeyValueToByteSlice(kvToAdd)
		mergeSize += len(byteKV)

		if offsetBefore == -1 || mergeSize-offsetBefore >= indexTerm {
			mergeSparseIndex = append(mergeSparseIndex, vars.SparseIndex{
				Key:    kvToAdd.Key,
				Offset: mergeSize,
			})

			offsetBefore = mergeSize
		}
		writeFD.Write(byteKV)
	}

	if mergeSparseIndex[len(mergeSparseIndex)-1].Key != kvToAdd.Key {
		byteKV := util.KeyValueToByteSlice(kvToAdd)
		mergeSparseIndex = append(mergeSparseIndex, vars.SparseIndex{
			Key:    kvToAdd.Key,
			Offset: mergeSize - len(byteKV),
		})
	}

	return Block{
		FileName: fullPath,
		Index:    mergeSparseIndex,
		Size:     mergeSize,
	}
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
