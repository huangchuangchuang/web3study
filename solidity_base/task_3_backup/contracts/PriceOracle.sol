// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "@chainlink/contracts/src/v0.8/shared/interfaces/AggregatorV3Interface.sol";

contract PriceOracle {
    mapping(address => address) public priceFeeds; // token => feed address
    
    constructor() {}
    
    function setPriceFeed(address token, address feed) external {
        priceFeeds[token] = feed;
    }
    
    function getPrice(address token) public view returns (int, uint8) {
        address feedAddress = priceFeeds[token];
        require(feedAddress != address(0), "Price feed not found");
        
        AggregatorV3Interface priceFeed = AggregatorV3Interface(feedAddress);
        (, int price, , , ) = priceFeed.latestRoundData();
        uint8 decimals = priceFeed.decimals();
        
        return (price, decimals);
    }
    
    function convertToUSD(uint256 amount, address token) public view returns (uint256) {
        (int price, uint8 decimals) = getPrice(token);
        return (amount * uint256(price)) / (10 ** decimals);
    }
}