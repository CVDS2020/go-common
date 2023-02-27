package pool

import (
	"gitee.com/sy_183/common/assert"
	"gitee.com/sy_183/common/log"
)

var DebugLogger = assert.Must(log.Config{
	Level: log.NewAtomicLevelAt(log.DebugLevel),
	Encoder: log.NewConsoleEncoder(log.ConsoleEncoderConfig{
		EncodeLevel: log.CapitalColorLevelEncoder,
	}),
}.Build())

type StaticDataPool struct {
	pool Pool[*Data]
	size uint
}

func NewStaticDataPool(size uint, poolProvider PoolProvider[*Data]) *StaticDataPool {
	return &StaticDataPool{
		pool: poolProvider(func(p Pool[*Data]) *Data {
			return newRefPoolData(p, size)
		}),
		size: size,
	}
}

func (p *StaticDataPool) Size() uint {
	return p.size
}

func (p *StaticDataPool) Alloc(len uint) (d *Data) {
	return p.AllocCap(len, len)
}

func (p *StaticDataPool) AllocCap(len, cap uint) (d *Data) {
	if cap > p.size {
		return NewData(make([]byte, len, cap))
	}
	d = p.pool.Get().Use()
	d.Data = d.raw[:len:cap]
	return d
}
