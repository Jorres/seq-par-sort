package main

import (
	"reflect"
	"testing"
)

func TestParScan(t *testing.T) {
	tests := []struct {
		a, want []int
	}{
		{[]int{0, 1, 2, 3}, []int{0, 1, 3, 6}},
		{[]int{1, 1, 1, 1, 1, 1, 1, 1, 1}, []int{1, 2, 3, 4, 5, 6, 7, 8, 9}},
		{[]int{0, 1, 2, 3, 5}, []int{0, 1, 3, 6, 11}},
		{[]int{0, -1, 0, 1, 0}, []int{0, -1, -1, 0, 0}},
		{[]int{100}, []int{100}},
	}

	for _, tt := range tests {
		if got := parScan(tt.a, 0, len(tt.a), sum, 0); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("parScan(%v) = %v; want %v", tt.a, got, tt.want)
		}
	}
}

func TestParMap(t *testing.T) {
	mult2 := func(x int) int {
		return x * 2
	}

	tests := []struct {
		a, want []int
	}{
		{[]int{0, 1, 2, 3}, []int{0, 2, 4, 6}},
		{[]int{0, 1, 2, 3, 5}, []int{0, 2, 4, 6, 10}},
		{[]int{0, -1, 0, 1, 0}, []int{0, -2, 0, 2, 0}},
		{[]int{100}, []int{200}},
	}

	for _, tt := range tests {
		ans := make([]int, len(tt.a))
		parMap(tt.a, ans, 0, len(tt.a), mult2)
		if !reflect.DeepEqual(ans, tt.want) {
			t.Errorf("parMap(%v) = %v; want %v", tt.a, ans, tt.want)
		}
	}
}

func TestParFilter(t *testing.T) {
	isEven := func(x int) bool {
		return x%2 == 0
	}
	lessThan1 := func(x int) bool {
		return x < 1
	}

	tests := []struct {
		a, want []int
		f       func(int) bool
	}{
		{[]int{0, 1, 2, 3}, []int{0, 2}, isEven},
		{[]int{0, 1, 2, 3, 5}, []int{0, 2}, isEven},
		{[]int{0, -1, 0, 1, 0}, []int{0, 0, 0}, isEven},
		{[]int{100}, []int{100}, isEven},
		{[]int{2, 1, 0}, []int{0}, lessThan1},
		{[]int{2, 1, 3, 0}, []int{0}, lessThan1},
	}

	for _, tt := range tests {
		ans := parFilter(tt.a, 0, len(tt.a), tt.f)
		if !reflect.DeepEqual(ans, tt.want) {
			t.Errorf("parFilter(%v) = %v; want %v", tt.a, ans, tt.want)
		}
	}
}

func TestParCopy(t *testing.T) {
	tests := []struct {
		a    []int
		l, r int
		want []int
	}{
		{[]int{0, 1, 2, 3}, 0, 4, []int{0, 1, 2, 3}},
		{[]int{0, 1, 2, 3}, 1, 3, []int{1, 2}},
		{[]int{0, 1, 2, 3}, 2, 3, []int{2}},
		{[]int{0, 1, 2, 3}, 3, 3, []int{}},
	}

	for _, tt := range tests {
		to := make([]int, tt.r-tt.l)
		parCopy(to, tt.a, 0, tt.l, tt.r)
		if !reflect.DeepEqual(to, tt.want) {
			t.Errorf("parFilter(%v) = %v; want %v", tt.a, to, tt.want)
		}
	}
}

func TestParQSort(t *testing.T) {
	// parQSort
	tests := []struct {
		a, want []int
	}{
		{[]int{2, 1, 3, 0}, []int{0, 1, 2, 3}},
		{[]int{1, 1, 1}, []int{1, 1, 1}},
		{[]int{100, 10, 123, -1, 0, 0, 0, 0, 0, 1, 1, 1}, []int{-1, 0, 0, 0, 0, 0, 1, 1, 1, 10, 100, 123}},
	}

	for _, tt := range tests {
		parQSort(tt.a, 0, len(tt.a)-1)
		if !reflect.DeepEqual(tt.a, tt.want) {
			t.Errorf("parQSort(..) = %v; want %v", tt.a, tt.want)
		}
	}
}
