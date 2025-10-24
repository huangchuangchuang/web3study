// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@chainlink/contracts-ccip/src/v0.8/ccip/interfaces/IRouterClient.sol";
import "@chainlink/contracts-ccip/src/v0.8/ccip/libraries/Client.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

contract CrossChainAuction is Ownable {
    // CCIP相关变量
    IRouterClient public immutable ccipRouter;  // CCIP路由合约地址
    uint64 public destinationChainSelector;  // 目标链选择器
    address public crossChainAuctionContract;  // 目标链上的拍卖合约地址
    
    // 拍卖状态
    uint256 public auctionId;  // 拍卖ID
    address public nftContract;  // NFT合约地址
    uint256 public nftId;
    uint256 public highestBid;  // 最高出价
    address public highestBidder;  // 最高出价者
    bool public auctionEnded;  // 拍卖是否结束
    
    // 出价记录
    mapping(address => uint256) public crossChainBids;  // 跨链出价记录
    
    event CrossChainBid(address bidder, uint256 amount, uint64 sourceChain);  // 跨链出价事件
    event CrossChainAuctionEnded(address winner, uint256 amount);  // 跨链拍卖结束事件
    
    constructor(
        address _ccipRouter,
        uint64 _destinationChainSelector,
        address _crossChainAuctionContract,
        uint256 _auctionId,
        address _nftContract,
        uint256 _nftId
    ) Ownable(msg.sender) {
        ccipRouter = IRouterClient(_ccipRouter);  // CCIP路由合约地址
        destinationChainSelector = _destinationChainSelector;  // 目标链选择器
        crossChainAuctionContract = _crossChainAuctionContract;  // 目标链上的拍卖合约地址
        auctionId = _auctionId;  // 拍卖ID
        nftContract = _nftContract;   // NFT合约地址
        nftId = _nftId;
    }
    
    // 接收跨链出价
    function receiveCrossChainBid(address bidder, uint256 amount) external {
        require(msg.sender == address(ccipRouter), "Only CCIP router can call this");
        require(!auctionEnded, "Auction ended");
        require(amount > highestBid, "Bid not high enough");
        
        if (highestBidder != address(0)) {
            crossChainBids[highestBidder] += highestBid;
        }
        
        highestBidder = bidder;
        highestBid = amount;
        
        emit CrossChainBid(bidder, amount, destinationChainSelector);
    }
    
    // 发送跨链出价到其他链
    function sendCrossChainBid(uint64 targetChainSelector, address targetContract, uint256 amount) external payable {
        Client.EVM2AnyMessage memory message = Client.EVM2AnyMessage({
            receiver: abi.encode(targetContract),
            data: abi.encodeWithSelector(
                CrossChainAuction.receiveCrossChainBid.selector,
                msg.sender,
                amount
            ),
            tokenAmounts: new Client.EVMTokenAmount[](0),
            feeToken: address(0),
            extraArgs: ""
        });  // 构造跨链消息
        
        uint256 fee = ccipRouter.getFee(targetChainSelector, message);  // 计算CCIP费用
        require(msg.value >= fee, "Insufficient CCIP fee");
        
        ccipRouter.ccipSend{value: msg.value}(targetChainSelector, message);  // 发送跨链消息
    }
    
    // 结束跨链拍卖
    function endCrossChainAuction() external onlyOwner {
        require(!auctionEnded, "Auction already ended");
        auctionEnded = true;
        
        emit CrossChainAuctionEnded(highestBidder, highestBid);
        
        // 这里可以添加将NFT转移给获胜者的逻辑
        // 或者发送消息到其他链通知拍卖结束
    }
    
    // 提取未获胜的出价
    function withdraw() external {
        uint256 amount = crossChainBids[msg.sender];  // 获取出价金额
        require(amount > 0, "No funds to withdraw");
        crossChainBids[msg.sender] = 0;  // 重置出价记录
        
        (bool success, ) = payable(msg.sender).call{value: amount}("");
        require(success, "Transfer failed");
    }
    
    // CCIP费用接收函数
    receive() external payable {}
}