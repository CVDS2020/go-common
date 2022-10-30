package uns

import "testing"

func TestBufferWriter_Write(t *testing.T) {
	w := NewBufferWriter(make([]byte, 0, 3))
	w.Write(StringToBytes("hello"))
}
