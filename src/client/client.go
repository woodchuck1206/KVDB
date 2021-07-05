package client

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

const (
	BASE_URL string = "localhost:7777"
	GET_URL  string = BASE_URL + "/get/"
	PUT_URL  string = BASE_URL + "/put"
	DEL_URL  string = BASE_URL + "/del/"
)

type Command struct {
	CMD Operation
	Key string
	Val string
}

type Operation int

const (
	GET Operation = iota
	PUT
	DEL
)

type Response struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Error string `json:"error"`
}

func intro() {
	fmt.Println("Hello GOKVDB!")
}

func parseArgs(ss []string) {
	cmd := ss[0]
	ss = ss[1:]
	switch strings.ToLower(cmd) {
	case "get":
		get(ss)
	case "put":

	case "del":

	}

}

func get(args []string) {
	if len(args) != 1 {
		fmt.Println("USAGE ERROR!")
		return
	}
	resp, err := http.Get(GET_URL + args[0])
	if err != nil {
		fmt.Println("NETWORK ERROR!")
	}
	fmt.Println(resp)
}

func put(args []string) {
	if len(args) != 2 {
		fmt.Println("USAGE ERROR!")
		return
	}
	resp, err := http.Post(PUT_URL, "application/json")
}

func del(args []string) {
	if len(args) != 1 {
		fmt.Println("USAGE ERROR!")
		return
	}

	resp, err := http.NewRequest(http.MethodDelete, DEL_URL+args[0])
}

func main() {
	params := os.Args
	if len(params) > 0 {
		parseArgs(os.Args)
	} else {
		intro()
	}
}
