package errors

import (
	"fmt"
	"testing"
)

func TestErrors(t *testing.T) {
	arr := Errors{}
	fmt.Println(arr.Error())
}
