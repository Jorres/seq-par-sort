package main

import (
	"fmt"
	// "log"
	"math/rand"
	// "os"
	// "runtime/pprof"
	"sync"
	"time"
)

const arraySize = 1e8
const stopForkingSize = 2000
const startBubbleSortSize = 20

func swap(arr []int, i, j int) {
	arr[i], arr[j] = arr[j], arr[i]
}

func doPartition(arr []int, low, high int) int {
	m := arr[high]
	i := low

	for j := low; j < high; j++ {
		if arr[j] < m {
			swap(arr, i, j)
			i++
		}
	}
	swap(arr, i, high)
	return i
}

func bubbleSort(arr []int, i, j int) {
	for end := j + 1; end > i; end-- {
		swapped := false
		for current := i; current < end-1; current++ {
			if arr[current] > arr[current+1] {
				arr[current], arr[current+1] = arr[current+1], arr[current]
				swapped = true
			}
		}
		if !swapped {
			return
		}
	}
}

func qsortSeq(arr []int, low, high int) {
	if high-low < startBubbleSortSize {
		bubbleSort(arr, low, high)
		return
	}
	if low < high {
		pivot := doPartition(arr, low, high)
		qsortSeq(arr, low, pivot-1)
		qsortSeq(arr, pivot+1, high)
	}
}

func qsortPar(arr []int, low, high int) {
	if high-low < stopForkingSize {
		qsortSeq(arr, low, high)
		return
	}

	if low < high {
		pivot := doPartition(arr, low, high)
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			qsortPar(arr, low, pivot-1)
			wg.Done()
		}()
		qsortPar(arr, pivot+1, high)
		wg.Wait()
	}
}

func generateRandomArray(size int) []int {
	array := make([]int, size)
	for i := range array {
		array[i] = rand.Intn(size) // up to size, why not. It's arbitrary.
	}
	return array
}

func isSortedAsc(arr []int) bool {
	for i := 0; i < len(arr)-1; i++ {
		if arr[i] > arr[i+1] {
			return false
		}
	}
	return true
}

func doTest(sortingFunc func([]int, int, int), testName string) {
	const nTests = 5
	var totalTime time.Duration

	fmt.Printf("%v with %v elements, averaged over %v launches\n\n", testName, arraySize, nTests)

	for i := 0; i < nTests; i++ {
		arr := generateRandomArray(arraySize)

		start := time.Now()

		// f, err := os.Create("cpu.prof")
		// if err != nil {
		// 	log.Fatal("could not create CPU profile: ", err)
		// }
		// defer f.Close() // error handling omitted for example
		// if err := pprof.StartCPUProfile(f); err != nil {
		// 	log.Fatal("could not start CPU profile: ", err)
		// }

		sortingFunc(arr, 0, len(arr)-1)

		// pprof.StopCPUProfile()

		elapsed := time.Since(start)
		totalTime += elapsed

		if !isSortedAsc(arr) {
			fmt.Println("Error: array is not sorted properly.")
			return
		}

		fmt.Printf("Launch %v: %v\n", i+1, elapsed)
	}

	avgTime := totalTime / nTests
	fmt.Printf("\nAverage time: %v\n", avgTime)
}

func main() {
	doTest(qsortPar, "Quicksort parallel")
	// doTest(qsortSeq, "Quicksort sequential")
}
