package internal

import (
	"fmt"
)

type errorCode uint

const (
	NullError errorCode = iota
	ConnectionErrorCode
	IncorrectServiceErrorCode
)

func ConnectionError(addr string, args ...interface{}) error {
	return fmt.Errorf(
		"%d: connection to address %s failed, trace: %v",
		ConnectionErrorCode,
		addr,
		args,
	)
}

func IncorrectServiceError(expected Service, recieved Service, args ...interface{}) error {
	return fmt.Errorf(
		"%d: connected to incorrect service, expected %s, recieved %s, trace: %v",
		IncorrectServiceErrorCode,
		expected.String(),
		recieved.String(),
		args,
	)
}
