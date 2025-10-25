// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/token/ERC20/extensions/ERC20Permit.sol";

contract BaseERC20 is ERC20, ERC20Permit {
    constructor() ERC20("MyToken", "MTK") ERC20Permit("MyToken") {
        _mint(msg.sender, 100000 * 10 ** 18);
    }
}

// ERC20Permit 是一个ERC20代币的扩展功能，实现了EIP-2612标准。

// 主要作用：
// 1. 元交易支持
// 允许用户通过签名授权他人代为操作代币
// 无需用户直接发送交易，节省gas费用
// 2. 单次批准和转账
// 可以在一个交易中完成批准和转账操作
// 避免传统的两次交易流程（先approve再transferFrom）
// 3. 签名机制
// 用户对授权信息进行签名，生成permit
// 第三方可以使用该签名执行代币操作