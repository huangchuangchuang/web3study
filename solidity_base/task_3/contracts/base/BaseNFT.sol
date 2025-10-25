// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import "@openzeppelin/contracts/token/ERC721/extensions/ERC721Enumerable.sol";
import "@openzeppelin/contracts/access/Ownable.sol";


// 同质化NFT场景：适合活动门票、会员资格、收藏卡牌等场景
// 1. 批量铸造和管理：支持批量铸造NFT，方便管理大量同质化资产
// 2. 可枚举性：通过ERC721Enumerable扩展，支持按索引枚举所有NFT，便于查询和展示
// 3. 所有权控制：通过Ownable合约，只有合约拥有者可以执行铸造和URI设置操作，确保安全性
contract BaseNFT is ERC721Enumerable, Ownable {
    string private _tokenURI;

    constructor() ERC721("Magic NFT", "MNFT") Ownable(msg.sender) {}


    function mint(address to, uint256 tokenId) external onlyOwner {
        _mint(to, tokenId);
    }

    function tokenURI(uint256 tokenId) public view override returns (string memory) {
        _requireOwned(tokenId); // 检查token是否存在且被拥有
        return _tokenURI;
    }

    function setTokenURI(string memory newTokenURI) external onlyOwner {
        _tokenURI = newTokenURI;
    }
}