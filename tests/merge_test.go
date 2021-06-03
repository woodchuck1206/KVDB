package tests

import (
	"errors"
	"io"
	"os"
	"testing"
)

const (
	bufferSize = 1
)

type MultiMergeType struct {
	Buffer []byte
	Idx    int
	Init   bool
	Finish bool
	File   *os.File
}

func (this MultiMergeType) IsDone() bool {
	return this.Finish
}

func (this MultiMergeType) LoadMore() {
	_, err := this.File.Read(this.Buffer)
	this.Idx = 0
	if err == io.EOF {
		this.Finish = true
	}
}

func (this MultiMergeType) Get() (byte, error) {
	var ret byte
	if this.Finish {
		return ret, errors.New("Finished")
	}

	if this.Idx == len(this.Buffer) || this.Init {
		this.Init = true
		this.LoadMore()
	}

	ret = this.Buffer[this.Idx]
	if ret == 10 {
		this.Finish = true
		return ret, errors.New("Finished")
	}
	this.Idx++
	return ret, nil
}

func TestMultiMerge(t *testing.T) {
	buffers := []MultiMergeType{}

	fileName1 := "merge1.txt"
	fileName2 := "merge2.txt"

	f1, err := os.Open(fileName1)
	defer f1.Close()
	if err != nil {
		t.Fail()
	}
	buffers = append(buffers, MultiMergeType{
		Buffer: make([]byte, bufferSize),
		Idx:    0,
		Init:   false,
		Finish: false,
		File:   f1,
	})

	f2, err := os.Open(fileName2)
	defer f2.Close()
	if err != nil {
		t.Fail()
	}
	buffers = append(buffers, MultiMergeType{
		Buffer: make([]byte, bufferSize),
		Idx:    0,
		Init:   false,
		Finish: false,
		File:   f2,
	})

	for {

		allDone := true
		for i := 0; i < len(buffers); i++ {
			allDone = allDone && buffers[i].IsDone()
		}
		if allDone {
			break
		}
	}
}
