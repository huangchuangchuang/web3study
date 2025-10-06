// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract ReverseString {
    // 定义事件用于调试输出
    event DebugInfo(string message, bytes data);
    event StringInfo(string original, uint256 length);

    /**
     * @dev 反转字符串
     * @param str 输入字符串 abcde
     * @return 反转后的字符串 edcba
     */

    // 定义事件用于调试输出

    function reverseString(string memory str) public returns (string memory) {
        bytes memory strBytes = bytes(str);
        uint256 length = strBytes.length;
        emit StringInfo(str, length);
        emit DebugInfo("Original bytes", strBytes);
        // 创建新的字节数组
        bytes memory reversed = new bytes(length);
        
        // 反转字符
        for (uint256 i = 0; i < length; i++) {
            reversed[i] = strBytes[length - 1 - i];
        }
        emit DebugInfo("reversed bytes", reversed);
        return string(reversed);
    }
}