package sstable

import (
	"os"

	"github.com/woodchuckchoi/KVDB/src/engine/vars"
)

type SSTable struct {
	r      int
	tables [][]Table
}

type Table struct { // stored on disk
	ptr *os.File
}

func NewSsTable(r int) *SSTable {
	return &SSTable{
		r:      r,
		tables: [][]Table{},
	}
}

func (sstable *SSTable) Get(key string) (string, error) {
	for _, level := range sstable.tables {
		for _, table := range level {
			if val, err := table.Get(key); err == nil {
				return val, nil
			}
		}
	}
	return "", vars.GET_FAIL_ERROR
}

func (table Table) Get(key string) (string, error) {

}
