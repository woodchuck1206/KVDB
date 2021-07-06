# Go KV-DB (In Progress)
LSM(Log-Structure Merge Tree) Style Key-Value database written in Go. Currently only supports tiering compaction.

# How To

* Library
```go
  import (
    "github.com/woodchuckchoi/KVDB/src/engine"
  )

  memtableSize := 1024 // in bytes
	r := 3 // multiplier for the size of each tier's block
	e := engine.NewEngine(memtableSize, r)

  key := "this is a key"
  value := "and this is a value"

  err := e.Put(key, value)
  if err != nil {
    // ...
  }

  got, err := e.Get(key)
  if err != nil || value != got {
    // ...
  }

  err = e.Delete(key)
  if err != nil {
    // ...
  }

  got, err = e.Get(key)
  if err == nil || got == value {
    // ...
  }

  e.Status() // print status in stdout
  e.CleanAll() // remove all SSTable blocks
```

* Server
```bash
  go build # or install and run the binary

  curl -XPOST localhost:7777/put -d 'key=hello' -d 'value=world'

  curl localhost:7777/get/hello
  # {"key":"hello","value":"world","error":""}

  curl -XDELETE localhost:7777/del/hello

  curl localhost:7777/get/hello
  # {"key":"","value":"","error":"GET FAIL ERROR"}
```

# Caution
DB stores SStable blocks in /tmp/gokvdb at the moment. Will take out all the carved-in-stone variables later.