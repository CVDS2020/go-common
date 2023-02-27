package stream2

import (
	"fmt"
	"gitee.com/sy_183/common/uns"
	"testing"
)

func TestStream(t *testing.T) {
	stream := NewQueueStream(0)
	for i := 0; i < 32; i++ {
		stream.Write(uns.StringToBytes(fmt.Sprintf("hello%d", i)), nil)
	}
	if s, err := stream.ReadString(214); err == nil {
		fmt.Println(s)
	}
}
