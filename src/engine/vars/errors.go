package vars

import (
	"errors"
)

var (
	PUT_FAIL_ERROR = errors.New("PUT FAIL ERROR")
	GET_FAIL_ERROR = errors.New("GET FAIL ERROR")

	FILE_CREATE_ERROR = errors.New("FILE CREATE ERROR")
	FILE_WRITE_ERROR  = errors.New("FILE WRITE ERROR")
	FILE_READ_ERROR   = errors.New("FILE READ ERROR")
	FILE_EOF_ERROR    = errors.New("EOF")

	FORMAT_ERROR = errors.New("FORMAT ERROR")

	MEM_TBL_FULL_ERROR  = errors.New("MEM TABLE FULL")
	KEY_NOT_FOUND_ERROR = errors.New("KEY NOT FOUND")
)
