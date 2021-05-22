package sstable

import (
	"github.com/woodchuckchoi/KVDB/src/engine/util"
	"github.com/woodchuckchoi/KVDB/src/engine/vars"
)

type SSTable struct {
	r      int
	levels []*Level
}

type Level struct {
	blocks []*Block
}

type Block struct {
	fileName string
	index    []vars.SparseIndex
	size     int
}

func (this *SSTable) Get(key string) (string, error) {
	for _, level := range this.levels {
		blockIdx := 0

		for ; blockIdx < len(level.blocks); blockIdx++ {
			if blockIdx == len(level.blocks)-1 || level.blocks[blockIdx].index[0].Key > key {
				break
			}
		}

		val, err := level.blocks[blockIdx].Get(key)
		if err == nil {
			return val, err
		}
	}
	return "", vars.GET_FAIL_ERROR
}

func (this *SSTable) L0Merge(keyValuePairs []vars.KeyValue) error {
	// order := len(this.levels[0].blocks)
	// fileName := util.GenerateFileName(0, order)
	return this.merge(0, keyValuePairs)
	// util.WriteKeyValuePairs()
	// return nil
}

func (this *SSTable) merge(level int, keyValuePairs []vars.KeyValue) error {
	order := 0
	if len(this.levels) > level {
		order = len(this.levels[level].blocks)
	} else {
		this.levels = append(this.levels, &Level{
			blocks: []*Block{},
		})
	}

	fileName := util.GenerateFileName(level, order)
	byteSlice, sparseIndex := util.KeyValueSliceToByteSliceAndSparseIndex(keyValuePairs)
	if util.WriteByteSlice(fileName, byteSlice) != nil {
		return vars.FILE_CREATE_ERROR
	}

	this.levels[level].blocks = append(this.levels[level].blocks, &Block{
		fileName: fileName,
		index:    sparseIndex,
		size:     len(byteSlice),
	})
	return nil
}

func (this *Block) Get(key string) (string, error) {
	from, till := 0, -1

	for _, keyOffsetPair := range this.index {
		if key >= keyOffsetPair.Key {
			from = keyOffsetPair.Offset
		}
		if key < keyOffsetPair.Key {
			till = keyOffsetPair.Offset
			break
		}
	}

	keyValuePairs, err := util.ReadKeyValuePairs(this.fileName, from, till)
	if err != nil {
		return "", err
	}

	return util.BinarySearchKeyValuePairs(keyValuePairs, key)
}

func NewSsTable(r int) *SSTable {
	return &SSTable{
		r:      r,
		levels: []*Level{},
	}
}

func CleanAll(ssTable *SSTable) {
	for _, level := range ssTable.levels {
		for _, block := range level.blocks {
			util.RemoveFile(block.fileName)
		}
	}
}

// type SSTable struct {
// 	r      int
// 	levels []*Level
// }

// type Level struct {
// 	tables []*Table
// 	index  [][]sIndex
// }

// type Table struct { // stored on disk
// 	ptr         *os.File
// 	sparseIndex []sIndex
// }

// type sIndex struct {
// 	key    string
// 	offset int
// }

// func NewSsTable(r int) *SSTable {
// 	return &SSTable{
// 		r:      r,
// 		levels: []*Level{newLevel()},
// 	}
// }

// func newLevel() *Level {
// 	return &Level{
// 		tables: []*Table {},
// 		index: 	[][]sIndex {},
// 	}
// }

// // func newLevel(table *Table) *Level {
// // 	return &Level{
// // 		tables: []*Table{table},
// // 		index:  [][]sIndex{}, // what should it be?
// // 	}
// // }

// func newTable(l, o, size int, data []vars.KeyValue) (*Table, error) {
// 	fileName := generateTableFileName(l, o)
// 	f, err := os.Create(fileName)
// 	if err != nil {
// 		return nil, vars.FILE_CREATE_ERROR
// 	}

// 	toWrite, sparseIndex := parseKeyValue(data, size)
// 	_, err = f.Write(toWrite)
// 	if err != nil {
// 		return nil, vars.FILE_WRITE_ERROR
// 	}

// 	return &Table{
// 		ptr:         f,
// 		sparseIndex: sparseIndex,
// 	}, nil
// }

// func (this *SSTable) L0Merge(keyValuePairs []vars.KeyValue) error {
// 	level, order := 0, len(this.levels[0].tables, keyValuePairs)
// 	table, err := newTable(level, order, keyValuePairs)
// 	if err != nil {
// 		return err
// 	}

// 	this.levels[0].tables = append(this.levels[0].tables, table)
// 	this.levels[0].index 	= append(this.levels[0].index, )

// 	byteSliceToWrite, sparseIndex := parseKeyValue(keyValuePairs)

// 	if this.levels
// }

// func getDiskSizeOfKeyValuePairs(data []vars.KeyValue) int {
// 	length := 0
// 	for _, keyValuePair := range data {
// 		length += len(keyValuePair.Key) + len(keyValuePair.Value) + 2 // for separator and delimiter
// 	}
// 	return length
// }

