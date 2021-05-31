package compaction

import (
	
)

type Compactor struct {
	pool		<-chan 
}

type MergeRequest interface {
	GetLevels() []*Level
}