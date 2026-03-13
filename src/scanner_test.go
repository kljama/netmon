package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
)

func BenchmarkActiveCount_SyncMapRange(b *testing.B) {
	sizes := []int{100, 1000, 10000}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("%d", size), func(b *testing.B) {
			var m sync.Map
			for i := 0; i < size; i++ {
				m.Store(i, true)
			}
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				count := 0
				m.Range(func(key, value interface{}) bool {
					count++
					return true
				})
			}
		})
	}
}

func BenchmarkActiveCount_AtomicCounter(b *testing.B) {
	sizes := []int{100, 1000, 10000}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("%d", size), func(b *testing.B) {
			var count atomic.Int64
			count.Store(int64(size))
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = count.Load()
			}
		})
	}
}
