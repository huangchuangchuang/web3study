// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC721/IERC721.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import "@openzeppelin/contracts/proxy/utils/UUPSUpgradeable.sol";
import "@chainlink/contracts/src/v0.8/shared/interfaces/AggregatorV3Interface.sol";

contract NFTChainlinkAuctionUUPS is ReentrancyGuard, Ownable, UUPSUpgradeable {
    IERC721 public nft;
    uint256 public nftId;
    address public seller;
    address public highestBidder;
    uint256 public highestBid;
    bool public ended;
    uint256 public endTime;
    uint256 public reservePrice; // 最低出价
    uint256 public feePercentage;  // 手续费百分比
    address public feeRecipient;  // 手续费接收地址
    
    // Chainlink相关变量
    address public ethUsdPriceFeed;  // ETH/USD价格预言机
    address public erc20Token;  // ERC20代币地址
    address public erc20UsdPriceFeed;  // ERC20/USD价格预言机
    
    mapping(address => uint256) public bids;
    
    event Bid(address indexed bidder, uint256 amount, bool isETH);
    event AuctionEnded(address winner, uint256 amount);
    event AuctionCancelled();

    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    function initialize(
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
    ) public initializer {
        __Ownable_init();
        __UUPSUpgradeable_init();
        
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
        
        ethUsdPriceFeed = _ethUsdPriceFeed;
        erc20Token = _erc20Token;
        erc20UsdPriceFeed = _erc20UsdPriceFeed;
    }

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

    function bidERC20(uint256 amount) public nonReentrant {
        require(erc20Token != address(0), "ERC20 token not supported");
        require(erc20UsdPriceFeed != address(0), "ERC20 price feed not available");
        require(block.timestamp < endTime, "Auction has ended");
        require(!ended, "Auction already ended");
        
        require(IERC20(erc20Token).transferFrom(msg.sender, address(this), amount), "ERC20 transfer failed");
        
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

    function convertERC20ToETH(uint256 erc20Amount) public view returns (uint256) {
        if (erc20Token == address(0) || erc20UsdPriceFeed == address(0) || ethUsdPriceFeed == address(0)) {
            return 0;
        }
        
        int256 erc20Price = getChainlinkPrice(erc20UsdPriceFeed);
        int256 ethPrice = getChainlinkPrice(ethUsdPriceFeed);
        
        return (erc20Amount * uint256(erc20Price) * 1e18) / uint256(ethPrice) / 1e8;
    }

    function getChainlinkPrice(address priceFeed) internal view returns (int256) {
        AggregatorV3Interface priceFeedContract = AggregatorV3Interface(priceFeed);
        (, int256 price, , , ) = priceFeedContract.latestRoundData();
        return price;
    }

    function getUSDValueByETH(uint256 ethAmount) public view returns (uint256) {
        if (ethUsdPriceFeed == address(0)) return 0;
        
        int256 ethPrice = getChainlinkPrice(ethUsdPriceFeed);
        return (ethAmount * uint256(ethPrice)) / 1e8 / 1e10;
    }

    function endAuction() public {
        require(block.timestamp >= endTime || highestBid >= reservePrice, "Auction not ended yet");
        require(!ended, "Auction already ended");
        ended = true;

        if (highestBidder != address(0)) {
            uint256 fee = (highestBid * feePercentage) / 100;
            uint256 sellerProceeds = highestBid - fee;
            
            if (fee > 0) {
                (bool feeSuccess, ) = payable(feeRecipient).call{value: fee}("");
                require(feeSuccess, "Failed to send fee");
            }
            
            (bool success, ) = payable(seller).call{value: sellerProceeds}("");
            require(success, "ETH transfer to seller failed");
            
            nft.transferFrom(address(this), highestBidder, nftId);
        } else {
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
        
        if (highestBidder != address(0)) {
            (bool success, ) = payable(highestBidder).call{value: highestBid}("");
            require(success, "Failed to refund bidder");
            highestBidder = address(0);
            highestBid = 0;
        }
        
        nft.transferFrom(address(this), seller, nftId);
        
        emit AuctionCancelled();
    }
    
    function _authorizeUpgrade(address newImplementation) internal override onlyOwner {}
}