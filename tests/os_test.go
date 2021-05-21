package tests

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestRW(t *testing.T) {
	key := "alarm:test"
	value := "\"some\nValue that needs to be stored!\""
	data := keyValueStore(key, value)
	err := ioutil.WriteFile("test.data", data, 0777)
	defer func() {
		os.Remove("test.data")
	}()
	if err != nil {
		fmt.Println(err)
		return
	}

	read, err := ioutil.ReadFile("test.data")
	if err != nil {
		t.Error("FILE READ ERROR")
		return
	}

	k, v := storeToKeyValue(read)
	if k != key || v != value {
		t.Error("KV PARSE FAIL")
	}
}

func keyValueStore(key, value string) []byte {
	ret := make([]byte, len(key)+len(value)+2)
	idx := 0
	for _, val := range []byte(key) {
		ret[idx] = val
		idx++
	}
	ret[idx] = 3
	idx++
	for _, val := range []byte(value) {
		ret[idx] = val
		idx++
	}
	ret[idx] = 0
	return ret
}

func storeToKeyValue(b []byte) (string, string) {
	var (
		key   string
		value string
	)
	for i, val := range b {
		if val == 0 {
			key = string(b[:i])
			value = string(b[i+1 : len(b)])
			break
		}
	}
	return key, value
}
