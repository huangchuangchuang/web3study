// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

/**
 * @title Counter
 * @dev 一个简单的计数器合约，支持增加、减少和重置计数
 */
contract Counter {
    // 状态变量，存储当前计数值
    uint256 private count;

    // 事件，当计数改变时触发
    event CountChanged(uint256 newCount);

    // 事件，当计数重置时触发
    event CountReset();

    /**
     * @dev 构造函数，初始化计数器为0
     */
    constructor() {
        count = 0;
    }

    /**
     * @dev 增加计数
     * @param _value 要增加的值
     */
    function increment(uint256 _value) public {
        count += _value;
        emit CountChanged(count);
    }

    /**
     * @dev 减少计数
     * @param _value 要减少的值
     */
    function decrement(uint256 _value) public {
        require(count >= _value, "Counter: decrement overflow");
        count -= _value;
        emit CountChanged(count);
    }

    /**
     * @dev 重置计数为0
     */
    function reset() public {
        count = 0;
        emit CountReset();
    }

    /**
     * @dev 获取当前计数值
     * @return 当前计数值
     */
    function getCount() public view returns (uint256) {
        return count;
    }
}