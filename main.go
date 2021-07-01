package main

import (
	"github.com/woodchuckchoi/KVDB/src/controller"
	"github.com/woodchuckchoi/KVDB/src/engine"
)

func main() {
	engine := engine.NewEngine(2048, 3)
	server := controller.Init(engine)
	server.Run()
}
