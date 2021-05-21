package util

import (
	"github.com/woodchuckchoi/KVDB/src/engine/vars"
)

func KeyValueToByteSlice(kv vars.KeyValue) []byte {
	ret := make([]byte, len(kv.Key)+len(kv.Value)+2)
	idx := 0
	for _, val := range []byte(kv.Key) {
		ret[idx] = val
		idx++
	}
	ret[idx] = 2
	idx++
	for _, val := range []byte(kv.Value) {
		ret[idx] = val
		idx++
	}
	ret[idx] = 0
	return ret
}

// func KeyValueSliceToByteSlice(kvs []vars.KeyValue) []byte {

// }

func IntPow(a, b int) int {
	ret := 1
	for i := 0; i < b; i++ {
		ret *= b
	}
	return ret
}
