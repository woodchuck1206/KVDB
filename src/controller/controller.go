package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/woodchuckchoi/KVDB/src/engine"
	"github.com/woodchuckchoi/KVDB/src/engine/vars"
)

type KVServer struct {
	Engine *engine.Engine
	Echo   *echo.Echo
}

type Response struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Error string `json:"error"`
}

func generateResponse(key, value string, err error) *Response {
	if err != nil {
		return &Response{
			Error: err.Error(),
		}
	}
	return &Response{
		Key:   key,
		Value: value,
	}
}

func Init(engine *engine.Engine) *KVServer {
	e := echo.New()
	kvServer := &KVServer{
		Engine: engine,
		Echo:   e,
	}
	kvServer.mapEcho()
	return kvServer
}

func (this *KVServer) Run() {
	this.Echo.Logger.Fatal(this.Echo.Start(vars.DEFAULT_PORT))
}

func (this *KVServer) mapEcho() {
	this.Echo.GET("/get/:key", this.get)
	this.Echo.POST("/put", this.put)
	this.Echo.DELETE("/del/:key", this.del)
}

func (this *KVServer) put(c echo.Context) error {
	key := c.FormValue("key")
	value := c.FormValue("value")

	err := this.Engine.Put(key, value)
	resp := generateResponse(key, value, err)
	return c.JSON(http.StatusOK, resp)
}

func (this *KVServer) get(c echo.Context) error {
	key := c.Param("key")

	value, err := this.Engine.Get(key)
	resp := generateResponse(key, value, err)
	return c.JSON(http.StatusOK, resp)
}

func (this *KVServer) del(c echo.Context) error {
	key := c.Param("key")

	err := this.Engine.Delete(key)
	resp := generateResponse(key, "", err)
	return c.JSON(http.StatusOK, resp)
}
