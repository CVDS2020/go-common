package container

import (
	"fmt"
	"github.com/gofrs/uuid"
	"testing"
)

func TestLinkedMap(t *testing.T) {
	m := NewLinkedMap[int, string](10000)
	for i := 0; i < 1000000; i++ {
		m.Put(i, uuid.Must(uuid.NewV4()).String())
	}
	fmt.Println(m.Size())

	var next *LinkedMapEntry[int, string]
	for first := m.FirstEntry(); first != nil; first = next {
		next = first.Next()
		m.RemoveEntry(first)
		fmt.Printf("remove %s\n", first.Value())
	}

	fmt.Println(m.Size())
}
