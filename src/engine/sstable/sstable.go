package sstable

import (
	"fmt"
	"os"

	"github.com/woodchuckchoi/KVDB/src/engine/vars"
)

type SSTable struct {
	r      int
	levels []*Level
}

type Level struct {
	tables []*Table
	index  []string
}

type Table struct { // stored on disk
	ptr *os.File
}

func NewSsTable(r int) *SSTable {
	return &SSTable{
		r:      r,
		levels: []*Level{},
	}
}

func newLevel(r int) *Level {
	return &Level{
		tables: []*Table{},
		index:  []string{},
	}
}

func newTable(l, o int, data []vars.KeyValue) *Table {
	fileName := generateTableFileName(l, o)
}

func generateTableFileName(level, order int) string {
	return fmt.Sprintf("db:level:order.db", level, order)
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
