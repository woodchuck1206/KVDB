package tests

import (
	"testing"

	"github.com/woodchuckchoi/KVDB/src/engine/sstable"
	"github.com/woodchuckchoi/KVDB/src/engine/vars"
)

type checkData struct {
	key     string
	boolean bool
}

func TestSSTable(t *testing.T) {
	tempR := 3
	testTable := sstable.NewSsTable(tempR)
	testData := []vars.KeyValue{
		{Key: "abraham", Value: "Lincoln"},
		{Key: "aztek", Value: "Pyramid"},
		{Key: "birth", Value: "Certificate"},
		{Key: "critical", Value: "Status"},
		{Key: "detroit", Value: "Factory"},
		{Key: "europe", Value: "Euro"},
		{Key: "frank", Value: "sausages"},
		{Key: "great", Value: "job"},
		{Key: "humongous", Value: "shepherd pie"},
		{Key: "zealous", Value: "pioneers"},
	}
	testTable.L0Merge(testData)
	checkData := []checkData{
		{"abraham", true},
		{"guacamole", false},
		{"bunny", false},
		{"aztek", true},
		{"critical", true},
		{"kelogue", false},
		{"bottom", false},
		{"zealous", true},
	}
	for _, datum := range checkData {
		if val, err := testTable.Get(datum.key); (err == nil) != datum.boolean {
			t.Errorf("WRONG FETCH %s %s %v %v", datum.key, val, err == nil, datum.boolean)
		}
	}

	sstable.CleanAll(testTable)
}
