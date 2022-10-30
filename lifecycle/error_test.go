package lifecycle

import (
	"fmt"
	"testing"
)

func TestIsStateError(t *testing.T) {
	var err error
	fmt.Println(IsStateError(err))
}
