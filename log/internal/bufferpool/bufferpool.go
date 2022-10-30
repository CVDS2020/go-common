package bufferpool

import (
	"gitee.com/sy_183/common/pool"
)

var (
	_pool = pool.NewBufferPool(1024)
	// Get retrieves a buffer from the pool, creating one if necessary.
	Get = _pool.Get
)
