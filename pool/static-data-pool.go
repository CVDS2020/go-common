package pool

import (
	"gitee.com/sy_183/common/assert"
	"gitee.com/sy_183/common/log"
	"gitee.com/sy_183/common/option"
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

func NewStaticDataPool(size uint, poolProvider PoolProvider[*Data], poolOptions ...option.AnyOption) *StaticDataPool {
	return &StaticDataPool{
		pool: poolProvider(func(p Pool[*Data]) *Data {
			return newRefPoolData(p, size)
		}, poolOptions...),
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
		return nil
	}
	d = p.pool.Get()
	if d == nil {
		return nil
	}
	d.AddRef()
	d.Data = d.raw[:len:cap]
	return d
}
