// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC721/IERC721.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";


contract NFTAuction is ReentrancyGuard, Ownable {
    IERC721 public nft;
    uint256 public nftId;
    address public seller;
    address public highestBidder;
    uint256 public highestBid;
    bool public ended;
    uint256 public endTime;
    uint256 public reservePrice;  // 保留价格
    uint256 public feePercentage;  // 手续费百分比
    address public feeRecipient;  // 手续费接收者地址
    
    mapping(address => uint256) public bids;
    
    event Bid(address indexed bidder, uint256 amount);  // 出价事件
    event AuctionEnded(address winner, uint256 amount);  // 拍卖结束事件
    event AuctionCancelled();  // 拍卖取消事件

    constructor(
        address _nft,
        uint256 _nftId,
        address _seller,
        uint256 _duration,
        uint256 _reservePrice,
        uint256 _feePercentage,
        address _feeRecipient
    ) Ownable(msg.sender){
        require(_nft != address(0), "Invalid NFT address");
        require(_seller != address(0), "Invalid seller address");
        require(_duration > 0, "Duration must be positive");
        require(_feePercentage <= 100, "Fee percentage too high");
        require(_feeRecipient != address(0), "Invalid fee recipient");
        
        nft = IERC721(_nft);   // 设置NFT合约地址
        nftId = _nftId;  // 设置NFT ID
        seller = _seller;  // 设置卖家地址
        reservePrice = _reservePrice;  // 设置保留价格
        feePercentage = _feePercentage;  // 设置手续费百分比
        feeRecipient = _feeRecipient;  // 设置手续费接收者
        endTime = block.timestamp + _duration;
    }

    function bid() public payable nonReentrant {
        require(block.timestamp < endTime, "Auction has ended");
        require(!ended, "Auction already ended");
        require(msg.value >= reservePrice, "Bid below reserve price");
        require(msg.value > highestBid, "Bid not high enough");

        if (highestBidder != address(0)) {
            bids[highestBidder] += highestBid;
        }

        highestBidder = msg.sender;
        highestBid = msg.value;

        emit Bid(msg.sender, msg.value);
    }

    function endAuction() public {
        require(block.timestamp >= endTime || highestBid >= reservePrice, "Auction not ended yet");
        require(!ended, "Auction already ended");
        ended = true;

        if (highestBidder != address(0)) {
            // 计算手续费
            uint256 fee = (highestBid * feePercentage) / 100;  // 手续费金额
            uint256 sellerProceeds = highestBid - fee;  // 卖家所得
            
            // 转移手续费给平台
            if (fee > 0) {
                (bool feeSuccess, ) = payable(feeRecipient).call{value: fee}("");
                require(feeSuccess, "Failed to send fee");
            }
            
            // 转移剩余资金给卖家
            (bool success, ) = payable(seller).call{value: sellerProceeds}("");
            require(success, "ETH transfer to seller failed");
            
            // 转移NFT给获胜者
            nft.transferFrom(address(this), highestBidder, nftId);
        } else {
            // 如果没有出价，NFT退还给卖家
            nft.transferFrom(address(this), seller, nftId);
        }

        emit AuctionEnded(highestBidder, highestBid);
    }

    function withdraw() public nonReentrant {
        uint256 amount = bids[msg.sender];
        require(amount > 0, "No funds to withdraw");
        bids[msg.sender] = 0;
        
        (bool success, ) = payable(msg.sender).call{value: amount}("");
        require(success, "ETH transfer failed");
    }
    
    // 紧急取消
    function emergencyCancel() public onlyOwner {
        require(!ended, "Auction already ended");
        ended = true;
        
        // 退还当前最高出价
        if (highestBidder != address(0)) {
            (bool success, ) = payable(highestBidder).call{value: highestBid}("");
            require(success, "Failed to refund bidder");
            highestBidder = address(0);
            highestBid = 0;
        }
        
        // NFT退还给卖家
        nft.transferFrom(address(this), seller, nftId);
        
        emit AuctionCancelled();
    }
    
    // 获取拍卖状态的辅助函数
    function getAuctionInfo() public view returns (
        address _nft,
        uint256 _nftId,
        address _seller,
        address _highestBidder,
        uint256 _highestBid,
        uint256 _endTime,
        uint256 _reservePrice,
        bool _ended
    ) {
        return (
            address(nft),
            nftId,
            seller,
            highestBidder,
            highestBid,
            endTime,
            reservePrice,
            ended
        );
    }
}