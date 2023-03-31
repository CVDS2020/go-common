package uns

import (
	"gitee.com/sy_183/common/generic"
	"gitee.com/sy_183/common/uns/goarch"
	"gitee.com/sy_183/common/uns/goos"
	"gitee.com/sy_183/common/uns/math"
	"unsafe"
)

const (
	// _64bit = 1 on 64-bit systems, 0 on 32-bit systems
	_64bit = 1 << (^uintptr(0) >> 63) / 2

	// heapAddrBits is the number of bits in a heap address. On
	// amd64, addresses are sign-extended beyond heapAddrBits. On
	// other arches, they are zero-extended.
	//
	// On most 64-bit platforms, we limit this to 48 bits based on a
	// combination of hardware and OS limitations.
	//
	// amd64 hardware limits addresses to 48 bits, sign-extended
	// to 64 bits. Addresses where the top 16 bits are not either
	// all 0 or all 1 are "non-canonical" and invalid. Because of
	// these "negative" addresses, we offset addresses by 1<<47
	// (arenaBaseOffset) on amd64 before computing indexes into
	// the heap arenas index. In 2017, amd64 hardware added
	// support for 57 bit addresses; however, currently only Linux
	// supports this extension and the kernel will never choose an
	// address above 1<<47 unless mmap is called with a hint
	// address above 1<<47 (which we never do).
	//
	// arm64 hardware (as of ARMv8) limits user addresses to 48
	// bits, in the range [0, 1<<48).
	//
	// ppc64, mips64, and s390x support arbitrary 64 bit addresses
	// in hardware. On Linux, Go leans on stricter OS limits. Based
	// on Linux's processor.h, the user address space is limited as
	// follows on 64-bit architectures:
	//
	// Architecture  Name              Maximum Value (exclusive)
	// ---------------------------------------------------------------------
	// amd64         TASK_SIZE_MAX     0x007ffffffff000 (47 bit addresses)
	// arm64         TASK_SIZE_64      0x01000000000000 (48 bit addresses)
	// ppc64{,le}    TASK_SIZE_USER64  0x00400000000000 (46 bit addresses)
	// mips64{,le}   TASK_SIZE64       0x00010000000000 (40 bit addresses)
	// s390x         TASK_SIZE         1<<64 (64 bit addresses)
	//
	// These limits may increase over time, but are currently at
	// most 48 bits except on s390x. On all architectures, Linux
	// starts placing mmap'd regions at addresses that are
	// significantly below 48 bits, so even if it's possible to
	// exceed Go's 48 bit limit, it's extremely unlikely in
	// practice.
	//
	// On 32-bit platforms, we accept the full 32-bit address
	// space because doing so is cheap.
	// mips32 only has access to the low 2GB of virtual memory, so
	// we further limit it to 31 bits.
	//
	// On ios/arm64, although 64-bit pointers are presumably
	// available, pointers are truncated to 33 bits in iOS <14.
	// Furthermore, only the top 4 GiB of the address space are
	// actually available to the application. In iOS >=14, more
	// of the address space is available, and the OS can now
	// provide addresses outside of those 33 bits. Pick 40 bits
	// as a reasonable balance between address space usage by the
	// page allocator, and flexibility for what mmap'd regions
	// we'll accept for the heap. We can't just move to the full
	// 48 bits because this uses too much address space for older
	// iOS versions.
	// TODO(mknyszek): Once iOS <14 is deprecated, promote ios/arm64
	// to a 48-bit address space like every other arm64 platform.
	//
	// WebAssembly currently has a limit of 4GB linear memory.
	heapAddrBits = (_64bit*(1-goarch.IsWasm)*(1-goos.IsIos*goarch.IsArm64))*48 + (1-_64bit+goarch.IsWasm)*(32-(goarch.IsMips+goarch.IsMipsle)) + 40*goos.IsIos*goarch.IsArm64

	// maxAlloc is the maximum size of an allocation. On 64-bit,
	// it's theoretically possible to allocate 1<<heapAddrBits bytes. On
	// 32-bit, however, this is one less than 1<<32 because the
	// number of bytes in the address space doesn't actually fit
	// in a uintptr.
	maxAlloc = (1 << heapAddrBits) - (1-_64bit)*1
)

//go:linkname panicmakeslicelen runtime.panicmakeslicelen
func panicmakeslicelen()

//go:linkname panicmakeslicecap runtime.panicmakeslicecap
func panicmakeslicecap()

type SliceStruct struct {
	Ptr unsafe.Pointer
	Len int
	Cap int
}

func SliceStructOf[E any](s []E) SliceStruct {
	return *(*SliceStruct)(unsafe.Pointer(&s))
}

func SliceStructOfPtr[E any](sp *[]E) *SliceStruct {
	return (*SliceStruct)(unsafe.Pointer(sp))
}

func SliceStructToSlice[E any](ss SliceStruct) []E {
	return *(*[]E)(unsafe.Pointer(&ss))
}

func MakeSliceUnchecked[E any](ptr unsafe.Pointer, len, cap int) (es []E) {
	*ConvertPointer[[]E, SliceStruct](&es) = SliceStruct{Ptr: ptr, Len: len, Cap: cap}
	return
}

func MakeSlice[E any](ptr unsafe.Pointer, len, cap int) []E {
	es := generic.Size[E]()
	mem, overflow := math.MulUintptr(es, uintptr(cap))
	if overflow || mem > maxAlloc || len < 0 || len > cap {
		// NOTE: Produce a 'len out of range' error instead of a
		// 'cap out of range' error when someone does make([]T, bignumber).
		// 'cap out of range' is true too, but since the cap is only being
		// supplied implicitly, saying len is clearer.
		// See golang.org/issue/4085.
		mem, overflow := math.MulUintptr(es, uintptr(len))
		if overflow || mem > maxAlloc || len < 0 {
			panicmakeslicelen()
		}
		panicmakeslicecap()
	}

	return MakeSliceUnchecked[E](ptr, len, cap)
}

func SlicePointer[E any](s []E) unsafe.Pointer {
	return SliceStructOfPtr(&s).Ptr
}

func SliceConvert[EI, EO any](s []EI, len, cap int) []EO {
	return MakeSlice[EO](SlicePointer(s), len, cap)
}

func SliceMerge[E any](s1, s2 []E) []E {
	var e E
	es := unsafe.Sizeof(e)
	ss1 := SliceStructOfPtr(&s1)
	s1Ptr := uintptr(ss1.Ptr)
	s2Ptr := uintptr(SliceStructOfPtr(&s2).Ptr)
	if s1Ptr+uintptr(len(s1))*es == s2Ptr {
		return MakeSliceUnchecked[E](ss1.Ptr, len(s1)+len(s2), len(s1)+cap(s2))
	}
	return nil
}
