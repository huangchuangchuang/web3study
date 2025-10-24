// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "./NFTChainlinkAuction.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import "@openzeppelin/contracts/token/ERC721/IERC721.sol";

contract NFTChainlinkAuctionFactory is Ownable, ReentrancyGuard {
    // 存储所有创建的拍卖合约地址
    address[] public allAuctions;
    
    // 映射：NFT合约地址 => NFT ID => 拍卖合约地址
    mapping(address => mapping(uint256 => address)) public getAuction;
    
    // 映射：拍卖合约地址 => 是否为有效拍卖
    mapping(address => bool) public isAuction;
    
    // 手续费相关
    uint256 public feePercentage = 2;
    address public feeRecipient;
    
    // Chainlink预言机地址（可配置）
    address public ethUsdPriceFeed;
    address public defaultERC20Token;
    address public defaultERC20UsdPriceFeed;
    
    // 事件
    event AuctionCreated(
        address indexed auction,
        address indexed nftContract,
        uint256 indexed nftId,
        address seller,
        uint256 duration,
        uint256 fee
    );
    
    event FeePercentageUpdated(uint256 oldFee, uint256 newFee);
    event FeeRecipientUpdated(address oldRecipient, address newRecipient);
    event ChainlinkFeedsUpdated(address ethUsdFeed, address erc20Token, address erc20UsdFeed);

    constructor(
        address _feeRecipient,
        address _ethUsdPriceFeed,
        address _defaultERC20Token,
        address _defaultERC20UsdPriceFeed
    ) Ownable(msg.sender) {
        feeRecipient = _feeRecipient;
        ethUsdPriceFeed = _ethUsdPriceFeed;
        defaultERC20Token = _defaultERC20Token;
        defaultERC20UsdPriceFeed = _defaultERC20UsdPriceFeed;
    }

    function createAuction(
        address _nft,
        uint256 _nftId,
        uint256 _duration,
        uint256 _reservePrice
    ) external nonReentrant returns (address) {
        require(_duration > 0, "Duration must be positive");
        require(_nft != address(0), "Invalid NFT address");
        require(_duration <= 365 days, "Duration too long");
        
        address existingAuction = getAuction[_nft][_nftId];
        require(existingAuction == address(0) || !isAuction[existingAuction], 
                "Auction already exists for this NFT");
        
        // 创建新的拍卖合约（包含Chainlink参数）
        NFTAuction auction = new NFTAuction(
            _nft,
            _nftId,
            msg.sender,
            _duration,
            _reservePrice,
            feePercentage,
            feeRecipient,
            ethUsdPriceFeed,              // ETH/USD价格预言机
            defaultERC20Token,            // 默认ERC20代币
            defaultERC20UsdPriceFeed      // ERC20/USD价格预言机
        );
        
        address auctionAddress = address(auction);
        
        allAuctions.push(auctionAddress);
        getAuction[_nft][_nftId] = auctionAddress;
        isAuction[auctionAddress] = true;
        
        // 将NFT转移到拍卖合约
        IERC721(_nft).transferFrom(msg.sender, auctionAddress, _nftId);
        
        emit AuctionCreated(
            auctionAddress,
            _nft,
            _nftId,
            msg.sender,
            _duration,
            feePercentage
        );
        
        return auctionAddress;
    }

    // 其他函数保持不变...

    // 更新Chainlink预言机地址
    function updateChainlinkFeeds(
        address _ethUsdPriceFeed,
        address _defaultERC20Token,
        address _defaultERC20UsdPriceFeed
    ) external onlyOwner {
        ethUsdPriceFeed = _ethUsdPriceFeed;
        defaultERC20Token = _defaultERC20Token;
        defaultERC20UsdPriceFeed = _defaultERC20UsdPriceFeed;
        emit ChainlinkFeedsUpdated(_ethUsdPriceFeed, _defaultERC20Token, _defaultERC20UsdPriceFeed);
    }
}