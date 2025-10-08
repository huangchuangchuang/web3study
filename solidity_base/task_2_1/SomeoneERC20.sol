// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;
//任务：参考 openzeppelin-contracts/contracts/token/ERC20/IERC20.sol实现一个简单的 ERC20 代币合约。要求：
//合约包含以下标准 ERC20 功能：
//balanceOf：查询账户余额。
//transfer：转账。
//approve 和 transferFrom：授权和代扣转账。
//使用 event 记录转账和授权操作。
//提供 mint 函数，允许合约所有者增发代币。
//提示：
//使用 mapping 存储账户余额和授权信息。
//使用 event 定义 Transfer 和 Approval 事件。
//部署到 sepolia 测试网，导入到自己的钱包
contract SomeoneERC20 {
    // 代币名称
    string public name;
    // 代币符号
    string public symbol;
    // 小数位数，通常为18
    uint8 public decimals = 18;
    // 总供应量
    uint256 public totalSupply;

    // 存储账户余额
    mapping(address => uint256) public balanceOf;
    // 存储授权信息: allowance[owner][spender] = amount
    mapping(address => mapping(address => uint256)) public allowance;

    // 合约所有者
    address public owner;

    // 转账事件
    event Transfer(address indexed from, address indexed to, uint256 value);
    // 授权事件
    event Approval(address indexed owner, address indexed spender, uint256 value);

    // 构造函数，初始化代币信息
    constructor(string memory _name, string memory _symbol, uint256 _initialSupply) {
        name = _name;
        symbol = _symbol;
        owner = msg.sender;
        // 考虑小数位数，初始供应量乘以10^decimals
        totalSupply = _initialSupply * (10 **uint256(decimals));
        // 将初始供应量分配给合约部署者
        balanceOf[msg.sender] = totalSupply;
        emit Transfer(address(0), msg.sender, totalSupply);
    }

    // 限制只有所有者可以调用
    modifier onlyOwner() {
        require(msg.sender == owner, "Only owner can call this function");
        _;
    }

    // 转账函数
    function transfer(address to, uint256 amount) public returns (bool success) {
        require(balanceOf[msg.sender] >= amount, "Insufficient balance");
        require(to != address(0), "Cannot transfer to the zero address");

        balanceOf[msg.sender] -= amount;
        balanceOf[to] += amount;

        emit Transfer(msg.sender, to, amount);
        return true;
    }
    //直接用 transfer 只能 “自己转自己的钱”，而 approve + transferFrom 实现了 “让别人代转自己的钱”，这在去中心化场景中非常重要：
    // 授权函数  授权 spender（被授权方）从自己（调用者）的账户中最多转出 amount 数量的代币。
    function approve(address spender, uint256 amount) public returns (bool success) {
        require(spender != address(0), "Cannot approve the zero address");

        allowance[msg.sender][spender] = amount;

        emit Approval(msg.sender, spender, amount);
        return true;
    }

    // 授权转账函数 让 spender（调用者）从 from（授权者）的账户中，向 to（接收者）转出 amount 数量的代币。
    function transferFrom(address from, address to, uint256 amount) public returns (bool success) {
        require(balanceOf[from] >= amount, "Insufficient balance");
        require(allowance[from][msg.sender] >= amount, "Allowance exceeded");
        require(to != address(0), "Cannot transfer to the zero address");

        balanceOf[from] -= amount;
        balanceOf[to] += amount;
        allowance[from][msg.sender] -= amount;

        emit Transfer(from, to, amount);
        return true;
    }

    // 增发代币函数，只有所有者可以调用
    function mint(address to, uint256 amount) public onlyOwner returns (bool success) {
        require(to != address(0), "Cannot mint to the zero address");

        totalSupply += amount;
        balanceOf[to] += amount;

        emit Transfer(address(0), to, amount);
        return true;
    }
}