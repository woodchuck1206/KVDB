package vars

// --- memtable, sstable

type KeyValue struct {
	Key   string
	Value string
}

type SparseIndex struct {
	Key    string
	Offset int
}
