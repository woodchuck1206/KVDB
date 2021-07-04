package client

import (
	"fmt"
	"os"
	"strings"
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

func intro() {
	fmt.Println("Hello GOKVDB!")
}

func parseArgs(ss []string) {
	switch strings.ToLower(ss[0]) {
	case "get":

	case "put":

	case "del":

	}

}

func main() {
	params := os.Args
	if len(params) > 0 {
		parseArgs(os.Args)
	} else {
		intro()
	}
}
