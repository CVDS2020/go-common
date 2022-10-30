package assert

import "testing"

func TestNotNil(t *testing.T) {
	NotEmpty[bool](true, "test")
}
