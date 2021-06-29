package tests

import (
	"math/rand"
	"testing"
	"time"

	"github.com/woodchuckchoi/KVDB/src/engine"
	"github.com/woodchuckchoi/KVDB/src/engine/vars"
)

type Log struct {
	i     int
	key   string
	value string
}

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
		switch curActionSelector {
		case 0:
			err = e.Put(key, value)
			if err != nil {
				t.Error(err)
			}
			t.Logf("%05d PUT KEY: %v VALUE: %v\n", i, key, value)
			freshValue, err := e.Get(key)
			if err != nil {
				t.Error(err)
			}
			if freshValue != value {
				t.Logf("FreshValue Mismatch! VALUE SHOULD BE %v, BUT %v\n", value, freshValue)
				e.Status()
				t.FailNow()
			}
			record[key] = value
			break
		case 1:
			valueFromEngine, err = e.Get(key)
			t.Logf("%05d GET KEY: %v VALUE: %v\n", i, key, valueFromEngine)
			valueFromMap, ok := record[key]
			if err == vars.GET_FAIL_ERROR && ok {
				e.Status()
				t.Errorf("%vth run! %v should exist in DB! EngineValue: %v MapValue: %v\n", i, key, valueFromEngine, valueFromMap)
				t.FailNow()
			}
			if err != vars.GET_FAIL_ERROR && valueFromEngine != valueFromMap {
				e.Status()
				t.Errorf("%v value does not match\nDB: %v\nMAP: %v\n", key, valueFromEngine, valueFromMap)
				t.FailNow()
			}
			break
		case 2:
			t.Logf("%05d DEL KEY: %v \n", i, key)
			err = e.Delete(key)
			delete(record, key)
			break
		}

	}
	if !t.Failed() {
		e.CleanAll()
	}
}

// func TestCompaction(t *testing.T) {

// }
