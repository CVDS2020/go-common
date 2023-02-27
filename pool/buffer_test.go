package pool

import (
	"gitee.com/sy_183/common/unit"
	"testing"
)

func TestBuffer(t *testing.T) {
	buffer := NewBuffer(unit.MeBiByte, 2048)
	for {
		buffer.Get()
	}
}
