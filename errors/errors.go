package Errors

import "fmt"

type errorCode uint

const (
	Nonce errorCode = iota
	ConnectionErrorCode
)

func ConnectionError(addr string, args ...interface{}) error {
	return fmt.Errorf("error %d: connection to address %s failed, trace: %v", ConnectionErrorCode, addr, args)
}
