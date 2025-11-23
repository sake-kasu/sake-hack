package utils

import (
	"fmt"
	"math"
)

// IntToInt32 は int を int32 に安全に変換する
// CWE-190 Integer Overflow 対策
func IntToInt32(value int) (int32, error) {
	if value > math.MaxInt32 || value < math.MinInt32 {
		return 0, fmt.Errorf("integer overflow: value %d is out of int32 range", value)
	}
	return int32(value), nil
}

// Int64ToInt32 は int64 を int32 に安全に変換する
// CWE-190 Integer Overflow 対策
func Int64ToInt32(value int64) (int32, error) {
	if value > math.MaxInt32 || value < math.MinInt32 {
		return 0, fmt.Errorf("integer overflow: value %d is out of int32 range", value)
	}
	return int32(value), nil
}

// IntToUint32 は int を uint32 に安全に変換する
func IntToUint32(value int) (uint32, error) {
	if value < 0 || value > math.MaxUint32 {
		return 0, fmt.Errorf("integer overflow: value %d is out of uint32 range", value)
	}
	return uint32(value), nil
}

// Int64ToUint32 は int64 を uint32 に安全に変換する
func Int64ToUint32(value int64) (uint32, error) {
	if value < 0 || value > math.MaxUint32 {
		return 0, fmt.Errorf("integer overflow: value %d is out of uint32 range", value)
	}
	return uint32(value), nil
}
