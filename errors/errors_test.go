package errors

import (
	"fmt"
	"testing"
)

func TestErrors(t *testing.T) {
	arr := Errors{}
	e := New("hello")
	fmt.Println(arr.Error(), e)
}
