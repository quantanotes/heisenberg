package err

import "fmt"

func BucketExistErr(name string) error {
	return fmt.Errorf("bucket %s already exists", name)
}

func BucketNotExistErr(name string) error {
	return fmt.Errorf("bucket %s does not exist", name)
}
