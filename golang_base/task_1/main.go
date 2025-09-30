package main

import (
	"fmt"
)

func AnswerQuesion1(nums []int) (bool, int) {
	// 只出现一次的数字：给定一个非空整数数组，除了某个元素只出现一次以外，其余每个元素均出现两次。
	// 找出那个只出现了一次的元素。可以使用 for 循环遍历数组，结合 if 条件判断和 map 数据结构来解决，
	// 例如通过 map 记录每个元素出现的次数，然后再遍历 map 找到出现次数为1的元素。
	fmt.Println(" <----------------------- AnswerQuesion1 -------------------->")
	counter_map := make(map[int]int)
	for _, v := range nums {
		counter_map[v]++
	}
	for k, v := range counter_map {
		if v == 1 {
			fmt.Println("只出现一次的数字是:", k)
			return true, k
		}
	}
	fmt.Println("没有只出现一次的数字")
	return false, 0
}

func AnswerQuesion2(num int) bool {
	// 判断一个整数是否是回文数.  121 是回文数，而 -121 不是。
	// 你可以将整数转换为字符串，然后使用双指针法从字符串的两端向中间移动，比较对应位置的字符是否相同。
	// 如果所有对应位置的字符都相同，则该整数是回文数；否则不是。
	fmt.Println(" <----------------------- AnswerQuesion2 -------------------->")
	if num < 0 {
		fmt.Printf("%d 不是回文数\n", num)
		return false
	}
	numStr := fmt.Sprintf("%d", num)
	left, right := 0, len(numStr)-1
	for left < right {
		if numStr[left] != numStr[right] {
			fmt.Printf("%d 不是回文数\n", num)
			return false
		}
		left++
		right--
	}
	fmt.Printf("%d 是回文数\n", num)
	return true
}

var map_str map[string]string = map[string]string{
	"(": ")",
	"[": "]",
	"{": "}",
}

func getValue(key string) (string, bool) {
	for k, v := range map_str {
		if k == key {
			return v, true
		}

	}
	return "", false
}

func matchStr(str_rune []rune) ([]rune, bool) {
	if len(str_rune) < 2 {
		return str_rune, false
	}
	value, ok := getValue(string(str_rune[0]))
	if !ok {
		return str_rune, false
	}

	if value == string(str_rune[1]) {
		return str_rune[2:], true
	} else if value == string(str_rune[len(str_rune)-1]) {
		return str_rune[1 : len(str_rune)-1], true
	}
	return str_rune, false
}

func loopStr(str_rune []rune) ([]rune, bool) {
	rune_str, ok := matchStr(str_rune)
	if ok {
		if len(rune_str) == 0 {
			return rune_str, true
		}
		return loopStr(rune_str)

	}
	return rune_str, false

}

func AnswerQuesion3(str string) {
	// 给定一个只包括 '('，')'，'{'，'}'，'['，']' 的字符串 s ，判断字符串是否有效。

	// 有效字符串需满足：
	// 1.左括号必须用相同类型的右括号闭合。
	// 2.左括号必须以正确的顺序闭合。
	// 3.每个右括号都有一个对应的相同类型的左括号。
	// 示例 1： 输入：s = "()"      输出：true
	// 示例 2： 输入：s = "()[]{}"  输出：true
	// 示例 3： 输入：s = "(]"      输出：false
	// 示例 4： 输入：s = "([])".   输出：true
	// 示例 5： 输入：s = "([)]"   输出：false

	fmt.Println(" <----------------------- AnswerQuesion3 -------------------->")
	rune_str := []rune(str)
	length := len(rune_str)
	if length == 0 || length%2 != 0 {
		fmt.Printf("%s a不是有效字符串\n", str)
		return
	}
	for _, value := range rune_str {
		if value != '(' && value != ')' && value != '[' && value != ']' && value != '{' && value != '}' {
			fmt.Printf("%s b不是有效字符串\n", str)
			return
		}
	}
	_, ok := loopStr(rune_str)
	if ok {
		fmt.Printf("%s 是有效字符串\n", str)
	} else {
		fmt.Printf("%s 不是有效字符串\n", str)
	}
}

