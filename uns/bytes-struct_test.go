package uns

import (
	"fmt"
	"testing"
)

func TestBytesStruct(t *testing.T) {
	bs := []byte{65, 66, 67, 68, 69, 70, 71, 72}
	bss := BytesStructOf(bs)
	t.Log(bss)
	t.Log(bss.Bytes())
	t.Log("=========================")
	t.Log(bss.ToString())
	t.Log(bss.ToStringStruct())
}

func TestSliceMerge(t *testing.T) {
	bs := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	bs1 := SliceMerge(bs[1:3], bs[4:5])
	fmt.Println(bs1)
}
