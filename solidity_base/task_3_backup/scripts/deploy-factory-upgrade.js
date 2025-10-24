async function upgrade() {
  const FactoryV2 = await ethers.getContractFactory("NFTChainlinkAuctionFactoryUUPSv2");
  
  // 升级到新版本
  const factory = await upgrades.upgradeProxy(
    "0x...", // 现有代理地址
    FactoryV2
  );
  
  console.log("Factory upgraded");
}

upgrade().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});