func minString(strs []string) string {
	if len(strs) == 0 {
		return ""
	}

	min_str := strs[0]
	for _, v := range strs {
		if len(v) < len(min_str) {
			min_str = v
		}
	}
	return min_str
}

func AnswerQuesion4(x []string) {
	// 编写一个函数来查找字符串数组中的最长公共前缀。
	// 如果不存在公共前缀，返回空字符串 ""。

	// 示例 1： 输入：strs = ["flower","flow","flight"] 输出："fl"
	// 示例 2： 输入：strs = ["dog","racecar","car"]  输出："" 解释：输入不存在公共前缀。
	// 边界条件检查
	fmt.Println(" <----------------------- AnswerQuesion4 -------------------->")
	if len(x) == 0 {
		fmt.Println("最长公共前缀: \"\"")
		return
	}

	if len(x) == 1 {
		fmt.Printf("最长公共前缀: \"%s\"\n", x[0])
		return
	}
	// 找到最短字符串作为比较基准
	min_str := minString(x)
	if len(min_str) == 0 {
		fmt.Println("最长公共前缀: \"\"")
		return
	}
	// 逐个字符比较
	for i := 0; i < len(min_str); i++ {
		char := min_str[i]
		// 检查所有字符串在位置i是否都有相同字符
		for _, str := range x {
			if str[i] != char {
				// 找到不匹配的位置，返回前面的公共前缀
				if i == 0 {
					fmt.Println("最长公共前缀: \"\"")
					return
				}
				fmt.Printf("最长公共前缀: \"%s\"\n", min_str[:i])
				return
			}
		}
	}

	// 所有字符都匹配，最短字符串就是公共前缀
	fmt.Printf("最长公共前缀: \"%s\"\n", min_str)

}

func AnswerQuesion5(nums []int) []int {
	// todo
	// 给定一个表示 大整数 的整数数组 digits，其中 digits[i] 是整数的第 i 位数字。这些数字按从左到右，从最高位到最低位排列。这个大整数不包含任何前导 0。
	// 将大整数加 1，并返回结果的数字数组。

	// 示例 1：
	// 输入：digits = [1,2,3]
	// 输出：[1,2,4]
	// 解释：输入数组表示数字 123。
	// 加 1 后得到 123 + 1 = 124。
	// 因此，结果应该是 [1,2,4]。

	// 示例 2：
	// 输入：digits = [4,3,2,1]
	// 输出：[4,3,2,2]
	// 解释：输入数组表示数字 4321。
	// 加 1 后得到 4321 + 1 = 4322。
	// 因此，结果应该是 [4,3,2,2]。

	// 示例 3：
	// 输入：digits = [9]
	// 输出：[1,0]
	// 解释：输入数组表示数字 9。
	// 加 1 得到了 9 + 1 = 10。
	// 因此，结果应该是 [1,0]。

	// 提示：
	// 1 <= digits.length <= 100
	// 0 <= digits[i] <= 9
	// digits 不包含任何前导 0。
	fmt.Println(" <----------------------- AnswerQuesion5 -------------------->")
	new_nums := []int{}

	if len(nums) == 0 {
		fmt.Println("输入的数组为空")
		return new_nums
	}
	if len(nums) == 1 && nums[0] == 9 {
		new_nums = append(new_nums, 1, 0)
		fmt.Println("加1后的数组:", new_nums)
		return new_nums
	}
	if nums[len(nums)-1] == 9 {
		new_nums = append(new_nums, nums[:len(nums)-2]...)
		new_nums = append(new_nums, nums[len(nums)-2]+1)
		new_nums = append(new_nums, 0)
	} else {
		new_nums = append(new_nums, nums[:len(nums)-1]...)
		new_nums = append(new_nums, nums[len(nums)-1]+1)
	}
	fmt.Println("加1后的数组:", new_nums)
	return new_nums

}

