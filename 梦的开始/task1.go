package main

import (
	"fmt"
	"strconv"
)

func main() {
	// 1. single_number test
	singleNumberResult := single_number([]int{4, 1, 2, 1, 2})
	fmt.Println("单一数字是:", singleNumberResult)

	// 2. isPalindrome test
	num2 := 1234321
	palindromeResult := isPalindrome(123321)
	fmt.Printf("数字: %d 是回文数：%t\n", num2, palindromeResult)

	// 3. isValidParentheses test
	str3 := "[((](){}"
	validParenthesesResult := isValidParentheses(str3)
	fmt.Printf("字符串 %s 中的括号是有效的：%t\n", str3, validParenthesesResult)

	// 4. longestCommonPrefix test
	longestCommonPrefixResult := longestCommonPrefix([]string{"flower", "flow", "flight", "flag"})
	fmt.Printf("最长公共字符串为: %s\n", longestCommonPrefixResult)

	// 5. remove_duplicates_from_sorted_array test
	duplicateArray := []int{1, 1, 2, 2, 3, 4, 4}
	removeDuplicatesResult := remove_duplicates_from_sorted_array(duplicateArray)
	fmt.Printf("去重后的数组长度为: %d\n", removeDuplicatesResult)

	// 6. plusOne test
	plusOneArray := []int{1, 2, 3}
	plusOneResult := plusOne(plusOneArray)
	fmt.Printf("加一后的数组为: %v\n", plusOneResult)

	// 7. merge_intervals test
	intervals := [][]int{{1, 3}, {2, 6}, {8, 10}, {15, 18}, {17, 20}}
	mergedIntervalsResult := merge_intervals(intervals)
	fmt.Printf("合并后的区间为: %v\n", mergedIntervalsResult)

	// 8. two_sum test
	nums := []int{2, 7, 11, 15}
	target := 9
	twoSumResult := two_sum(nums, target)
	fmt.Printf("和为 %d 的两个数的索引为: %v\n", target, twoSumResult)
}

func single_number(nums []int) int {
	var m map[int]int = make(map[int]int)

	for _, num := range nums {
		m[num] = m[num] + 1
	}

	for num, count := range m {
		if count == 1 {
			return num
		}
	}

	return 0
}

func isPalindrome(x int) bool {
	var s string = strconv.Itoa(x)
	var runes []rune = []rune(s)
	for i := 0; i < len(runes)/2; i++ {
		if runes[i] != runes[len(runes)-1-i] {
			return false
		}
	}
	return true
}

func isValidParentheses(s string) bool {
	var stack []rune

	for _, char := range s {
		if char == '(' || char == '{' || char == '[' {
			stack = append(stack, char)
		} else {
			if len(stack) == 0 {
				return false
			}
			top := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			if (char == ')' && top != '(') ||
				(char == '}' && top != '{') ||
				(char == ']' && top != '[') {
				return false
			}
		}
	}

	return len(stack) == 0
}

func longestCommonPrefix(strs []string) string {
	if len(strs) == 0 {
		return ""
	}
	prefix := strs[0]
	for _, str := range strs[1:] {
		for len(str) < len(prefix) || str[:len(prefix)] != prefix {
			prefix = prefix[:len(prefix)-1]
			if len(prefix) == 0 {
				return ""
			}
		}
	}
	return prefix
}

func remove_duplicates_from_sorted_array(nums []int) int {
	if len(nums) == 0 {
		return 0
	}

	writeIndex := 1

	for i := 1; i < len(nums); i++ {
		if nums[i] != nums[i-1] {
			nums[writeIndex] = nums[i]
			writeIndex++
		}
	}

	return writeIndex
}

func plusOne(nums []int) []int {
	for i := len(nums) - 1; i >= 0; i-- {
		if nums[i] < 9 {
			nums[i]++
			return nums
		}
		nums[i] = 0
	}
	return append([]int{1}, nums...)
}

func merge_intervals(intervals [][]int) [][]int {
	if len(intervals) == 0 {
		return intervals
	}

	for i := 0; i < len(intervals); i++ {
		for j := i + 1; j < len(intervals); j++ {
			if intervals[i][0] > intervals[j][0] {
				temp := intervals[i]
				intervals[i] = intervals[j]
				intervals[j] = temp
			}
		}
	}

	merged := [][]int{intervals[0]}
	for i := 1; i < len(intervals); i++ {
		if intervals[i][0] <= merged[len(merged)-1][1] {
			merged[len(merged)-1] = []int{merged[len(merged)-1][0], intervals[i][1]}
		} else {
			merged = append(merged, intervals[i])
		}
	}

	return merged
}

// 给定一个整数数组 nums 和一个整数目标值 target，请你在该数组中找出 和为目标值 target  的那 两个 整数，并返回它们的数组下标。
func two_sum(nums []int, target int) []int {

	m := make(map[int]int)

	for i, num := range nums {
		value, ok := m[target-num]
		if ok {
			return []int{value, i}
		}
		m[num] = i
	}
	return nil
}
