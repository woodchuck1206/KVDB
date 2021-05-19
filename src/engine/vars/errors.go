package vars

import (
	"errors"
)

var (
	PUT_FAIL_ERROR    = errors.New("PUT FAIL ERROR")
	GET_FAIL_ERROR    = errors.New("GET FAIL ERROR")
	FILE_CREATE_ERROR = errors.New("FILE CREATE ERROR")
)
