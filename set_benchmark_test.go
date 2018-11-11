package ribbons

import (
	"testing"
)

func BenchmarkSetHit_1(b *testing.B)       { benchSet(b, 1, 0) }
func BenchmarkMapHit_1(b *testing.B)       { benchMap(b, 1, 0) }
func BenchmarkSetHit_10(b *testing.B)      { benchSet(b, 10, 8) }
func BenchmarkMapHit_10(b *testing.B)      { benchMap(b, 10, 8) }
func BenchmarkSetHit_100(b *testing.B)     { benchSet(b, 100, 88) }
func BenchmarkMapHit_100(b *testing.B)     { benchMap(b, 100, 88) }
func BenchmarkSetHit_1000(b *testing.B)    { benchSet(b, 1000, 888) }
func BenchmarkMapHit_1000(b *testing.B)    { benchMap(b, 1000, 888) }
func BenchmarkSetHit_100000(b *testing.B)  { benchSet(b, 100000, 88888) }
func BenchmarkMapHit_100000(b *testing.B)  { benchMap(b, 100000, 88888) }
func BenchmarkSetHit_1000000(b *testing.B) { benchSet(b, 1000000, 888888) }
func BenchmarkMapHit_1000000(b *testing.B) { benchMap(b, 1000000, 888888) }

func BenchmarkSetMiss_1(b *testing.B)       { benchSet(b, 1, 0, 0) }
func BenchmarkMapMiss_1(b *testing.B)       { benchMap(b, 1, 0, 0) }
func BenchmarkSetMiss_10(b *testing.B)      { benchSet(b, 10, 8, 8) }
func BenchmarkMapMiss_10(b *testing.B)      { benchMap(b, 10, 8, 8) }
func BenchmarkSetMiss_100(b *testing.B)     { benchSet(b, 100, 88, 88) }
func BenchmarkMapMiss_100(b *testing.B)     { benchMap(b, 100, 88, 88) }
func BenchmarkSetMiss_1000(b *testing.B)    { benchSet(b, 1000, 888, 888) }
func BenchmarkMapMiss_1000(b *testing.B)    { benchMap(b, 1000, 888, 888) }
func BenchmarkSetMiss_100000(b *testing.B)  { benchSet(b, 100000, 88888, 88888) }
func BenchmarkMapMiss_100000(b *testing.B)  { benchMap(b, 100000, 88888, 88888) }
func BenchmarkSetMiss_1000000(b *testing.B) { benchSet(b, 1000000, 888888, 888888) }
func BenchmarkMapMiss_1000000(b *testing.B) { benchMap(b, 1000000, 888888, 888888) }

func newSet(n uint64, miss ...uint64) *UINT64Set {
	s := &UINT64Set{}
	for i := uint64(0); i < n; i++ {
		if len(miss) > 0 && miss[0] == i {
			continue
		}
		s.Add(i)
	}
	return s
}

func newMap(n uint64, miss ...uint64) map[uint64]struct{} {
	s := make(map[uint64]struct{}, n)
	for i := uint64(0); i < n; i++ {
		if len(miss) > 0 && miss[0] == i {
			continue
		}
		s[i] = struct{}{}
	}
	return s
}

func benchSet(b *testing.B, size, check uint64, miss ...uint64) {
	s := newSet(size, miss...)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = s.Has(check)
	}
}

func benchMap(b *testing.B, size, check uint64, miss ...uint64) {
	s := newMap(size, miss...)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = s[check]
	}
}
