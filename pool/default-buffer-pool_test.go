package pool

import (
	"crypto/md5"
	"gitee.com/sy_183/common/unit"
	"math/rand"
	"testing"
)

func TestDefaultBufferPool(t *testing.T) {
	fill := func(buf []byte, ch byte, len int) {
		for i := 0; i < len; i++ {
			buf[i] = ch
		}
	}
	pool := NewDefaultBufferPool(256*unit.KiBiByte, 2048, ProvideSlicePool[*Buffer])
	var chunks []*Data
	var md5s [][16]byte
	ch := byte('0')
	for i := 0; i < 1000000; i++ {
		buf := pool.Get()
		l := rand.Intn(800) + 600
		fill(buf, ch, l)
		md5s = append(md5s, md5.Sum(buf[:l]))
		data := pool.Alloc(uint(l))
		chunks = append(chunks, data)
		if len(chunks) > 1000 {
			for i, chunk := range chunks {
				if md5.Sum(chunk.Data) != md5s[i] {
					panic("")
				}
				chunk.Release()
			}
			chunks = chunks[:0]
			md5s = md5s[:0]
		}
		if ch++; ch > 'z' {
			ch = '0'
		}
	}
}
