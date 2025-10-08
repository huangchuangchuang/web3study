// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

// ✅ 作业 1：ERC20 代币
// 任务：参考 openzeppelin-contracts/contracts/token/ERC20/IERC20.sol实现一个简单的 ERC20 代币合约。要求：
// 合约包含以下标准 ERC20 功能：
// balanceOf：查询账户余额。
// transfer：转账。
// approve 和 transferFrom：授权和代扣转账。
// 使用 event 记录转账和授权操作。
// 提供 mint 函数，允许合约所有者增发代币。
// 提示：
// 使用 mapping 存储账户余额和授权信息。
// 使用 event 定义 Transfer 和 Approval 事件。
// 部署到sepolia 测试网，导入到自己的钱包

contract MyERC20 {
    // 小数位数的倍数。使用 constant 恒定（无 gas 消耗）
    uint256 private constant DECIMALS_MULTIPLIER = 10 ** 18;
    // 精度，最小位数; 用public测试网测试识别最小位，不然要自己手动输入
    uint8 public decimals = 18;
    // 合约部署者
    address public owner;
    // 代币名称
    string public name;
    // 代币符号
    string public symbol;
    // 总供应量
    uint256 public totalSupply;
    // 存储账户余额
    mapping(address => uint256) public balanceOf;
    // 存储授权信息:
    mapping(address from => mapping(address => uint256) to) public allowance;

    // 转帐事件
    event Tranfer(address indexed from, address indexed to, uint256 amount);
    // 授权事件
    event Approve(
        address indexed owner,
        address indexed spender,
        uint256 value
    );

    //
    constructor(
        string memory _name,
        string memory _symbol,
        uint _initialSupply
    ) {
        name = _name;
        symbol = _symbol;
        owner = msg.sender;
        // 初始化代币数量
        totalSupply = _initialSupply * DECIMALS_MULTIPLIER;
        // 将初始化代币分配给部署者
        balanceOf[msg.sender] = totalSupply;
        emit Tranfer(address(0), msg.sender, totalSupply);
    }

    modifier onlyOwner() {
        require(msg.sender == owner, "Only owner can call this function");
        _;
    }

    // 转账函数
    function transfer(address to, uint256 amount) public returns (bool) {
        require(to != address(0), "transfer to the zero address");
        require(
            balanceOf[msg.sender] >= amount,
            "transfer amount exceeds balance"
        );
        balanceOf[msg.sender] -= amount;
        balanceOf[to] += amount;
        emit Tranfer(msg.sender, to, amount);
        return true;
    }

    // 授权函数, spender 被授权者
    function approve(address spender, uint256 amount) public returns (bool) {
        require(spender != address(0), "Cannot approve the zero address");
        allowance[msg.sender][spender] = amount;
        emit Approve(msg.sender, spender, amount);
        return true;
    }

    // 授权转账函数
    function transferFrom(address from,address to,uint256 amount) public returns (bool) {
        require(from != address(0), "Cannot transfer the zero address");
        require(to != address(0), "Cannot transfer the zero address");
        require(allowance[from][msg.sender] >= amount, "The transfer amount exceeds allowance");
        require(balanceOf[from] >= amount, "The transfer amount exceeds balance");
        balanceOf[from] -= amount;
        balanceOf[to] += amount;
        allowance[from][msg.sender] -= amount;

        emit Tranfer(from, to, amount);
        return true;
    }

    // 增发代币函数，只有所有者可以调用
    function mint(uint256 amount) public onlyOwner returns (bool){
        totalSupply += amount;
        balanceOf[owner] += amount;
        
        emit Tranfer(address(0), owner, amount);
        return true;
    }

    // 便捷显示函数 - 货币供应总量
    function totalSupplyDisplay() public view returns (uint256) {
        return totalSupply / DECIMALS_MULTIPLIER;
    }
    // 便捷显示函数 - 余额
    function balanceOfDisplay() public view returns (uint256) {
        return balanceOf[msg.sender] / DECIMALS_MULTIPLIER;
    }
    // 便捷函数 - 转帐
    function transferTokens(
        address to,
        uint256 tokensAmount
    ) public returns (bool success) {
        uint256 weiAmount = tokensAmount * DECIMALS_MULTIPLIER;
        return transfer(to, weiAmount);
    }
}
