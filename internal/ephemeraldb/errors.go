package ephemeraldb

import "fmt"

type ContainerError struct {
	msg   string
	inner error
}

func (err ContainerError) Error() string {
	return fmt.Sprintf("%s ==> %s", err.msg, err.inner)
}
