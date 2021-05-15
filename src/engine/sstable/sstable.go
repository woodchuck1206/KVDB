package sstable

import (

)

type SSTable struct {	
	r			int
	addr	string
	
}

func NewSsTable() *SSTable {

}

func (sstable *SSTable) Get(key string) (string, error) {

}

func 