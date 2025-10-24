const { ethers, upgrades } = require("hardhat");

async function main() {
  console.log("Deploying upgradeable contracts...");
  
  // 部署可升级的拍卖合约实现
  const NFTChainlinkAuctionUUPS = await ethers.getContractFactory("NFTChainlinkAuctionUUPS");
  const auctionImplementation = await upgrades.deployImplementation(NFTChainlinkAuctionUUPS, {
    kind: "uups"
  });
  console.log("NFTChainlinkAuctionUUPS implementation deployed to:", auctionImplementation.address);

  // 部署可升级的工厂合约
  const NFTChainlinkAuctionFactoryUUPS = await ethers.getContractFactory("NFTChainlinkAuctionFactoryUUPS");
  const factory = await upgrades.deployProxy(
    NFTChainlinkAuctionFactoryUUPS,
    [
      "0x0000000000000000000000000000000000000000", // feeRecipient (will be updated)
      "0x0000000000000000000000000000000000000000", // ethUsdPriceFeed
      "0x0000000000000000000000000000000000000000", // defaultERC20Token
      "0x0000000000000000000000000000000000000000", // defaultERC20UsdPriceFeed
      auctionImplementation.address
    ],
    { 
      initializer: 'initialize',
      kind: "uups"
    }
  );
  
  await factory.deployed();
  console.log("NFTChainlinkAuctionFactoryUUPS deployed to:", factory.address);
  
  // 获取实现地址
  const implementationAddress = await upgrades.erc1967.getImplementationAddress(factory.address);
  console.log("Factory implementation address:", implementationAddress);
  
  console.log("Deployment completed!");
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});