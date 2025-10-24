// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

contract MyNFT is ERC721, Ownable {
    uint256 private _nextTokenId = 1;
    
    // 存储每个 token 的 URI
    mapping(uint256 => string) private _tokenURIs;

    constructor() 
    ERC721("MyArtCollection", "ART")
    Ownable(msg.sender)
    {}

    // 铸造 NFT
    function mintNFT(address to, string memory uri) public onlyOwner returns (uint256) {
        // to = recipient 你的钱包地址。
        // uri = tokenURI 元数据的 IPFS 链接。
        uint256 tokenId = _nextTokenId++;
        _mint(to, tokenId);
        _tokenURIs[tokenId] = uri;
        return tokenId;
    }

    // 重写 tokenURI 方法以返回自定义 URI
    function tokenURI(uint256 tokenId) public view override returns (string memory) {
        _requireOwned(tokenId); // 检查 token 是否存在
        return _tokenURIs[tokenId];
    }

    // 可选：添加批量铸造功能
    function batchMintNFT(address to, string[] memory uris) public onlyOwner returns (uint256[] memory) {
        uint256[] memory tokenIds = new uint256[](uris.length);
        for (uint256 i = 0; i < uris.length; i++) {
            tokenIds[i] = mintNFT(to, uris[i]);
        }
        return tokenIds;
    }
}