package sstable

import (
	"fmt"

	"github.com/bits-and-blooms/bloom/v3"
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
	Bloom    bloom.BloomFilter // bloom filter added
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

func (this *SSTable) GetSelf() *SSTable {
	return this
}

func (this *SSTable) L0Merge(keyValuePairs []vars.KeyValue) (int, error) {
	// order := len(this.levels[0].blocks)
	// fileName := util.GenerateFileName(0, order)
	return this.Merge(0, keyValuePairs)
	// util.WriteKeyValuePairs()
	// return nil
}

func (this *SSTable) Cleanse(level int) {
	if level >= 0 && level < len(this.levels) {
		for _, block := range this.levels[level].Blocks {
			util.RemoveFile(block.FileName)
		}
		this.levels[level].Blocks = []*Block{}
	}
}

func (this *SSTable) MergeBlock(level int, block Block) (int, error) {
	if level == len(this.levels) {
		newLevel := &Level{
			Blocks: []*Block{&block},
		}
		this.levels = append(this.levels, newLevel)
	} else {
		this.levels[level].Blocks = append([]*Block{&block}, this.levels[level].Blocks...)
	}

	// this.levels[level].Blocks = append(this.levels[level].Blocks, &block)
	if len(this.levels[level].Blocks) == util.GetMaxBlockSizeOfLevel(level, this.r) {
		return level, vars.SS_TBL_LVL_FULL_ERROR
	}
	return -1, nil
}

func (this *SSTable) Merge(level int, keyValuePairs []vars.KeyValue) (int, error) {
	if len(this.levels) <= level {
		this.levels = append(this.levels, &Level{
			Blocks: []*Block{},
		})
	}
	fileName := util.GenerateFileName(level)
	fullPath := util.GetFullPathOf(level, fileName)
	byteSlice, sparseIndex := util.KeyValueSliceToByteSliceAndSparseIndex(keyValuePairs)
	if util.WriteByteSlice(fullPath, byteSlice) != nil {
		return -1, vars.FILE_CREATE_ERROR
	}

	newBlock := Block{
		FileName: fullPath,
		Index:    sparseIndex,
		Size:     len(byteSlice),
	}

	this.levels[level].Blocks = append([]*Block{&newBlock}, this.levels[level].Blocks...)
	// this.levels[level].Blocks = append(this.levels[level].Blocks, &Block{
	// 	FileName: fullPath,
	// 	Index:    sparseIndex,
	// 	Size:     len(byteSlice),
	// })

	if len(this.levels[level].Blocks) == util.GetMaxBlockSizeOfLevel(level, this.r) {
		return level, vars.SS_TBL_LVL_FULL_ERROR // compaction should kick in
	}
	return -1, nil
}

func (this *SSTable) Status() {
	for idx, level := range this.levels {
		fmt.Printf("Level %v [", idx)
		for _, b := range level.Blocks {
			fmt.Printf("%v, ", b.FileName)
		}
		fmt.Printf("]\n")
		fmt.Printf("Index %v [", idx)
		for _, b := range level.Blocks {
			fmt.Printf("%v, ", b.Index)
		}
		fmt.Printf("]\n")
	}
}

func (this *SSTable) CleanAll() {
	CleanAll(this)
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
