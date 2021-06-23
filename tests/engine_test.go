package tests

import (
	"math/rand"
	"testing"
	"time"

	"github.com/woodchuckchoi/KVDB/src/engine"
	"github.com/woodchuckchoi/KVDB/src/engine/vars"
)

const (
	keyLength   = 3 // to make it more likely to produce duplicate values
	valueLength = 20
)

func getRandomString(length int) string {
	ret := make([]byte, length)
	for i := 0; i < length; i++ {
		ret[i] = getRandomAlphabet()
	}
	return string(ret)
}

func getRandomAlphabet() byte {
	return 'a' + byte(rand.Intn(26))
}

func TestPutGetDelete(t *testing.T) {
	memtableSize := 1024
	r := 3
	e := engine.NewEngine(memtableSize, r)

	seed := time.Now().Nanosecond()
	rand.Seed(int64(seed))

	record := map[string]string{}
	for i := 0; i < 1000; i++ {
		key := getRandomString(keyLength)
		value := getRandomString(valueLength)

		curActionSelector := rand.Intn(3)
		var err error
		var valueFromEngine string
		t.Logf("%04dth run %v key: %v value: %v\n", i, curActionSelector, key, value)
		switch curActionSelector {
		case 0:
			err = e.Put(key, value)
			record[key] = value
			break
		case 1:
			valueFromEngine, err = e.Get(key)
			valueFromMap, ok := record[key]
			if err == vars.GET_FAIL_ERROR && ok {
				t.Errorf("%v does not exist in DB!\n", key)
			}
			if valueFromEngine != valueFromMap {
				t.Errorf("%v value does not match\nDB: %v\nMAP: %v\n", key, valueFromEngine, valueFromMap)
			}
			break
		case 2:
			err = e.Delete(key)
			delete(record, key)
			break
		}

	}

}

// func TestCompaction(t *testing.T) {

// }
