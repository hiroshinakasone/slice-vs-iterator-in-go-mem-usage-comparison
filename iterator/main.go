package main

import (
	"iter"
	"math"
	"runtime"
	"fmt"
)

const (
	NumOfElements = 10_000_000
)

func sliceNumbers(n int) []int {
	s := make([]int, n)
	for i := range len(s) {
		s[i] = i
	}
	return s
}

func iterFilter(seq []int, predicate func(int) bool) iter.Seq[int] {
	return func(yield func(int) bool) {
		for i := range seq {
			if predicate(i) {
				if !yield(i) {
					return
				}
			}
		}
	}
}

func benchmarkIterator(n int) {
	numbers := sliceNumbers(n)
	even := iterFilter(numbers, func(i int) bool { return i%2 == 0 })
	for i := range even {
		_ = i
	}
}

func main() {
	var aveCount uint64
	var aveTotalAlloc uint64
	var maxTotalAlloc uint64 = 0
	var minTotalAlloc uint64 = math.MaxUint64
	var aveMallocs uint64
	var maxMallocs uint64 = 0
	var minMallocs uint64 = math.MaxUint64

	var m1, m2 runtime.MemStats
	for i := 1; i < 10_000; i++ {
		runtime.GC()
		runtime.ReadMemStats(&m1)
		benchmarkIterator(NumOfElements)
		runtime.ReadMemStats(&m2)

		diffTotalAlloc := m2.TotalAlloc - m1.TotalAlloc
		diffMallocs := m2.Mallocs - m1.Mallocs

		aveTotalAlloc = (aveTotalAlloc*aveCount + diffTotalAlloc) / uint64(i)
		maxTotalAlloc = max(maxTotalAlloc, diffTotalAlloc)
		minTotalAlloc = min(minTotalAlloc, diffTotalAlloc)

		aveMallocs = (aveMallocs*aveCount + diffMallocs) / uint64(i)
		maxMallocs = max(maxMallocs, diffMallocs)
		minMallocs = min(minMallocs, diffMallocs)

		aveCount++
	}

	fmt.Printf("NumOfElements: %d\n", NumOfElements)
	fmt.Println()
	fmt.Printf("AveTotalAlloc: %d B/op\n", aveTotalAlloc)
	fmt.Printf("MaxTotalAlloc: %d B\n", maxTotalAlloc)
	fmt.Printf("MinTotalAlloc: %d B\n", minTotalAlloc)
	fmt.Println()
	fmt.Printf("AveMallocs: %d allocs/op\n", aveMallocs)
	fmt.Printf("MaxMallocs: %d allocs\n", maxMallocs)
	fmt.Printf("MinMallocs: %d allocs\n", minMallocs)
}