// func parseKeyValue(data []vars.KeyValue) ([]byte, []sIndex) {
// 	size := getDiskSizeOfKeyValuePairs(data)
// 	bdta := make([]byte, size)
// 	sidx := []sIndex{}
// 	offset, lastOffset := 0, 0

// 	for _, pair := range data {
// 		byteString := util.KeyValueToByteSlice(pair)
// 		byteCopyHelper(byteString, &bdta, offset)

// 		if offset-lastOffset >= INDEX_TERM || offset == 0 {
// 			sidx = append(sidx, sIndex{
// 				key:    pair.Key,
// 				offset: offset,
// 			})
// 			lastOffset = offset
// 		}
// 		offset += len(byteString)
// 	}

// 	return bdta, sidx
// }

// // ---

// func (sstable *SSTable) SaveFlushOnDisk(data []vars.KeyValue) error {
// 	var order int
// 	if len(sstable.levels) == 0 {
// 		order = 0
// 	} else {
// 		order = len(sstable.levels[0].tables)
// 	}

// 	table, err := newTable(0, order, sstable.flushSize, data)
// 	if err != nil {
// 		return err
// 	}

// 	if len(sstable.levels) == 0 {
// 		new
// 		sstable.levels = append(sstable.levels)
// 	}
// }

// func byteCopyHelper(src []byte, dest *[]byte, offset int) {
// 	for idx, byteCharacter := range src {
// 		if len(*dest) < len(src)-idx+offset {
// 			*dest = append(*dest, byteCharacter)
// 		} else {
// 			(*dest)[offset+idx] = byteCharacter
// 		}
// 	}
// }

// func generateTableFileName(level, order int) string {
// 	fileName := fmt.Sprintf("db:%d:%d.db", level, order)
// 	return path.Join(BASE_DIR, fileName)
// }

// func (sstable *SSTable) Get(key string) (string, error) {
// 	for _, level := range sstable.levels {
// 		partition := 0
// 		for partition < len(level.index)-1 {
// 			if key >= level.index[partition][0].key && key < level.index[partition+1][0].key {
// 				break
// 			}
// 			partition++
// 		}
// 		val, err := level.tables[partition].Get(key)
// 		if err == nil {
// 			return val, nil
// 		}
// 	}
// 	return "", vars.GET_FAIL_ERROR
// }

// func (table Table) Get(key string) (string, error) {
// 	partition := 0
// 	from, till := -1, -1
// 	for partition < len(table.sparseIndex)-1 {
// 		if key >= table.sparseIndex[partition].key && key < table.sparseIndex[partition+1].key {
// 			from = table.sparseIndex[partition].offset
// 			break
// 		}
// 		partition++
// 	}

// 	if partition != len(table.sparseIndex)-1 {
// 		till = table.sparseIndex[partition+1].offset
// 	}

// 	return GetFromFile(table.ptr, from, till, key)
// }

// func GetFromFile(f *os.File, from, till int, key string) (string, error) {
// 	byteChunk, err := loadChunk(f, from, till)
// 	if err != nil {
// 		return "", err
// 	}
// 	kvPairs, err := parseBytes(byteChunk)
// 	return binSearch(kvPairs, key)
// }

// func loadChunk(f *os.File, from, till int) ([]byte, error) {
// 	_, err := f.Seek(int64(from), 0)
// 	if err != nil {
// 		return nil, vars.FILE_READ_ERROR
// 	}

// 	var size int

// 	if till != -1 {
// 		size = till - from
// 	} else {
// 		info, _ := f.Stat()
// 		size = int(info.Size()) - from
// 	}

// 	ret := make([]byte, size)

// 	_, err = f.Read(ret)
// 	if err != nil {
// 		return nil, vars.FILE_READ_ERROR
// 	}
// 	return ret, nil
// }

// func parseBytes(bSlice []byte) ([]vars.KeyValue, error) {
// 	byteRecords := bytes.Split(bSlice, []byte{vars.DELIMITER})
// 	ret := make([]vars.KeyValue, len(byteRecords))
// 	for idx, byteRecord := range byteRecords {
// 		keyValue := bytes.Split(byteRecord, []byte{vars.SEPARATOR})
// 		if len(keyValue) != 2 {
// 			return nil, vars.FORMAT_ERROR
// 		}
// 		ret[idx] = vars.KeyValue{
// 			Key:   string(keyValue[0]),
// 			Value: string(keyValue[1]),
// 		}
// 	}
// 	return ret, nil
// }

// func binSearch(s []vars.KeyValue, key string) (string, error) {
// 	left, right := 0, len(s)
// 	for left < right {
// 		mid := (left + right) / 2
// 		if s[mid].Key == key {
// 			return s[mid].Value, nil
// 		}

// 		if s[mid].Key < key {
// 			left = mid + 1
// 		} else {
// 			right = mid
// 		}
// 	}
// 	return "", vars.KEY_NOT_FOUND_ERROR
// }

// // func loadChunk(f *os.File, from, till int) ([]byte, error) {
// // 	_, err := f.Seek(int64(from), 0)
// // 	if err != nil {
// // 		return nil, vars.FILE_READ_ERROR
// // 	}
// // 	f.Close()
// // 	ret := make([]byte, )
// // 	f.Read()
// // }
