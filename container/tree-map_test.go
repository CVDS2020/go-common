package container

import (
	"fmt"
	"github.com/gofrs/uuid"
	"testing"
)

func TestTreeMap(t *testing.T) {
	m := NewOrderedKeyTreeMap[int, string]()
	for i := 0; i < 10000000; i++ {
		m.Put(i, uuid.Must(uuid.NewV4()).String())
	}
	fmt.Println(m.Size())

	var next *TreeMapEntry[int, string]
	for first := m.GetFirstEntry(); first != nil; first = next {
		next = first.Next()
		m.RemoveEntry(first)
		//fmt.Printf("remove %s\n", first)
	}

	fmt.Println(m.Size())
}

func TestMap(t *testing.T) {
	m := make(map[int]string)
	for i := 0; i < 10000000; i++ {
		m[i] = uuid.Must(uuid.NewV4()).String()
	}
	fmt.Println(len(m))

	for i := 0; i < 10000000; i++ {
		delete(m, i)
	}

	fmt.Println(len(m))
}
