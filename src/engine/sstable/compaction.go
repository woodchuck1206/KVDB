package sstable

import (
	"io"
	"os"

	"github.com/woodchuckchoi/KVDB/src/engine/vars"
)

type Compactor struct {
	Receiver <-chan MergeSignal
}

type MergeSignal struct {
	Level    int
	LevelRef *SSTable
}

type MergeUnit struct {
	Buffer []byte
	Index  int
	Init   bool
	Finish bool
	File   *os.File
}

func (this *MergeUnit) Load() { // load size should be managed // each key value pairs sizes differ
	_, err := this.File.Read(this.Buffer)
	this.Index = 0
	if err == io.EOF {
		this.Finish = true
	}
}

func (this *MergeUnit) Get() (byte, error) { // should return keyvalue pair
	var ret byte
	if this.Finish {
		return ret, vars.FILE_EOF_ERROR
	}

	if this.Index >= len(this.Buffer) || !this.Init {
		this.Init = true
		this.Index = 0
		this.Load()
	}

	ret = this.Buffer[this.Index]
	// if ret...
}

func (this *MergeUnit) Pop() (byte, error) {
	ret, err := this.Get()
	this.Index++
	return ret, err
}

func NewCompactor(chanBuffer int) (*Compactor, chan<- MergeSignal) {
	if chanBuffer <= 0 {
		chanBuffer = 42 // arbitrary default number for buffer size
	}
	channel := make(chan MergeSignal, chanBuffer)
	return &Compactor{Receiver: channel}, channel
}

func MultiMerge(level *Level) Block {

}

// func (this *Compactor) Run() {
// 	for {
// 		select {
// 		case sig := <-this.Receiver:
// 			// merge using sig and update the reference
// 			levelToMerge := sig.Level
// 			sstable := sig.LevelRef

// 			sstable.levels[0]
// 		}
// 	}
// }

// func mergeLevel(level *Level) *Block {

// }

// func multiMerge()
