package sstable

type Compactor struct {
	Receiver <-chan MergeSignal
}

type MergeSignal struct {
	Level    int
	LevelRef *SSTable
}

func NewCompactor(chanBuffer int) (*Compactor, chan<- MergeSignal) {
	if chanBuffer <= 0 {
		chanBuffer = 42
	}
	channel := make(chan MergeSignal, chanBuffer)
	return &Compactor{Receiver: channel}, channel
}

func (this *Compactor) Run() {
	for {
		select {
		case sig := <-this.Receiver:
			// merge using sig and update the reference
		}
	}
}
