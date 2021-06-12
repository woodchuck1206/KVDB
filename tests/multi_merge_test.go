package tests

import (
	"fmt"
	"math/rand"
	"sort"
	"testing"

	"github.com/woodchuckchoi/KVDB/src/engine/sstable"
	"github.com/woodchuckchoi/KVDB/src/engine/util"
	"github.com/woodchuckchoi/KVDB/src/engine/vars"
)

func generateRandomString(length int) string {
	ret := make([]byte, length)
	for i := 0; i < length; i++ {
		ret[i] = generateRandomAlphabet()
	}
	return string(ret)
}

func generateRandomAlphabet() byte {
	return 'a' + byte(rand.Intn(26))
}

func TestCompaction(t *testing.T) {

	block1KeyValue := []vars.KeyValue{}
	block2KeyValue := []vars.KeyValue{}
	block3KeyValue := []vars.KeyValue{}

	testBlocks := []*sstable.Block{}
	for idx, block := range [][]vars.KeyValue{block1KeyValue, block2KeyValue, block3KeyValue} {
		for i := 0; i < 100; i++ {
			key, value := generateRandomString(20), generateRandomString(50)
			block = append(block, vars.KeyValue{
				Key:   key,
				Value: value,
			})
		}

		sort.Slice(block, func(i, j int) bool {
			if block[i].Key < block[j].Key {
				return true
			}
			return false
		})
		byteSlice, sparseIndex := util.KeyValueSliceToByteSliceAndSparseIndex(block)
		fileName := fmt.Sprintf("block%v.data", idx)
		util.WriteByteSlice(fileName, byteSlice)

		t.Logf("Wrote %v\n", fileName)

		block := sstable.Block{
			FileName: fileName,
			Index:    sparseIndex,
			Size:     len(byteSlice),
		}

		testBlocks = append(testBlocks, &block)
	}

	testLevel := sstable.Level{
		Blocks: testBlocks,
	}

	mergedBlock := sstable.MultiMerge(&testLevel, 1)
	t.Errorf("%v\n%v\n%v", mergedBlock.FileName, mergedBlock.Index, mergedBlock.Size)
}
