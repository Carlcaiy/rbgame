package main

import (
	"fmt"
	"sort"
	"testing"
)

func TestSlotLine(t *testing.T) {
	start()
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func dif_midd(arr []int) int {
	sum := 0
	for i := range arr {
		if i == 0 {
			sum += abs(arr[i] - half)
		} else {
			sum += abs(arr[i] - arr[i-1])
		}
	}
	return sum
}

func dif_head(arr []int) int {
	sum := 0
	for i := range arr {
		if i > 0 {
			sum += abs(arr[i]-arr[i-1]) * abs(i-len(arr)/2) * 10
		}
	}
	return sum
}

func key(arr []int) int {
	key := (arr[0] << 16) | (arr[1] << 12) | (arr[2] << 8) | (arr[3] << 4) | arr[4]
	return key
}

func dasdadadasd(key int) []int {
	array := make([]int, 5)
	array[4] = key & 0xf
	array[3] = (key >> 4) & 0xf
	array[2] = (key >> 8) & 0xf
	array[1] = (key >> 12) & 0xf
	array[0] = (key >> 16) & 0xf
	return array
}

func ldc(k int) int {
	arr := dasdadadasd(k)
	for i := range arr {
		arr[i] = height - arr[i] - 1
	}
	key := (arr[0] << 16) | (arr[1] << 12) | (arr[2] << 8) | (arr[3] << 4) | arr[4]
	return key
}

const height int = 3

var half int = height / 2

const midline int = 0x11111

func start() {
	var lines = make(map[int]bool)
	for i := 0; i <= half; i++ {
		for j := 0; j < height; j++ {
			if i == half && j > half {
				continue
			}
			for k := 0; k < height; k++ {
				if i == half && j == half && k > half {
					continue
				}
				for l := 0; l < height; l++ {
					if i == half && j == half && k == half && l > half {
						continue
					}
					for m := 0; m < height; m++ {
						if i == half && j == half && k == half && l == half && m > half {
							continue
						}
						lines[key([]int{i, j, k, l, m})] = true
					}
				}
			}
		}
	}

	array := make([][]int, 0, 125)
	for k := range lines {
		arr := dasdadadasd(k)
		if arr[0] == arr[4] && arr[1] == arr[3] {
			array = append(array, dasdadadasd(k))
			lines[k] = false
		}
	}

	sort.Slice(array, func(i int, j int) bool {
		if dif_midd(array[i]) == dif_midd(array[j]) {
			return dif_head(array[i]) < dif_head(array[j])
		}
		return dif_midd(array[i]) < dif_midd(array[j])
	})

	count := 0
	for i, v := range array {
		if (i-1)%10 == 0 {
			fmt.Println()
		}
		if key(v) != midline {
			fmt.Printf("0x%05x, 0x%05x, ", key(v), ldc(key(v)))
			count += 2
		} else {
			fmt.Printf("0x%05x", key(v))
			count += 1
		}
	}
	fmt.Println("\n左右对称 count:", count)

	array = make([][]int, 0, 24)
	for k := range lines {
		if lines[k] {
			arr := dasdadadasd(k)
			if arr[0]+arr[4] == (height-1) && arr[1]+arr[3] == (height-1) && arr[2] == half {
				array = append(array, arr)
				lines[k] = false
			}
		}
	}

	sort.Slice(array, func(i int, j int) bool {
		if dif_midd(array[i]) == dif_midd(array[j]) {
			return dif_head(array[i]) < dif_head(array[j])
		}
		return dif_midd(array[i]) < dif_midd(array[j])
	})

	count = 0
	for i, v := range array {
		if (i-1)%10 == 0 {
			fmt.Println()
		}
		if key(v) != midline {
			fmt.Printf("0x%05x, 0x%05x, ", key(v), ldc(key(v)))
			count += 2
		} else {
			fmt.Printf("0x%05x\n", key(v))
			count += 1
		}
	}

	fmt.Println("中心对称 count:", count)

	count = 0
	array = make([][]int, 0, 1488)
	for k, v := range lines {
		if v {
			array = append(array, dasdadadasd(k))
		}
	}
	sort.Slice(array, func(i int, j int) bool {
		if dif_midd(array[i]) == dif_midd(array[j]) {
			return dif_head(array[i]) < dif_head(array[j])
		}
		return dif_midd(array[i]) < dif_midd(array[j])
	})
	for i, v := range array {
		if i%20 == 0 {
			fmt.Println()
		}
		fmt.Printf("0x%05x, 0x%05x, ", key(v), ldc(key(v)))
		count += 2
	}
	fmt.Println("没对称 count:", count)
}
