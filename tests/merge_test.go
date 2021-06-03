package tests

import (
	"errors"
	"io"
	"os"
	"sort"
	"testing"
)

const (
	bufferSize = 2
	merge1Data = "abdejjlmoooruv"
	merge2Data = "ceeeeeeefguuuz"
)

type MultiMergeType struct {
	Buffer []byte
	Idx    int
	Init   bool
	Finish bool
	File   *os.File
}

func (this *MultiMergeType) IsDone() bool {
	return this.Finish
}

func (this *MultiMergeType) LoadMore() {
	_, err := this.File.Read(this.Buffer)
	this.Idx = 0
	if err == io.EOF {
		this.Finish = true
	}
}

func (this *MultiMergeType) Peek() (byte, error) {
	var ret byte
	if this.Finish {
		return ret, errors.New("Finished")
	}

	if this.Idx >= len(this.Buffer) || !this.Init {
		this.Init = true
		this.Idx = 0
		this.LoadMore()
	}

	ret = this.Buffer[this.Idx]
	if ret == 10 { // comes at the end of a line && only for test
		this.Finish = true
		return ret, errors.New("Finished")
	}
	return ret, nil
}

func (this *MultiMergeType) Get() (byte, error) {
	ret, err := this.Peek()
	this.Idx++
	return ret, err
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

	merged := []byte{}
	for {
		var smallestBuffer *MultiMergeType
		for i := 0; i < len(buffers); i++ {
			ret, err := buffers[i].Peek()
			if err != nil {
				continue
			}
			if smallestBuffer != nil {
				origRet, _ := smallestBuffer.Peek()
				if origRet <= ret {
					continue
				}
			}

			smallestBuffer = &(buffers[i])
		}

		if smallestBuffer == nil {
			break
		}
		toAdd, _ := smallestBuffer.Get()
		merged = append(merged, toAdd)
	}
	answer := []byte{}
	answer = append(answer, []byte(merge1Data)...)
	answer = append(answer, []byte(merge2Data)...)
	sort.Slice(answer, func(i, j int) bool {
		if answer[i] < answer[j] {
			return true
		}
		return false
	})

	if string(merged) != string(answer) {
		t.Fail()
	}
}
