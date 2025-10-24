const { ethers, upgrades } = require("hardhat");

async function main() {
  // 部署 UUPS 逻辑合约
  const Factory = await ethers.getContractFactory("NFTChainlinkAuctionFactoryUUPS");
  
  // 部署可升级代理
  const factory = await upgrades.deployProxy(Factory, [
    feeRecipient,
    ethUsdPriceFeed,
    defaultERC20Token,
    defaultERC20UsdPriceFeed,
    auctionImplementation
  ], {
    initializer: 'initialize'
  });
  
  console.log("Factory deployed to:", factory.address);
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});