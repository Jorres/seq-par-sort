package main

import (
	"math"
	"math/rand"
	"sync"
)

const parBlockSize = 2000
const scanBlockSize = 2000

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func sum(a, b int) int {
	return a + b
}

func parFor(n int, f func(int)) {
	var wg sync.WaitGroup
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func(curBlock int) {
			f(curBlock)
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func parScan(a []int, l, r int, f func(int, int) int, startVal int) []int {
	if r-l < scanBlockSize {
		ans := make([]int, r-l)

		curVal := startVal
		for i := l; i < r; i++ {
			curVal = f(curVal, a[i])
			ans[i-l] = curVal
		}

		return ans
	}

	blocks := int(math.Ceil(float64(r-l) / scanBlockSize))
	sums := make([]int, blocks)

	parFor(blocks, func(curBlock int) {
		curBlockVal := 0
		for k := l + curBlock*scanBlockSize; k < min(l+(curBlock+1)*scanBlockSize, r); k++ {
			curBlockVal = f(curBlockVal, a[k])
		}
		sums[curBlock] = curBlockVal
	})

	sums = parScan(sums, 0, len(sums), sum, 0) // span X/blockSize
	ans := make([]int, r-l)

	parFor(blocks, func(curBlock int) {
		curBlockVal := 0
		if curBlock > 0 {
			curBlockVal = sums[curBlock-1]
		}

		for k := l + curBlock*scanBlockSize; k < min(l+(curBlock+1)*scanBlockSize, r); k++ {
			curBlockVal = f(curBlockVal, a[k])
			ans[k-l] = curBlockVal
		}
	})

	return ans
}

func parMap(a, b []int, l, r int, f func(int) int) {
	if r-l < parBlockSize {
		for i := l; i < r; i++ {
			b[i] = f(a[i])
		}
		return
	}

	m := l + (r-l)/2

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		parMap(a, b, l, m, f)
		wg.Done()
	}()
	parMap(a, b, m, r, f)
	wg.Wait()
}

func parFilter(a []int, l, r int, f func(int) bool) []int {
	if r-l < parBlockSize {
		ans := make([]int, parBlockSize / 2)
		for i := l; i < r; i++ {
			if f(a[i]) {
				ans = append(ans, a[i])
			}
		}
		return ans
	}

	flags := make([]int, r-l)
	parMap(a, flags, l, r, func(x int) int {
		if f(x) {
			return 1
		}
		return 0
	})

	blocks := int(math.Ceil(float64(r-l) / parBlockSize))
	sums := make([]int, blocks)

	parFor(blocks, func(curBlock int) {
		curBlockVal := 0
		for k := l + curBlock*parBlockSize; k < min(l+(curBlock+1)*parBlockSize, r); k++ {
			curBlockVal = curBlockVal + flags[k]
		}
		sums[curBlock] = curBlockVal
	})

	sums = parScan(sums, 0, len(sums), sum, 0)

	// I used to do this for a separate answer,
	// but reusing flags allocation makes sense
	// ans := make([]int, sums[len(sums)-1])
	parFor(blocks, func(curBlock int) {
		shift := 0
		if curBlock > 0 {
			shift = sums[curBlock-1]
		}
		lastWritten := shift
		for k := l + curBlock*parBlockSize; k < min(l+(curBlock+1)*parBlockSize, r); k++ {
			if flags[k] == 1 {
				if (k < lastWritten) {
					panic("Aaaa, reusing flags turned out wrong")
				}
				flags[lastWritten] = a[k]
				lastWritten += 1
			}
		}
	})

	return flags
}

func parCopy(to []int, from []int, startPos int, l, r int) {
	if r-l < parBlockSize {
		for i := l; i < r; i++ {
			to[startPos+i-l] = from[i]
		}
		return
	}

	m := l + (r-l)/2

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		parCopy(to, from, startPos, l, m)
		wg.Done()
	}()
	parCopy(to, from, startPos+(m-l), m, r)
	wg.Wait()
}

func parQSort(a []int, l, _r int) {
	r := _r + 1
	if r-l < parBlockSize {
		QsortSeq(a, l, r-1)
		return
	}

	m := rand.Intn(r-l) + l

	left := parFilter(a, l, r, func(x int) bool {
		return x < a[m]
	})
	middle := parFilter(a, l, r, func(x int) bool {
		return x == a[m]
	})
	right := parFilter(a, l, r, func(x int) bool {
		return x > a[m]
	})

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		parQSort(left, 0, len(left)-1)
		wg.Done()
	}()
	parQSort(right, 0, len(right)-1)
	wg.Wait()

	wg.Add(2)
	go func() {
		parCopy(a, left, 0, 0, len(left))
		wg.Done()
	}()
	go func() {
		parCopy(a, middle, len(left), 0, len(middle))
		wg.Done()
	}()
	parCopy(a, right, len(left)+len(middle), 0, len(right))
	wg.Wait()
}
