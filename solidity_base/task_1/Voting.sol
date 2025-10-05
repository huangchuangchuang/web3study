// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;
//  创建一个名为Voting的合约，包含以下功能：
// 一个mapping来存储候选人的得票数 一个vote函数，
// 允许用户投票给某个候选人 
// 一个getVotes函数，返回某个候选人的得票数 
// 一个resetVotes函数，重置所有候选人的得票数
contract Voting {
    // 存储候选人的得票数
    mapping(string => uint256) private votes;
    
    // 存储已投票地址，防止重复投票
    mapping(address => bool) private hasVoted;

    string[] private candidateList;  // 维护键列表
    address[] private voterList;  // 维护投票者列表
    
    // 合约所有者
    address private owner;
    
    // 构造函数，设置合约所有者
    constructor() {
        owner = msg.sender;
    }
    
    /**
     * @dev 允许用户投票给某个候选人
     * @param candidate 候选人姓名
     */
    function vote(string memory candidate) public {
        // 检查用户是否已经投过票
        require(!hasVoted[msg.sender], "You have already voted");
        
        // 检查候选人姓名不能为空
        require(bytes(candidate).length > 0, "Candidate name cannot be empty");
        
        // 记录用户已投票
        hasVoted[msg.sender] = true;
        voterList.push(msg.sender);  // ✅ 添加这行：记录投票者
        // 如果是新候选人，添加到列表
        if (votes[candidate] == 0) {
            candidateList.push(candidate);
        }
        // 增加候选人得票数
        votes[candidate] += 1;
    }
    
    /**
     * @dev 返回某个候选人的得票数
     * @param candidate 候选人姓名
     * @return 候选人的得票数
     */
    function getVotes(string memory candidate) public view returns (uint256) {
        return votes[candidate];
    }
    
    /**
     * @dev 重置所有候选人的得票数（仅合约所有者可以调用）
     */
    function resetVotes() public {
        // 检查调用者是否为合约所有者
        require(msg.sender == owner, "Only owner can reset votes");
        
        // 重置所有投票记录（注意：这只是简单示例，在实际应用中可能需要更复杂的重置机制）
        // 由于 mapping 无法直接清空，这里只是演示基本思路
        for (uint256 i = 0; i < candidateList.length; i++) {
            votes[candidateList[i]] = 0;
        }
        // 重置所有用户的投票状态
        for (uint256 i = 0; i < voterList.length; i++) {
            hasVoted[voterList[i]] = false;
        }
    }
    
    /**
     * @dev 获取合约所有者地址
     * @return 合约所有者地址
     */
    function getOwner() public view returns (address) {
        return owner;
    }
}