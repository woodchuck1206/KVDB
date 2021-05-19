package sstable

import (
	"fmt"
	"os"
	"path"

	"github.com/woodchuckchoi/KVDB/src/engine/util"
	"github.com/woodchuckchoi/KVDB/src/engine/vars"
)

var (
	BASE_DIR   = "/tmp/gokvdb"
	INDEX_TERM = 1 << 12
)

type SSTable struct {
	flushSize int
	r         int
	levels    []*Level
}

type Level struct {
	tables []Table
	index  [][]sIndex
}

type Table struct { // stored on disk
	ptr         *os.File
	sparseIndex []sIndex
}

type sIndex struct {
	key    string
	offset int
}

func NewSsTable(r, flushSize int) *SSTable {
	return &SSTable{
		flushSize: flushSize,
		r:         r,
		levels:    []*Level{},
	}
}

func newLevel(r int) *Level {
	return &Level{
		tables: []Table{},
		index:  [][]sIndex{},
	}
}

func newTable(l, o, size int, data []vars.KeyValue) (*Table, error) {
	fileName := generateTableFileName(l, o)
	f, err := os.Create(fileName)
	if err != nil {
		return nil, vars.FILE_CREATE_ERROR
	}

	toWrite, sparseIndex := parseKeyValue(data, size)
	_, err = f.Write(toWrite)
	if err != nil {
		return nil, vars.FILE_WRITE_ERROR
	}

	return &Table{
		ptr:         f,
		sparseIndex: sparseIndex,
	}, nil
}

func parseKeyValue(data []vars.KeyValue, size int) ([]byte, []sIndex) {
	bdta := make([]byte, size)
	sidx := []sIndex{}
	offset, lastOffset := 0, 0

	for _, pair := range data {
		byteString := util.KeyValueToByteSlice(pair)
		byteCopyHelper(byteString, &bdta, offset)

		if offset-lastOffset >= INDEX_TERM || offset == 0 {
			sidx = append(sidx, sIndex{
				key:    pair.Key,
				offset: offset,
			})
			lastOffset = offset
		}
		offset += len(byteString)
	}

	return bdta, sidx
}

func byteCopyHelper(src []byte, dest *[]byte, offset int) {
	for idx, byteCharacter := range src {
		if len(*dest) < len(src)-idx+offset {
			*dest = append(*dest, byteCharacter)
		} else {
			(*dest)[offset+idx] = byteCharacter
		}
	}
}

func generateTableFileName(level, order int) string {
	fileName := fmt.Sprintf("db:%d:%d.db", level, order)
	return path.Join(BASE_DIR, fileName)
}

func (sstable *SSTable) Get(key string) (string, error) {
	for _, level := range sstable.levels {
		partition := 0
		for partition < len(level.index)-1 {
			if key >= level.index[partition][0].key && key < level.index[partition+1][0].key {
				break
			}
			partition++
		}
		val, err := level.tables[partition].Get(key)
		if err == nil {
			return val, nil
		}
	}
	return "", vars.GET_FAIL_ERROR
}

func (table Table) Get(key string) (string, error) {
	partition := 0
	from, till := -1, -1
	for partition < len(table.sparseIndex)-1 {
		if key >= table.sparseIndex[partition].key && key < table.sparseIndex[partition+1].key {
			from = table.sparseIndex[partition].offset
			break
		}
		partition++
	}

	if partition != len(table.sparseIndex)-1 {
		till = table.sparseIndex[partition+1].offset
	}

	return GetFromFile(table.ptr, from, till, key)
}

func GetFromFile(f *os.File, from, till int, key string) (string, error) {

}

// func loadChunk(f *os.File, from, till int) ([]byte, error) {
// 	_, err := f.Seek(int64(from), 0)
// 	if err != nil {
// 		return nil, vars.FILE_READ_ERROR
// 	}
// 	f.Close()
// 	ret := make([]byte, )
// 	f.Read()
// }
