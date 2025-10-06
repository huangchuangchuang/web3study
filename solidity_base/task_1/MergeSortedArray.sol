// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract MergeSortedArray {
    /**
     * @dev 合并两个有序数组
     * @param nums1 第一个有序数组
     * @param nums2 第二个有序数组
     * @return 合并后的有序数组
     */
    function mergeSortedArrays(int256[] memory nums1, int256[] memory nums2) 
        public pure returns (int256[] memory) 
    {
        uint256 m = nums1.length;
        uint256 n = nums2.length;
        int256[] memory result = new int256[](m + n);
        
        uint256 i = 0; // nums1 的索引
        uint256 j = 0; // nums2 的索引
        uint256 k = 0; // result 的索引
        
        // 合并两个数组，下面i，j其中一组循环完成终止
        while (i < m && j < n) {
            if (nums1[i] <= nums2[j]) {
                result[k] = nums1[i];
                i++;
            } else {
                result[k] = nums2[j];
                j++;
            }
            k++;
        }
        
        // 另一组i或j为剩余元素，合并追加
        while (i < m) {
            result[k] = nums1[i];
            i++;
            k++;
        }
        
        while (j < n) {
            result[k] = nums2[j];
            j++;
            k++;
        }
        
        return result;
    }
}