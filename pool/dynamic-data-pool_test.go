package pool

import (
	"fmt"
	"gitee.com/sy_183/common/unit"
	"testing"
)

func TestDynamicDataPool(t *testing.T) {
	pool := NewDynamicDataPoolWithExp(64*unit.KiBiByte, unit.MeBiByte, "sync")
	fmt.Println(pool.pools)
}
