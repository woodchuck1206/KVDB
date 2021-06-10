package tests

import (
	"testing"

	"github.com/woodchuckchoi/KVDB/src/engine/sstable"
	"github.com/woodchuckchoi/KVDB/src/engine/util"
	"github.com/woodchuckchoi/KVDB/src/engine/vars"
)

func TestCompaction(t *testing.T) {
	block1Name := "block1.data"
	block2Name := "block2.data"
	block3Name := "block3.data"

	block1KeyValue := []vars.KeyValue {
		vars.KeyValue{

		},
		vars.KeyValue{

		},
		vars.KeyValue{

		},
		vars.KeyValue{

		},
		vars.KeyValue{

		},
	}

	block2KeyValue := []vars.KeyValue {
		vars.KeyValue{

		},
		vars.KeyValue{

		},
		vars.KeyValue{

		},
		vars.KeyValue{

		},
		vars.KeyValue{

		},
	}

	block3KeyValue := []vars.KeyValue {
		vars.KeyValue{

		},
		vars.KeyValue{

		},
		vars.KeyValue{

		},
		vars.KeyValue{

		},
		vars.KeyValue{

		},
	}

	util.WriteKeyValuePairs()

	block1 := sstable.Block {

	}

	block2 := sstable.Block {

	}

	block3 := sstable.Block {

	}
	
	level := sstable.Level {
		Blocks: ,
	}

	sstable.MultiMerge()
}