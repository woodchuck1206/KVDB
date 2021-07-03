package tests

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"sort"
	"testing"
	"time"

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
	rand.Seed(int64(time.Now().Second()))
	defer func() {
		for i := 0; i < 3; i++ {
			fileName := fmt.Sprintf("block%v.data", i)
			os.Remove(fileName)
		}
	}()

	block1KeyValue := []vars.KeyValue{}
	block2KeyValue := []vars.KeyValue{}
	block3KeyValue := []vars.KeyValue{}

	testBlocks := []*sstable.Block{}
	for idx, block := range []*[]vars.KeyValue{&block1KeyValue, &block2KeyValue, &block3KeyValue} {
		for i := 0; i < 100; i++ {
			key, value := generateRandomString(20), generateRandomString(50)
			*block = append(*block, vars.KeyValue{
				Key:   key,
				Value: value,
			})
		}

		sort.Slice(*block, func(i, j int) bool {
			if (*block)[i].Key < (*block)[j].Key {
				return true
			}
			return false
		})
		byteSlice, sparseIndex := util.KeyValueSliceToByteSliceAndSparseIndex(*block)
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

	t.Log("MergeUnits Ready")
	for i := 0; i < len(testLevel.Blocks); i++ {
		t.Logf("%v\n", testLevel.Blocks[i].FileName)
	}

	mergedBlock := sstable.MultiMerge(&testLevel, 1, 3)
	defer os.Remove(mergedBlock.FileName)
	mergedFile, _ := os.Open(mergedBlock.FileName)
	bytes, _ := ioutil.ReadAll(mergedFile)
	mergedKeyValues := util.ByteSliceToKeyValue(bytes)

	compareKeyValues := []vars.KeyValue{}
	for _, block := range [][]vars.KeyValue{block1KeyValue, block2KeyValue, block3KeyValue} {
		t.Log(len(block))
		compareKeyValues = append(compareKeyValues, block...)
	}

	sort.Slice(compareKeyValues, func(i, j int) bool {
		if compareKeyValues[i].Key < compareKeyValues[j].Key {
			return true
		} else if compareKeyValues[i].Key == compareKeyValues[j].Key && i < j {
			return true
		}
		return false
	})

	var before vars.KeyValue = compareKeyValues[0]
	overwritten := []vars.KeyValue{before}
	for i := 1; i < len(compareKeyValues); i++ {
		if compareKeyValues[i].Key == before.Key {
			continue
		}
		before = compareKeyValues[i]
		overwritten = append(overwritten, before)
	}

	if len(overwritten) != len(mergedKeyValues) {
		t.Error("Length Not Match", len(overwritten), len(mergedKeyValues))
	}

	t.Logf("Merged Index: %v\nMerged Size: %v\nBytes Size: %v\n", mergedBlock.Index, mergedBlock.Size, len(bytes))

	for i := 0; i < len(compareKeyValues); i++ {
		if overwritten[i] != mergedKeyValues[i] {
			t.Errorf("\n%vth Comparison\nOriginal: %v\nMerged: %v", i, overwritten[i], mergedKeyValues[i])
		}
	}

}
