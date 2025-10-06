// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract BinarySearch {
    /**
     * @dev 在有序数组中二分查找目标值
     * @param nums 有序数组
     * @param target 目标值
     * @return 目标值的索引，如果不存在返回-1
     */
    function binarySearch(int256[] memory nums, int256 target) 
        public pure returns (int256) 
    {
        uint256 left = 0;
        uint256 right = nums.length;
        
        while (left < right) {
            uint256 mid = left + (right - left) / 2;
            
            if (nums[mid] == target) {
                return int256(mid);
            } else if (nums[mid] < target) {
                left = mid + 1;
            } else {
                right = mid;
            }
        }
        
        return -1; // 未找到
    }
    
    /**
     * @dev 在有序数组中查找目标值（递归版本）
     * @param nums 有序数组
     * @param target 目标值
     * @return 目标值的索引，如果不存在返回-1
     */
    function binarySearchRecursive(int256[] memory nums, int256 target) 
        public pure returns (int256) 
    {
        return binarySearchHelper(nums, target, 0, int256(nums.length) - 1);
    }
    
    /**
     * @dev 二分查找辅助函数
     * @param nums 有序数组
     * @param target 目标值
     * @param left 左边界
     * @param right 右边界
     * @return 目标值的索引，如果不存在返回-1
     */
    function binarySearchHelper(
        int256[] memory nums, 
        int256 target, 
        int256 left, 
        int256 right
    ) internal pure returns (int256) {
        if (left > right) {
            return -1;
        }
        
        int256 mid = left + (right - left) / 2;
        
        if (nums[uint256(mid)] == target) {
            return mid;
        } else if (nums[uint256(mid)] < target) {
            return binarySearchHelper(nums, target, mid + 1, right);
        } else {
            return binarySearchHelper(nums, target, left, mid - 1);
        }
    }
}