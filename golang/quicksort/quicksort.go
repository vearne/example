package main

import "fmt"

func sort(nums []int) {
	if len(nums) <= 1 {
		return
	}
	var i, j, val int
	val = nums[0]
	j = len(nums) - 1

	for i < j {
		for nums[i] < val {
			i++
		}
		for nums[j] > val {
			j--
		}
		if nums[j] < val {
			nums[i], nums[j] = nums[j], nums[i]
		}
		if nums[i] > val {
			nums[i], nums[j] = nums[j], nums[i]
		}
	}

	sort(nums[0:i])
	sort(nums[i+1:])
}

func main() {
	nums := []int{}
	sort(nums)
	fmt.Println(nums)
}
