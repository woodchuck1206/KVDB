package sstable

import (
	"fmt"
	"os"

	"github.com/woodchuckchoi/KVDB/src/engine/util"
	"github.com/woodchuckchoi/KVDB/src/engine/vars"
)

var (
	BASE_DIR		= "/tmp/gokvdb"
	INDEX_TERM	= 1 << 12
)

type SSTable struct {
	r      int
	levels []*Level
}

type Level struct {
	tables []Table
	index  [][]sIndex
}

type Table struct { // stored on disk
	ptr 					*os.File
	sparseIndex		[]sIndex
}

type sIndex struct {
	key			string
	offset	int
}

func NewSsTable(r int) *SSTable {
	return &SSTable{
		r:      r,
		levels: []*Level{},
	}
}

func newLevel(r int) *Level {
	return &Level{
		tables: []Table{},
		index:  [][]sIndex{},
	}
}

func newTable(l, o int, data []vars.KeyValue) (*Table, err) {
	fileName := generateTableFileName(l, o)
	f, err := os.Create(fileName)
	if err != nil {
		return nil, vars.FILE_CREATE_ERROR
	}

	// merge these two functions, so it can be done in one loop
	f.Write(util.KeyValueSliceToByteSlice(data))
	sparseIndex := createSparseIndex(data)
	//

	return &Table{
		ptr: f,
		sparseIndex: sparseIndex,
	}
}

func createSparseIndex(data []vars.KeyValue) []sIndex {
	ret := []sIndex {
		sIndex{
			key: data[0].Key,
			offset: 0,
		},
	}
	offset := len(data[0].Key) + len(data[0].Value) + 2
	lastOffset := 0
	
	for idx, pair := range data[1:] {
		if offset >= lastOffset + INDEX_TERM {
			ret = append(ret, sIndex{
				key: pair.Key,
				offset: offset,
			})
			lastOffset = offset
		}
		offset += len(pair.Key) + len(pair.Value) + 2
	}

	return ret
}

func generateTableFileName(level, order int) string {
	return fmt.Sprintf("%s/db:%s:%s.db", BASE_DIR, level, order)
}

func (sstable *SSTable) Get(key string) (string, error) {
	for _, level := range sstable.levels { // index matching
		// for _, table := range level.tables {
		// 	if val, err := table.Get(key); err == nil {
		// 		return val, nil
		// 	}
		}
	}
	return "", vars.GET_FAIL_ERROR
}

func (table Table) Get(key string) (string, error) {

	return "", nil
}
