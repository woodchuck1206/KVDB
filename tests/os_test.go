package tests

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestAppend(t *testing.T) {
	fileName := "anotherTest.data"
	byteWrite := make([]byte, 1)
	cmp := []byte{}
	f, err := os.Create(fileName)
	if err != nil {
		t.Errorf("%v", err)
	}
	// f.Chmod(fs.ModeAppend)

	for i := 0; i < 10; i++ {
		byteWrite[0] = byte(int('a') + i)
		cmp = append(cmp, byte(int('a')+i))
		f.Write(byteWrite) // simple Write function appends to the file by default
	}
	f.Close()

	f, _ = os.Open(fileName)
	ret := make([]byte, 100)
	l, _ := f.Read(ret)
	os.Remove(fileName)
	if l != 10 {
		t.Errorf("length not match %v with values %v \ncmp %v", l, string(ret), string(cmp))
	}

	for i := 0; i < l; i++ {
		if cmp[i] != ret[i] {
			t.Errorf("byte not match at %v", i)
		}
	}
}

func TestReadChunk(t *testing.T) {
	fileName := "testfile.data"
	toWrite := "this will be\nwritten"
	byteWrite := []byte(toWrite)
	err := ioutil.WriteFile(fileName, byteWrite, 0777)
	defer func() {
		os.Remove(fileName)
	}()
	if err != nil {
		t.Error("WRITE FAIL")
	}

	file, err := os.Open(fileName)
	if err != nil {
		t.Error("READ FAIL")
	}
	defer func() {
		file.Close()
	}()

	receiver := make([]byte, 10)
	file.Seek(10, 0)
	n, err := file.Read(receiver)
	if string(receiver[:n]) != toWrite[10:] {
		t.Error("PARSING FAIL")
	}
}

func TestRW(t *testing.T) {
	key := "alarm:test"
	value := "\"some\nValue that needs to be stored!\""
	data := keyValueStore(key, value)
	err := ioutil.WriteFile("test.data", data, 0777)
	defer func() {
		os.Remove("test.data")
	}()
	if err != nil {
		t.Error("WRITE FAIL")
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
		if val == 3 {
			key = string(b[:i])
			value = string(b[i+1 : len(b)-1])
			break
		}
	}
	return key, value
}
