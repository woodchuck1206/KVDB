package rbtree

import (
	"fmt"
	"testing"
)

func TestRedBlackTree(t *testing.T) {
	rbtree := &RedBlackTree{}

	type toInsert struct {
		key   string
		value string
	}

	testInsert := []toInsert{
		{key: "ab", value: "2r3"},
		{key: "asdf", value: "adg"},
		{key: "fwe", value: "zb"},
		{key: "qewf", value: "asGD"},
		{key: "abd", value: "WYHR"},
		{key: "afng", value: "sfh"},
	}

	for _, val := range testInsert {
		rbtree.Insert(val.key, val.value)
	}

	fmt.Println(rbtree.Flush())
}
