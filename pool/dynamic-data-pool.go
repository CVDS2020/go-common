package pool

import (
	"gitee.com/sy_183/common/option"
	"sort"
)

type DynamicDataPool struct {
	pools []*StaticDataPool
	sizes []uint
}

func NewDynamicDataPool(pools ...*StaticDataPool) *DynamicDataPool {
	sort.Slice(pools, func(i, j int) bool {
		return pools[i].Size() < pools[j].Size()
	})
	p := new(DynamicDataPool)
	var last uint
	for _, pool := range pools {
		if size := pool.Size(); size > last {
			p.pools = append(p.pools, pool)
			p.sizes = append(p.sizes, size)
			last = size
		}
	}
	return p
}

func NewDynamicDataPoolWithThresholds(thresholds []int, poolProvider PoolProvider[*Data], poolOptions ...option.AnyOption) *DynamicDataPool {
	sort.Ints(thresholds)
	p := new(DynamicDataPool)
	var last int
	for _, threshold := range thresholds {
		if threshold > last {
			p.pools = append(p.pools, NewStaticDataPool(uint(threshold), poolProvider, poolOptions...))
			p.sizes = append(p.sizes, uint(threshold))
			last = threshold
		}
	}
	return p
}

func NewDynamicDataPoolWithExp(min, max uint, poolProvider PoolProvider[*Data], poolOptions ...option.AnyOption) *DynamicDataPool {
	var thresholds []int
	for i := 1; i < 63; i++ {
		if uint(1<<i) >= min && uint(1<<(i-1)) < max {
			thresholds = append(thresholds, 1<<i)
		} else if uint(1<<(i-1)) >= max {
			break
		}
	}
	return NewDynamicDataPoolWithThresholds(thresholds, poolProvider, poolOptions...)
}

func (p *DynamicDataPool) Alloc(len uint) *Data {
	for i, pool := range p.pools {
		if len <= p.sizes[i] {
			return pool.Alloc(len)
		}
	}
	return nil
}

func (p *DynamicDataPool) AllocCap(len, cap uint) *Data {
	for i, pool := range p.pools {
		if len <= p.sizes[i] {
			return pool.AllocCap(len, cap)
		}
	}
	return nil
}
