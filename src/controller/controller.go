package controller

import (
	"github.com/labstack/echo/v4"
	"github.com/woodchuckchoi/KVDB/src/engine"
	"github.com/woodchuckchoi/KVDB/src/engine/vars"
)

type KVServer struct {
	Engine *engine.Engine
	Echo   *echo.Echo
}

func Init(engine *engine.Engine) KVServer {
	e := echo.New()
	mapEcho(e)

	return KVServer{
		Engine: engine,
		Echo:   e,
	}
}

func (this KVServer) Run() {
	this.Echo.Logger.Fatal(this.Echo.Start(vars.DEFAULT_PORT))
}

func mapEcho(e *echo.Echo) {
	e.GET("/get/:key", get)
	e.PUT("/put", put)
	e.DELETE("/del", del)
}

func put(c echo.Context) error {
	key := c.FormValue("key")
	value := c.FormValue("value")

}

func get(c echo.Context) error {
	key := c.Param("key")
}

func del(c echo.Context) error {

}
