// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "./NFTAuction.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import "@openzeppelin/contracts/token/ERC721/IERC721.sol";

contract AuctionFactory is Ownable, ReentrancyGuard {
    // 存储所有创建的拍卖合约地址
    address[] public allAuctions;
    
    // 映射：NFT合约地址 => NFT ID => 拍卖合约地址
    // 类似于 Uniswap V2 的 getPair mapping
    mapping(address => mapping(uint256 => address)) public getAuction;
    
    // 映射：拍卖合约地址 => 是否为有效拍卖
    mapping(address => bool) public isAuction;
    
    // 手续费相关（可选功能）
    uint256 public feePercentage = 2; // 2% 手续费
    address public feeRecipient;  // 手续费接收者地址
    
    // 事件
    event AuctionCreated(
        address indexed auction,      // 新创建的拍卖合约地址
        address indexed nftContract,  // NFT合约地址
        uint256 indexed nftId,        // NFT标识符ID
        address seller,               // 卖家地址
        uint256 duration,             // 拍卖持续时间
        uint256 fee                   // 手续费百分比
    );
    
    event FeePercentageUpdated(uint256 oldFee, uint256 newFee);  // 手续费比例更新事件
    event FeeRecipientUpdated(address oldRecipient, address newRecipient);  // 手续费接收者更新事件

    constructor(address _feeRecipient) Ownable(msg.sender){
        feeRecipient = _feeRecipient;  // 设置初始手续费接收者
    }

    /**
     * @dev 创建新的拍卖
     * @param _nft NFT合约地址
     * @param _nftId NFT ID
     * @param _duration 拍卖持续时间（秒）
     * @param _reservePrice 保留价格（wei）
     */
    function createAuction(
        address _nft,
        uint256 _nftId,
        uint256 _duration,
        uint256 _reservePrice
    ) external nonReentrant returns (address) {
        require(_duration > 0, "Duration must be positive");  // 持续时间必须为正数
        require(_nft != address(0), "Invalid NFT address");
        require(_duration <= 365 days, "Duration too long");
        
        // 检查是否已存在该NFT的拍卖
        address existingAuction = getAuction[_nft][_nftId];
        require(existingAuction == address(0) || !isAuction[existingAuction], 
                "Auction already exists for this NFT");
        
        // 创建新的拍卖合约
        NFTAuction auction = new NFTAuction(
            _nft,
            _nftId,
            msg.sender,  // 卖家
            _duration,
            _reservePrice,
            feePercentage,
            feeRecipient
        );
        
        // 获取新拍卖合约地址
        address auctionAddress = address(auction);
        
        // 记录拍卖合约
        allAuctions.push(auctionAddress);  // 存储到拍卖列表
        getAuction[_nft][_nftId] = auctionAddress;  // 映射NFT到拍卖合约
        isAuction[auctionAddress] = true;  // 标记为有效拍卖
        
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

    /**
     * @dev 获取所有拍卖合约地址
     */
    function getAllAuctions() external view returns (address[] memory) {
        return allAuctions;
    }

    /**
     * @dev 获取拍卖数量
     */
    function getAuctionsCount() external view returns (uint256) {
        return allAuctions.length;
    }
    
    /**
     * @dev 更新手续费比例
     */
    function setFeePercentage(uint256 _feePercentage) external onlyOwner {
        require(_feePercentage <= 100, "Fee too high"); // 最多10%
        emit FeePercentageUpdated(feePercentage, _feePercentage);
        feePercentage = _feePercentage;
    }
    
    /**
     * @dev 更新手续费接收者
     */
    function setFeeRecipient(address _feeRecipient) external onlyOwner {
        require(_feeRecipient != address(0), "Invalid fee recipient");
        emit FeeRecipientUpdated(feeRecipient, _feeRecipient);
        feeRecipient = _feeRecipient;
    }
    
    /**
     * @dev 暂停特定拍卖（紧急情况）
     */
    function emergencyCancelAuction(address _nft, uint256 _nftId) external onlyOwner {
        address auctionAddress = getAuction[_nft][_nftId];
        require(auctionAddress != address(0), "Auction not found");
        require(isAuction[auctionAddress], "Invalid auction");
        
        NFTAuction auction = NFTAuction(auctionAddress);
        auction.emergencyCancel();
        
        // 标记拍卖为无效
        isAuction[auctionAddress] = false;
    }
}