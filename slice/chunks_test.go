package slice

import (
	"fmt"
	"testing"
)

func ignorePanic(fn func()) {
	defer func() {
		if e := recover(); e != nil {
			fmt.Println(e)
		}
	}()
	fn()
}

func TestChunks(t *testing.T) {
	a := Chunks[int]{{}, {1, 2, 3}, {}, {4, 5, 6}, {7, 8, 9, 10}, {}}
	fmt.Println(a.Slice(10, 10))
	//a := Chunks[int]{{}, {}, {}, {}, {}, {}}
	for i := -1; i < 12; i++ {
		ignorePanic(func() {
			fmt.Println(a.Slice(-1, i))
		})
	}
	for i := -1; i < 12; i++ {
		ignorePanic(func() {
			fmt.Println(a.Slice(i, -1))
		})
	}
	for i := 0; i < 12; i++ {
		for j := 0; j < 12; j++ {
			ignorePanic(func() {
				fmt.Println(a.Slice(i, j))
			})
		}
	}
}
