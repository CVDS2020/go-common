package stream

import (
	"errors"
	"gitee.com/sy_183/common/pool"
	"gitee.com/sy_183/common/uns"
	"io"
)

const StreamDefaultCapacity = 16

var IndexOutOfRange = errors.New("index out of range")

type RDStream interface {
	Size() uint

	Read(size uint) ([]byte, error)

	ReadAll() ([]byte, error)

	ReadByte() (byte, error)

	ReadString(size uint) (string, error)

	ReadChunks(size uint, chunks [][]byte) ([][]byte, error)

	Peek(size uint) ([]byte, error)

	PeekAll() ([]byte, error)

	PeekByte() (byte, error)

	PeekLastByte() (byte, error)

	PeekIndexByte(i uint) (byte, error)

	PeekString(size uint) (string, error)

	PeekChunks(size uint, chunks [][]byte) ([][]byte, error)

	Skip(size uint) error
}

type Stream interface {
	RDStream

	Write(data []byte, ref pool.Reference)

	WriteString(str string, ref pool.Reference)
}

func Must(err error) {
	if err != nil {
		panic(err)
	}
}

func MustBytes(data []byte, err error) []byte {
	if err != nil {
		panic(err)
	}
	return data
}

func MustByte(b byte, err error) byte {
	if err != nil {
		panic(err)
	}
	return b
}

func MustString(s string, err error) string {
	if err != nil {
		panic(err)
	}
	return s
}

func MustChunks(chunks [][]byte, err error) [][]byte {
	if err != nil {
		panic(err)
	}
	return chunks
}

type Node struct {
	next *Node
	prev *Node
	data []byte
	ref  pool.Reference
}

type QueueStream struct {
	head *Node
	tail *Node
	cap  uint
	len  uint
	size uint
}

func NewQueueStream(cap uint) *QueueStream {
	if cap == 0 {
		cap = StreamDefaultCapacity
	}
	s := new(QueueStream)
	s.grow(cap)
	return s
}

func (s *QueueStream) grow(grow uint) {
	nodes := make([]Node, grow)
	for i := range nodes {
		if i != 0 {
			nodes[i].prev = &nodes[i-1]
		}
		if i != int(grow)-1 {
			nodes[i].next = &nodes[i+1]
		}
	}
	if s.cap == 0 {
		nodes[0].prev = &nodes[grow-1]
		nodes[grow-1].next = &nodes[0]
		s.head = &nodes[0]
		s.tail = &nodes[0]
	} else {
		nodes[0].prev = s.head.prev
		nodes[grow-1].next = s.head
		s.head.prev.next = &nodes[0]
		s.head.prev = &nodes[grow-1]
		if s.len == s.cap {
			s.tail = &nodes[0]
		}
	}
	s.cap += grow
}

func (s *QueueStream) readNext() {
	if s.head.ref != nil {
		s.head.ref.Release()
	}
	s.head = s.head.next
	s.len--
}

func (s *QueueStream) Len() uint {
	return s.len
}

func (s *QueueStream) Cap() uint {
	return s.cap
}

func (s *QueueStream) Size() uint {
	return s.size
}

func (s *QueueStream) Write(data []byte, ref pool.Reference) {
	if len(data) == 0 {
		ref.Release()
		return
	}
	if s.len == s.cap {
		if s.cap < 1024 {
			s.grow(s.cap)
		} else {
			s.grow(s.cap >> 2)
		}
	}
	s.tail.data = data
	s.tail.ref = ref
	s.tail = s.tail.next
	s.len++
	s.size += uint(len(data))
}

func (s *QueueStream) WriteString(str string, ref pool.Reference) {
	s.Write(uns.StringToBytes(str), ref)
}

