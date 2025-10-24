// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/token/ERC721/IERC721.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";


contract NFTAutionSingle is ReentrancyGuard, Ownable {
    IERC721 public nft;
    uint256 public nftId;
    address public seller;
    address public highestBidder;
    uint256 public highestBid;
    bool public ended;

    mapping(address => uint256) public bids;

    event Bid(address indexed bidder, uint256 amount);
    event AuctionEnded(address winner, uint256 amount);

    constructor(address _nft, uint256 _nftId, address _seller) Ownable(msg.sender){
        nft = IERC721(_nft);
        nftId = _nftId;
        seller = _seller;
    }

    function bid() public payable nonReentrant {
        // 买家的钱进入合约账户
        require(!ended, "Auction ended");
        require(msg.value > highestBid, "Bid not high enough");

        if (highestBidder != address(0)) {
            bids[highestBidder] = highestBid; // 记录之前最高出价者的出价，方便后续提现
        }

        highestBidder = msg.sender; // 记录当前出价者
        highestBid = msg.value; // 记录当前出价金额

        emit Bid(msg.sender, msg.value);
    }

    function endAuction() public onlyOwner{
        // 合约将资金转给卖家，NFT转给最高出价者
        require(!ended, "Auction ended");
        ended = true;

        if (highestBidder != address(0)) {
            // // 先转移ETH
            // payable(seller).transfer(highestBid);

            // 先转移ETH，使用call方式处理可能的失败
            (bool success, ) = payable(seller).call{value: highestBid}("");
            require(success, "ETH transfer to seller failed");
            // 再转移NFT
            nft.transferFrom(address(this), highestBidder, nftId);
        } else {
            nft.transferFrom(address(this), seller, nftId);
        }

        emit AuctionEnded(highestBidder, highestBid);
    }

    function withdraw() public nonReentrant {
        // 退款阶段（合约 → 其他买家）

        // 合约将资金转回调用者（之前出价但未获胜的买家）
        uint256 amount = bids[msg.sender];
        require(amount > 0, "No funds to withdraw");
        bids[msg.sender] = 0;
        payable(msg.sender).transfer(amount);
    }

    
}
