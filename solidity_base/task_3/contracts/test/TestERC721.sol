// SPDX-License-Identifier: MIT
pragma solidity ^0.8.17;

import "@openzeppelin/contracts/token/ERC721/extensions/ERC721Enumerable.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

contract TestERC721 is ERC721Enumerable, Ownable {
    string private _tokenURI;

    // constructor() ERC721("Troll", "Troll") Ownable(msg.sender) {}
    // 参数化构造函数
    constructor(string memory name, string memory symbol) 
        ERC721(name, symbol) 
        Ownable(msg.sender) 
    {}

    function mint(address to, uint256 tokenId) external onlyOwner {
        _mint(to, tokenId);
    }

    function tokenURI(uint256 tokenId) public view override returns (string memory) {
        return _tokenURI;
    }

    function setTokenURI(string memory newTokenURI) external onlyOwner {
        _tokenURI = newTokenURI;
    }
}


// ERC721Enumerable 是 ERC721 标准的一个扩展，提供了枚举功能，允许按索引查询NFT信息。
// // 获取合约中NFT的总供应量
// function totalSupply() public view returns (uint256)

// // 根据索引获取tokenId
// function tokenByIndex(uint256 index) public view returns (uint256)

// // 获取某个拥有者拥有的NFT数量
// function balanceOf(address owner) public view returns (uint256)

// // 获取某个拥有者特定索引的tokenId
// function tokenOfOwnerByIndex(address owner, uint256 index) public view returns (uint256)