package sstable

import (
	"github.com/woodchuckchoi/KVDB/src/engine/util"
	"github.com/woodchuckchoi/KVDB/src/engine/vars"
)

type SSTable struct {
	r      int
	levels []*Level
}

type Level struct {
	Blocks []*Block
}

type Block struct {
	FileName string
	Index    []vars.SparseIndex
	Size     int
}

func (this *Block) Get(key string) (string, error) {
	from, till := 0, -1

	// get range from sparseIndex
	for _, keyOffsetPair := range this.Index {
		if key >= keyOffsetPair.Key {
			from = keyOffsetPair.Offset
		}
		if key < keyOffsetPair.Key {
			till = keyOffsetPair.Offset
			break
		}
	}

	keyValuePairs, err := util.ReadKeyValuePairs(this.FileName, from, till)
	if err != nil {
		return "", err
	}

	return BinarySearchKeyValuePairs(keyValuePairs, key)
}

func (this *Block) has(key string) bool {
	return key >= this.Index[0].Key && key <= this.Index[len(this.Index)-1].Key
}

func BinarySearchKeyValuePairs(binTree []vars.KeyValue, key string) (string, error) {
	left, right := 0, len(binTree)
	for left < right {
		mid := (left + right) / 2
		if binTree[mid].Key == key {
			return binTree[mid].Value, nil
		}

		if binTree[mid].Key < key {
			left = mid + 1
		} else {
			right = mid
		}
	}
	return "", vars.KEY_NOT_FOUND_ERROR
}

// leveling
func (this *SSTable) Get(key string) (string, error) {
	for _, level := range this.levels {
		for blockIdx := 0; blockIdx < len(level.Blocks); blockIdx++ {
			if level.Blocks[blockIdx].has(key) {
				val, err := level.Blocks[blockIdx].Get(key)
				if err == nil {
					return val, err
				}
			}
		}
	}
	return "", vars.GET_FAIL_ERROR
}

func (this *SSTable) L0Merge(keyValuePairs []vars.KeyValue) error {
	// order := len(this.levels[0].blocks)
	// fileName := util.GenerateFileName(0, order)
	return this.merge(0, keyValuePairs)
	// util.WriteKeyValuePairs()
	// return nil
}

func (this *SSTable) merge(level int, keyValuePairs []vars.KeyValue) error {
	order := 0
	if len(this.levels) > level {
		order = len(this.levels[level].Blocks)
	} else {
		this.levels = append(this.levels, &Level{
			Blocks: []*Block{},
		})
	}

	fileName := util.GenerateFileName(level, order)
	byteSlice, sparseIndex := util.KeyValueSliceToByteSliceAndSparseIndex(keyValuePairs)
	if util.WriteByteSlice(fileName, byteSlice) != nil {
		return vars.FILE_CREATE_ERROR
	}

	this.levels[level].Blocks = append(this.levels[level].Blocks, &Block{
		FileName: fileName,
		Index:    sparseIndex,
		Size:     len(byteSlice),
	})
	// if the level is full, compaction should kick in at this point
	return nil
}

func NewSsTable(r int) *SSTable {
	return &SSTable{
		r:      r,
		levels: []*Level{},
	}
}

func CleanAll(ssTable *SSTable) {
	for _, level := range ssTable.levels {
		for _, block := range level.Blocks {
			util.RemoveFile(block.FileName)
		}
	}
}

// tiering
// func (this *SSTable) Get(key string) (string, error) {
// 	for _, level := range this.levels {
// 		blockIdx := 0

// 		for ; blockIdx < len(level.blocks); blockIdx++ {
// 			if blockIdx == len(level.blocks)-1 || level.blocks[blockIdx].index[0].Key > key {
// 				break
// 			}
// 		}

// 		val, err := level.blocks[blockIdx].Get(key)
// 		if err == nil {
// 			return val, err
// 		}
// 	}
// 	return "", vars.GET_FAIL_ERROR
// }
