package main

import (
	"github.com/woodchuckchoi/KVDB/src/controller"
	"github.com/woodchuckchoi/KVDB/src/engine"

	"github.com/woodchuckchoi/KVDB/src/config"
)

func main() {
	c := config.ParseConfig("")
	engine := engine.NewEngine(c)
	server := controller.Init(engine)
	server.Run()
}