func (s *QueueStream) Read(size uint) ([]byte, error) {
	if size == 0 {
		return nil, nil
	} else if s.len == 0 || size > s.size {
		return nil, io.EOF
	}
	s.size -= size
	cur := s.head.data
	if size <= uint(len(cur)) {
		data, update := cur[:size], cur[size:]
		if len(update) == 0 {
			s.readNext()
		} else {
			s.head.data = update
		}
		return data, nil
	}
	data := make([]byte, size)
	var need uint
	for i := uint(0); size > 0; i += need {
		if need = uint(len(cur)); need > size {
			need = size
			s.head.data = cur[need:]
			copy(data[i:], cur[:need])
			break
		}
		copy(data[i:], cur)
		s.readNext()
		size -= need
		cur = s.head.data
	}
	return data, nil
}

func (s *QueueStream) ReadAll() ([]byte, error) {
	defer s.Clear()
	return s.PeekAll()
}

func (s *QueueStream) ReadByte() (byte, error) {
	if s.size < 1 {
		return 0, io.EOF
	}
	cur := s.head.data
	b := cur[0]
	if len(cur) == 1 {
		s.readNext()
	} else {
		s.head.data = cur[1:]
	}
	s.size--
	return b, nil
}

func (s *QueueStream) ReadString(size uint) (string, error) {
	bytes, err := s.Read(size)
	return uns.BytesToString(bytes), err
}

func (s *QueueStream) ReadChunk() ([]byte, error) {
	if s.len == 0 {
		return nil, io.EOF
	}
	data := s.head.data
	s.size -= uint(len(data))
	s.readNext()
	return data, nil
}

func (s *QueueStream) ReadChunks(size uint, chunks [][]byte) ([][]byte, error) {
	if size == 0 {
		return nil, nil
	} else if s.len == 0 || size > s.size {
		return nil, io.EOF
	}
	s.size -= size
	cur := s.head.data
	var need uint
	for i := uint(0); size > 0; i += need {
		if need = uint(len(cur)); need > size {
			need = size
			s.head.data = cur[need:]
			chunks = append(chunks, cur[:need])
			break
		}
		chunks = append(chunks, cur)
		s.readNext()
		size -= need
		cur = s.head.data
	}
	return chunks, nil
}

func (s *QueueStream) Peek(size uint) ([]byte, error) {
	if size == 0 {
		return nil, nil
	} else if s.len == 0 || size > s.size {
		return nil, io.EOF
	}
	curNode := s.head
	if size <= uint(len(curNode.data)) {
		return curNode.data[:size], nil
	}
	data := make([]byte, size)
	var need uint
	for i := uint(0); size > 0; i += need {
		if need = uint(len(curNode.data)); need > size {
			need = size
			copy(data[i:], curNode.data[:need])
			break
		}
		copy(data[i:], curNode.data)
		curNode = curNode.next
		size -= need
	}
	return data, nil
}

func (s *QueueStream) PeekAll() ([]byte, error) {
	if s.size == 0 {
		return nil, nil
	}
	if s.len == 1 {
		return s.head.data, nil
	} else {
		data := make([]byte, s.size)
		curNode := s.head
		for i := 0; curNode == s.tail; i += len(curNode.data) {
			copy(data[i:], curNode.data)
		}
		return data, nil
	}
}

func (s *QueueStream) PeekByte() (byte, error) {
	if s.size < 1 {
		return 0, io.EOF
	}
	return s.head.data[0], nil
}

func (s *QueueStream) PeekIndexByte(i uint) (byte, error) {
	if i >= s.size {
		return 0, IndexOutOfRange
	}
	curNode := s.head
	var need uint
	for i > 0 {
		if need = uint(len(curNode.data)); need > i {
			break
		}
		curNode = curNode.next
		i -= need
	}
	return curNode.data[i], nil
}

func (s *QueueStream) PeekLastByte() (byte, error) {
	if s.size < 1 {
		return 0, io.EOF
	}
	last := s.tail.prev.data
	return last[len(last)-1], nil
}

func (s *QueueStream) PeekString(size uint) (string, error) {
	bytes, err := s.Peek(size)
	return uns.BytesToString(bytes), err
}

