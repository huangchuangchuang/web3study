// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "./NFTChainlinkAuctionUUPS.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import "@openzeppelin/contracts/token/ERC721/IERC721.sol";
import "@openzeppelin/contracts/proxy/utils/UUPSUpgradeable.sol";
import "@openzeppelin/contracts/proxy/ERC1967/ERC1967Proxy.sol";


// UUPS 升级的是逻辑合约，而不是代理合约

// 代理合约 (Proxy)
// 地址固定：部署后地址永远不会改变
// 存储数据：保存所有状态变量
// 转发调用：将函数调用转发给逻辑合约
// 升级功能：包含指向逻辑合约的指针

// 逻辑合约 (Implementation)
// 可替换：可以被新版本替换
// 执行逻辑：包含实际的业务逻辑
// 升级控制：包含 _authorizeUpgrade 函数
contract NFTChainlinkAuctionFactoryUUPS is Ownable, ReentrancyGuard, UUPSUpgradeable {
    address public auctionImplementation;  // 逻辑合约地址
    address[] public allAuctions;  // 存储所有创建的拍卖合约地址
    
    mapping(address => mapping(uint256 => address)) public getAuction;  // NFT合约地址 => NFT ID => 拍卖合约地址
    mapping(address => bool) public isAuction;  // 拍卖合约地址 => 是否为有效拍卖
    
    uint256 public feePercentage = 2;  // 手续费百分比
    address public feeRecipient;  // 手续费接收地址
    
    address public ethUsdPriceFeed;  // ETH/USD价格预言机地址
    address public defaultERC20Token;  // ERC20代币地址
    address public defaultERC20UsdPriceFeed;  // ERC20/USD价格预言机地址
    
    event AuctionCreated(
        address indexed auction,
        address indexed nftContract,
        uint256 indexed nftId,
        address seller,
        uint256 duration,
        uint256 fee
    );
    
    event FeePercentageUpdated(uint256 oldFee, uint256 newFee);  // 手续费百分比更新事件
    event FeeRecipientUpdated(address oldRecipient, address newRecipient);  // 手续费接收地址更新事件
    event ChainlinkFeedsUpdated(address ethUsdFeed, address erc20Token, address erc20UsdFeed);  // Chainlink预言机地址更新事件
    event ImplementationUpdated(address indexed implementation);  // 逻辑合约地址更新事件

    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    function initialize(
        address _feeRecipient,
        address _ethUsdPriceFeed,
        address _defaultERC20Token,
        address _defaultERC20UsdPriceFeed,
        address _auctionImplementation
    ) public initializer {
        // 使用 initializer 修饰符确保这些设置只能执行一次。
        __Ownable_init();
        __UUPSUpgradeable_init();
        
        feeRecipient = _feeRecipient;
        ethUsdPriceFeed = _ethUsdPriceFeed;
        defaultERC20Token = _defaultERC20Token;
        defaultERC20UsdPriceFeed = _defaultERC20UsdPriceFeed;
        auctionImplementation = _auctionImplementation;
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
        
        bytes memory data = abi.encodeWithSelector(
            NFTChainlinkAuctionUUPS.initialize.selector,
            _nft,
            _nftId,
            msg.sender,
            _duration,
            _reservePrice,
            feePercentage,
            feeRecipient,
            ethUsdPriceFeed,
            defaultERC20Token,
            defaultERC20UsdPriceFeed
        );
        
        ERC1967Proxy proxy = new ERC1967Proxy(auctionImplementation, data);
        address auctionAddress = address(proxy);
        
        allAuctions.push(auctionAddress);
        getAuction[_nft][_nftId] = auctionAddress;
        isAuction[auctionAddress] = true;
        
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

    function getAllAuctions() external view returns (address[] memory) {
        return allAuctions;
    }

    function getAuctionsCount() external view returns (uint256) {
        return allAuctions.length;
    }
    
    function setFeePercentage(uint256 _feePercentage) external onlyOwner {
        require(_feePercentage <= 100, "Fee too high");
        emit FeePercentageUpdated(feePercentage, _feePercentage);
        feePercentage = _feePercentage;  // 更新手续费百分比
    }
    
    function setFeeRecipient(address _feeRecipient) external onlyOwner {
        require(_feeRecipient != address(0), "Invalid fee recipient");
        emit FeeRecipientUpdated(feeRecipient, _feeRecipient);
        feeRecipient = _feeRecipient;  // 更新手续费接收地址
    }
    
    function updateChainlinkFeeds(
        address _ethUsdPriceFeed,
        address _defaultERC20Token,
        address _defaultERC20UsdPriceFeed
    ) external onlyOwner {
        ethUsdPriceFeed = _ethUsdPriceFeed;  // 更新ETH/USD价格预言机地址
        defaultERC20Token = _defaultERC20Token;  // 更新默认ERC20代币地址
        defaultERC20UsdPriceFeed = _defaultERC20UsdPriceFeed;  // 更新Chainlink预言机地址
        emit ChainlinkFeedsUpdated(_ethUsdPriceFeed, _defaultERC20Token, _defaultERC20UsdPriceFeed);
    }
    
    function emergencyCancelAuction(address _nft, uint256 _nftId) external onlyOwner {
        address auctionAddress = getAuction[_nft][_nftId];  // 获取拍卖合约地址
        require(auctionAddress != address(0), "Auction not found");
        require(isAuction[auctionAddress], "Invalid auction");
        
        NFTChainlinkAuctionUUPS auction = NFTChainlinkAuctionUUPS(auctionAddress);  // 实例化拍卖合约
        auction.emergencyCancel();  // 调用拍卖合约的紧急取消函数
        
        isAuction[auctionAddress] = false;  // 标记拍卖为无效
    }
    
    function updateImplementation(address _newImplementation) external onlyOwner {
        require(_newImplementation != address(0), "Invalid implementation");
        auctionImplementation = _newImplementation;  // 更新逻辑合约地址
        emit ImplementationUpdated(_newImplementation);  // 触发事件
    }
    
    function _authorizeUpgrade(address newImplementation) internal override onlyOwner {}
}