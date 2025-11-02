// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import {AggregatorV3Interface} from "@chainlink/contracts/src/v0.8/shared/interfaces/AggregatorV3Interface.sol";


// mock 价格预言机合约，用于测试;
contract AggregatorV3 is AggregatorV3Interface {
    int256 answer;  // 价格答案
    uint80 public currentRoundId;  // 当前轮次ID
    
    constructor(int256 _answer) {
        answer = _answer;
    }
    
    function decimals() external pure returns (uint8) {
        return 18;
    }

    function description() external pure returns (string memory) {
        return "a test price feeder";
    }

    function version() external pure returns (uint256) {
        return 0;
    }

    function getRoundData(uint80 _roundId) external view returns (uint80, int256, uint256, uint256, uint80) {
        // uint80 roundId,           // 轮次ID
        // int256 answer,            // 价格答案（核心数据）
        // uint256 startedAt,        // 开始时间
        // uint256 updatedAt,        // 更新时间  
        // uint80 answeredInRound    // 回答的轮次

        // 在实际应用中，这里会查询历史数据
        require(_roundId <= currentRoundId && _roundId > 0, "Invalid round ID");
        return (_roundId, answer, 0, 0, 0);
    }

    function latestRoundData() public view returns (uint80, int256, uint256, uint256, uint80) {
        return (currentRoundId, answer, 0, 0, 0);
    }

    function setPrice(int256 _answer) public {
        answer = _answer;
        currentRoundId++;  // 每次更新价格，轮次ID增加
    }
}