func (s *QueueStream) PeekChunk() ([]byte, error) {
	if s.len == 0 {
		return nil, io.EOF
	}
	return s.head.data, nil
}

func (s *QueueStream) PeekChunks(size uint, chunks [][]byte) ([][]byte, error) {
	if size == 0 {
		return nil, nil
	} else if s.len == 0 || size > s.size {
		return nil, io.EOF
	}
	curNode := s.head
	var need uint
	for i := uint(0); size > 0; i += need {
		if need = uint(len(curNode.data)); need > size {
			need = size
			chunks = append(chunks, curNode.data[:need])
			break
		}
		chunks = append(chunks, curNode.data)
		curNode = curNode.next
		size -= need
	}
	return chunks, nil
}

func (s *QueueStream) Skip(size uint) error {
	if size == 0 {
		return nil
	} else if s.len == 0 || size > s.size {
		return io.EOF
	}
	s.size -= size
	cur := s.head.data
	var need uint
	for i := uint(0); size > 0; i += need {
		if need = uint(len(cur)); need > size {
			need = size
			s.head.data = cur[need:]
			break
		}
		s.readNext()
		size -= need
		cur = s.head.data
	}
	return nil
}

func (s *QueueStream) SkipChunk() error {
	if s.len == 0 {
		return io.EOF
	}
	s.size -= uint(len(s.head.data))
	s.readNext()
	return nil
}

func (s *QueueStream) Clear() {
	s.tail = s.head
	s.len = 0
	s.size = 0
}

type BytesStream []byte

func (s *BytesStream) Size() uint {
	return uint(len(*s))
}

func (s *BytesStream) Read(size uint) ([]byte, error) {
	if size > s.Size() {
		return nil, io.EOF
	}
	res := (*s)[:size]
	*s = (*s)[size:]
	return res, nil
}

func (s *BytesStream) ReadAll() ([]byte, error) {
	res := *s
	*s = (*s)[len(*s):]
	return res, nil
}

func (s *BytesStream) ReadByte() (byte, error) {
	if len(*s) == 0 {
		return 0, io.EOF
	}
	b := (*s)[0]
	*s = (*s)[1:]
	return b, nil
}

func (s *BytesStream) ReadString(size uint) (string, error) {
	bytes, err := s.Read(size)
	return uns.BytesToString(bytes), err
}

func (s *BytesStream) ReadChunks(size uint, chunks [][]byte) ([][]byte, error) {
	bytes, err := s.Read(size)
	if len(bytes) == 0 {
		return chunks, err
	}
	return append(chunks, bytes), err
}

func (s *BytesStream) Peek(size uint) ([]byte, error) {
	if size > s.Size() {
		return nil, io.EOF
	}
	return (*s)[:size], nil
}

func (s *BytesStream) PeekAll() ([]byte, error) {
	return *s, nil
}

func (s *BytesStream) PeekByte() (byte, error) {
	if len(*s) == 0 {
		return 0, io.EOF
	}
	return (*s)[0], nil
}

func (s *BytesStream) PeekLastByte() (byte, error) {
	if len(*s) == 0 {
		return 0, io.EOF
	}
	return (*s)[len(*s)-1], nil
}

func (s *BytesStream) PeekIndexByte(i uint) (byte, error) {
	if i >= s.Size() {
		return 0, IndexOutOfRange
	}
	return (*s)[i], nil
}

func (s *BytesStream) PeekString(size uint) (string, error) {
	bytes, err := s.Peek(size)
	return uns.BytesToString(bytes), err
}

func (s *BytesStream) PeekChunks(size uint, chunks [][]byte) ([][]byte, error) {
	bytes, err := s.Peek(size)
	if len(bytes) == 0 {
		return chunks, err
	}
	return append(chunks, bytes), err
}

func (s *BytesStream) Skip(size uint) error {
	if size > s.Size() {
		return io.EOF
	}
	*s = (*s)[size:]
	return nil
}
