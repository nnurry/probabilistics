package arch

import (
	"math"
	"strconv"
)

const PtrSize = 32 << uintptr(^uintptr(0)>>63)
const IntSize = strconv.IntSize

var Log2IntSize = uint(math.Log2(IntSize)) // log2(32) or log2(64) -> assert type = int
