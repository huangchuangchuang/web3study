// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/access/Ownable.sol";

contract BeggingContract is Ownable {
    // 记录每个捐赠者的捐赠金额
    mapping(address => uint256) public donations;
    
    // 记录总捐赠金额
    uint256 public totalDonations;
    
    // 捐赠事件
    event Donation(address indexed donor, uint256 amount);
    event Withdrawal(address indexed owner, uint256 amount);

    constructor() Ownable(msg.sender) {}

    /**
     * @dev 允许用户向合约捐赠以太币
     */
    function donate() public payable {
        require(msg.value > 0, "Donation amount must be greater than 0");
        
        donations[msg.sender] += msg.value;
        totalDonations += msg.value;
        
        emit Donation(msg.sender, msg.value);
    }

    /**
     * @dev 允许合约所有者提取所有资金
     */
    function withdraw() public onlyOwner {
        uint256 amount = address(this).balance;
        require(amount > 0, "No funds to withdraw");
        
        payable(owner()).transfer(amount);
        
        emit Withdrawal(owner(), amount);
    }

    /**
     * @dev 查询某个地址的捐赠金额
     * @param donor 捐赠者地址
     * @return 捐赠金额
     */
    function getDonation(address donor) public view returns (uint256) {
        return donations[donor];
    }

    /**
     * @dev 获取合约余额
     */
    function getContractBalance() public view returns (uint256) {
        return address(this).balance;
    }

    // 防止合约接收意外转账
    receive() external payable {
        donate();
    }
}