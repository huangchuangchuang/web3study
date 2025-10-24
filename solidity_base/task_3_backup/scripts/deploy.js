// deploy.js
const { ethers, network } = require("hardhat");

async function main() {
  // 根据网络设置Chainlink预言机地址
  let ethUsdPriceFeed, defaultERC20Token, defaultERC20UsdPriceFeed;
  
  if (network.name === "goerli") {
    // Goerli测试网的Chainlink预言机地址
    ethUsdPriceFeed = "0xD4a33860578De61DBAbDc8BFdb98FD742fA7028e"; // ETH/USD
    defaultERC20Token = "0x7af963cF44225351d23463793486C2d267c0dD85"; // 测试LINK代币
    defaultERC20UsdPriceFeed = "0x777D47a138136714A01f6B49DD1e31540776f43E"; // LINK/USD
  } else if (network.name === "sepolia") {
    // Sepolia测试网的Chainlink预言机地址
    ethUsdPriceFeed = "0x694AA1769357215DE4FAC081bf1f309aDC325306";
    defaultERC20Token = ethers.constants.AddressZero;
    defaultERC20UsdPriceFeed = ethers.constants.AddressZero;
  } else {
    // 本地测试环境
    ethUsdPriceFeed = ethers.constants.AddressZero;
    defaultERC20Token = ethers.constants.AddressZero;
    defaultERC20UsdPriceFeed = ethers.constants.AddressZero;
  }

  // 部署拍卖合约实现
  const NFTAuction = await ethers.getContractFactory("NFTAuction");
  const auctionImplementation = await NFTAuction.deploy();
  await auctionImplementation.deployed();
  console.log("NFTAuction implementation deployed to:", auctionImplementation.address);

  // 部署工厂合约
  const [deployer] = await ethers.getSigners();
  const AuctionFactory = await ethers.getContractFactory("AuctionFactory");
  const factory = await AuctionFactory.deploy(
    deployer.address,           // feeRecipient
    ethUsdPriceFeed,            // ethUsdPriceFeed
    defaultERC20Token,          // defaultERC20Token
    defaultERC20UsdPriceFeed    // defaultERC20UsdPriceFeed
  );
  await factory.deployed();
  console.log("AuctionFactory deployed to:", factory.address);

  console.log("Deployment completed!");
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});