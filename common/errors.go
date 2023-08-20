package common

import "fmt"

func InvalidBucket(bucket string, args ...interface{}) error {
	return fmt.Errorf(
		"bucket %s does not exist, trace: %v",
		bucket,
		args,
	)
}

func InvalidKey(key string, bucket string, args ...interface{}) error {
	return fmt.Errorf(
		"key %s does not exist in bucket %s, trace: %v",
		key,
		bucket,
		args,
	)
}

func InvalidIndexConfig(bucket string, args ...interface{}) error {
	return fmt.Errorf(
		"invalid index config for bucket %s, trace: %v",
		bucket,
		args,
	)
}

func NoMapping(bucket string, args ...interface{}) error {
	return fmt.Errorf(
		"mapping for bucket %s does not exist, trace: %v",
		bucket,
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
