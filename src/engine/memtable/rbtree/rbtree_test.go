package rbtree

import (
	"reflect"
	"sort"
	"testing"

	"github.com/woodchuckchoi/KVDB/src/engine/vars"
)

func TestRedBlackTree(t *testing.T) {
	rbtree := &RedBlackTree{}

	type toInsert struct {
		key   string
		value string
	}

	testInsert := []vars.KeyValue{
		{Key: "ab", Value: "2r3"},
		{Key: "asdf", Value: "adg"},
		{Key: "fwe", Value: "zb"},
		{Key: "qewf", Value: "asGD"},
		{Key: "abd", Value: "WYHR"},
		{Key: "afng", Value: "sfh"},
	}

	for _, val := range testInsert {
		rbtree.Insert(val.Key, val.Value)
	}

	sort.Slice(testInsert, func(i, j int) bool {
		if testInsert[i].Key < testInsert[j].Key {
			return true
		}
		return false
	})

	if !reflect.DeepEqual(rbtree.Flush(), testInsert) {
		panic("RBTREE ORDER ERROR")
	}
}
