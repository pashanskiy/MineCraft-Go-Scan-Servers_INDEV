package main

import (
	"encoding/hex"
	"fmt"
	"testing"
)

// func Benchmark_test(b *testing.B) {
// 	var i uint64
// 	for i = 0; i < b.N; i++ {
// 		// x = 127
// 		// _ = x
// 		// x = 0
// 		// _ = x
// 		// // fmt.Println(x)
// 		encodeVarint(i)
// 	}

// }

func Test_test(t *testing.T) {

	fmt.Println(hex.EncodeToString(toVarintTest(int32(len("ip")))))

}

func toVarintTest(x int32) []byte {
	var bytes [10]byte
	var n int
	for n = 0; x > 127; n++ {
		bytes[n] = 0x80 | uint8(x&0x7F)
		x >>= 7
	}
	bytes[n] = uint8(x)
	n++
	return bytes[0:n]
}

func test_Kek(t *testing.T) {
	x := 8 | 5
	fmt.Println(x)
}