func AnswerQuesion6(nums []int) int {
	// todo
	//给你一个 非严格递增排列 的数组 nums ，请你 原地 删除重复出现的元素，使每个元素 只出现一次 ，返回删除后数组的新长度。
	// 元素的 相对顺序 应该保持 一致 。然后返回 nums 中唯一元素的个数。

	// 示例 1：
	// 输入：nums = [1,1,2]
	// 输出：2, nums = [1,2,_]
	// 解释：函数应该返回新的长度 2 ，并且原数组 nums 的前两个元素被修改为 1, 2 。不需要考虑数组中超出新长度后面的元素。

	// 示例 2：
	// 输入：nums = [0,0,1,1,1,2,2,3,3,4]
	// 输出：5, nums = [0,1,2,3,4]
	// 解释：函数应该返回新的长度 5 ， 并且原数组 nums 的前五个元素被修改为 0, 1, 2, 3, 4 。不需要考虑数组中超出新长度后面的元素。
	fmt.Println(" <----------------------- AnswerQuesion6 -------------------->")
	new_mums := []int{}
	// 这个方案 相对顺序未能保持一致
	// new_map := make(map[int]int)
	// for _, v := range nums {
	// 	new_map[v]++
	// }
	// for k := range new_map {
	// 	new_mums = append(new_mums, k)
	// }
	// 这个方案 相对顺序保持一致
	for _, v := range nums {
		found := false
		for _, v2 := range new_mums {
			if v == v2 {
				found = true
				break
			}
		}
		if !found {
			new_mums = append(new_mums, v)
		}
	}

	fmt.Println("删除重复出现的元素后数组:", new_mums)
	fmt.Println("删除重复出现的元素后数组长度:", len(new_mums))
	return len(new_mums)
}

func AnswerQuesion7(nums []int, target int) []int {
	// todo
	// 给定一个整数数组 nums 和一个整数目标值 target，请你在该数组中找出 和为目标值 target  的那 两个 整数，并返回它们的数组下标。
	// 你可以假设每种输入只会对应一个答案，并且你不能使用两次相同的元素。
	// 你可以按任意顺序返回答案。

	// 	示例 1：
	// 输入：nums = [2,7,11,15], target = 9
	// 输出：[0,1]
	// 解释：因为 nums[0] + nums[1] == 9 ，返回 [0, 1] 。

	// 示例 2：
	// 输入：nums = [3,2,4], target = 6
	// 输出：[1,2]

	// 示例 3：
	// 输入：nums = [3,3], target = 6
	// 输出：[0,1]s
	fmt.Println(" <----------------------- AnswerQuesion7 -------------------->")
	new_nums := []int{}

	for i, v := range nums {
		found_num := target - v
		for j, v2 := range nums {
			if i != j && found_num == v2 {
				fmt.Printf("找到两个数 %d 和 %d 的下标是: [%d, %d]\n", v, v2, i, j)
				new_nums = append(new_nums, i, j)
				return new_nums
			}
		}
	}
	fmt.Println("没有找到")
	return new_nums
}

func main() {
	AnswerQuesion1([]int{2, 2, 3, 1, 3})
	AnswerQuesion2(121)
	AnswerQuesion2(-121)
	AnswerQuesion2(122321)
	AnswerQuesion2(1223221)
	AnswerQuesion3("()")
	AnswerQuesion3("()[]{}")
	AnswerQuesion3("(]")
	AnswerQuesion3("([)]")
	AnswerQuesion3("{[]}")
	AnswerQuesion4([]string{"flower", "flow", "flight"})
	AnswerQuesion4([]string{"dog", "racecar", "car"})
	AnswerQuesion5([]int{1, 2, 3})
	AnswerQuesion5([]int{4, 3, 2, 1})
	AnswerQuesion5([]int{9})
	AnswerQuesion5([]int{6, 2, 9})
	AnswerQuesion6([]int{1, 1, 2})
	AnswerQuesion6([]int{0, 0, 1, 1, 1, 2, 2, 3, 3, 4})
	AnswerQuesion7([]int{2, 7, 11, 15}, 9)
	AnswerQuesion7([]int{3, 2, 4}, 6)
}
