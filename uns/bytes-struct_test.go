package uns

import "testing"

func TestBytesStruct(t *testing.T) {
	bs := []byte{65, 66, 67, 68, 69, 70, 71, 72}
	bss := BytesStructOf(bs)
	t.Log(bss)
	t.Log(bss.Bytes())
	t.Log("=========================")
	t.Log(bss.ToString())
	t.Log(bss.ToStringStruct())
}
