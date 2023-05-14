package utils

import "fmt"

func InvalidCollection(collection string, args ...interface{}) error {
	return fmt.Errorf(
		"collection %s does not exist, trace: %v",
		collection,
		args,
	)
}

func InvalidKey(key string, collection string, args ...interface{}) error {
	return fmt.Errorf(
		"key %s does not exist in collection %s, trace: %v",
		key,
		collection,
		args,
	)
}

func InvalidIndexConfig(collection string, args ...interface{}) error {
	return fmt.Errorf(
		"invalid index config for collection %s, trace: %v",
		collection,
		args,
	)
}

func NoMapping(collection string, args ...interface{}) error {
	return fmt.Errorf(
		"mapping for collection %s does not exist, trace: %v",
		collection,
		args,
	)
}

func DimensionMismatch(expected int, recieved int, args ...interface{}) error {
	return fmt.Errorf(
		"dimension mismatch, expected: %d, recieved: %d, trace: %v",
		expected,
		recieved,
		args,
	)
}

func JsonUnmarshalError(err error, args ...interface{}) error {
	return fmt.Errorf(
		"json unmarshal error: %v, trace: %v",
		err,
		args,
	)
}
