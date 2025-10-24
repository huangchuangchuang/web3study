
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC721/IERC721.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import "@chainlink/contracts/src/v0.8/shared/interfaces/AggregatorV3Interface.sol";


contract NFTChainlinkAuction is ReentrancyGuard, Ownable {
    IERC721 public nft;
    uint256 public nftId;
    address public seller;
    address public highestBidder;
    uint256 public highestBid;
    bool public ended;
    uint256 public endTime;
    uint256 public reservePrice;
    uint256 public feePercentage;
    address public feeRecipient;
    
    // Chainlink相关变量
    address public ethUsdPriceFeed;  // ETH/USD价格预言机
    address public erc20Token;       // ERC20代币地址
    address public erc20UsdPriceFeed; // ERC20/USD价格预言机
    
    mapping(address => uint256) public bids;
    
    event Bid(address indexed bidder, uint256 amount, bool isETH);
    event AuctionEnded(address winner, uint256 amount);
    event AuctionCancelled();

    constructor(
        address _nft,
        uint256 _nftId,
        address _seller,
        uint256 _duration,
        uint256 _reservePrice,
        uint256 _feePercentage,
        address _feeRecipient,
        address _ethUsdPriceFeed,
        address _erc20Token,
        address _erc20UsdPriceFeed
    ) Ownable(msg.sender) {
        require(_nft != address(0), "Invalid NFT address");
        require(_seller != address(0), "Invalid seller address");
        require(_duration > 0, "Duration must be positive");
        require(_feePercentage <= 100, "Fee percentage too high");
        require(_feeRecipient != address(0), "Invalid fee recipient");
        
        nft = IERC721(_nft);
        nftId = _nftId;
        seller = _seller;
        reservePrice = _reservePrice;
        feePercentage = _feePercentage;
        feeRecipient = _feeRecipient;
        endTime = block.timestamp + _duration;
        
        // Chainlink预言机地址
        ethUsdPriceFeed = _ethUsdPriceFeed;
        erc20Token = _erc20Token;
        erc20UsdPriceFeed = _erc20UsdPriceFeed;
    }

    // ETH出价
    function bidETH() public payable nonReentrant {
        require(block.timestamp < endTime, "Auction has ended");
        require(!ended, "Auction already ended");
        require(msg.value >= reservePrice, "Bid below reserve price");
        
        uint256 usdValue = getUSDValueByETH(msg.value);
        uint256 currentHighestUSD = highestBid > 0 ? getUSDValueByETH(highestBid) : 0;
        require(usdValue > currentHighestUSD, "Bid not high enough in USD");

        if (highestBidder != address(0)) {
            bids[highestBidder] += highestBid;
        }

        highestBidder = msg.sender;
        highestBid = msg.value;

        emit Bid(msg.sender, msg.value, true);
    }

    // ERC20出价
    function bidERC20(uint256 amount) public nonReentrant {
        require(erc20Token != address(0), "ERC20 token not supported");
        require(erc20UsdPriceFeed != address(0), "ERC20 price feed not available");
        require(block.timestamp < endTime, "Auction has ended");
        require(!ended, "Auction already ended");
        
        // 转移ERC20代币到合约
        require(IERC20(erc20Token).transferFrom(msg.sender, address(this), amount), "ERC20 transfer failed");
        
        // 计算等值ETH金额
        uint256 ethValue = convertERC20ToETH(amount);
        require(ethValue >= reservePrice, "Bid below reserve price");
        
        uint256 usdValue = getUSDValueByETH(ethValue);
        uint256 currentHighestUSD = highestBid > 0 ? getUSDValueByETH(highestBid) : 0;
        require(usdValue > currentHighestUSD, "Bid not high enough in USD");

        if (highestBidder != address(0)) {
            bids[highestBidder] += highestBid;
        }

        highestBidder = msg.sender;
        highestBid = ethValue;

        emit Bid(msg.sender, amount, false);
    }

    // 将ERC20代币金额转换为ETH金额
    function convertERC20ToETH(uint256 erc20Amount) public view returns (uint256) {
        if (erc20Token == address(0) || erc20UsdPriceFeed == address(0) || ethUsdPriceFeed == address(0)) {
            return 0;
        }
        
        int256 erc20Price = getChainlinkPrice(erc20UsdPriceFeed);
        int256 ethPrice = getChainlinkPrice(ethUsdPriceFeed);
        
        // (ERC20数量 * ERC20价格) / ETH价格 = 等值ETH数量
        return (erc20Amount * uint256(erc20Price) * 1e18) / uint256(ethPrice) / 1e8;
    }

    // 获取Chainlink价格数据
    function getChainlinkPrice(address priceFeed) internal view returns (int256) {
        AggregatorV3Interface priceFeedContract = AggregatorV3Interface(priceFeed);
        (, int256 price, , , ) = priceFeedContract.latestRoundData();
        return price;
    }

    // 获取金额的USD价值
    function getUSDValueByETH(uint256 ethAmount) public view returns (uint256) {
        if (ethUsdPriceFeed == address(0)) return 0;
        
        int256 ethPrice = getChainlinkPrice(ethUsdPriceFeed);
        // ethAmount * ethPrice / 1e8 (因为Chainlink价格有8位小数)
        return (ethAmount * uint256(ethPrice)) / 1e8 / 1e10; // 调整为正确的USD值
    }

    function endAuction() public {
        require(block.timestamp >= endTime || highestBid >= reservePrice, "Auction not ended yet");
        require(!ended, "Auction already ended");
        ended = true;

        if (highestBidder != address(0)) {
            // 计算手续费
            uint256 fee = (highestBid * feePercentage) / 100;
            uint256 sellerProceeds = highestBid - fee;
            
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
    
    function getAuctionInfo() public view returns (
        address _nft,
        uint256 _nftId,
        address _seller,
        address _highestBidder,
        uint256 _highestBid,
        uint256 _endTime,
        uint256 _reservePrice,
        bool _ended,
        uint256 _ethUsdPrice,
        uint256 _erc20UsdPrice
    ) {
        uint256 ethUsdPrice = ethUsdPriceFeed != address(0) ? uint256(getChainlinkPrice(ethUsdPriceFeed)) : 0;
        uint256 erc20UsdPrice = erc20UsdPriceFeed != address(0) ? uint256(getChainlinkPrice(erc20UsdPriceFeed)) : 0;
        
        return (
            address(nft),
            nftId,
            seller,
            highestBidder,
            highestBid,
            endTime,
            reservePrice,
            ended,
            ethUsdPrice,
            erc20UsdPrice
        );
    }
}