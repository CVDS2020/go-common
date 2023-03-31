package container

import (
	"fmt"
	"testing"
)

func TestQueue(t *testing.T) {
	queue := NewQueue[int](10)
	queue.PushHead(9)
	queue.PushHead(8)
	queue.Push(0)
	queue.Push(1)
	queue.Push(2)
	queue.Pop()
	fmt.Println(queue.Get(1))
	fmt.Println(queue.Tail())
	fmt.Println(queue.Head())
	fmt.Println(queue.Slice(0, 1))
}

func TestPanic(t *testing.T) {
	a := make([]int, 1, 10)
	//goPanicSliceB(1, 9)
	i, j := -1, 9
	_, _ = i, j
	_ = a[i:]
}
