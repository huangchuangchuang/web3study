const { deployments, ethers } = require("hardhat");

module.exports = async () => {
    const { deploy } = deployments;
    
    // 部署 ETH 价格预言机 - 设置价格为 10000 USD
    const ethPriceFeed = await deploy("EthPriceFeed", {
        contract: "AggregatorV3",
        from: (await ethers.getSigners())[0].address,
        args: [ethers.parseEther("10000")], // 1 ETH = 10000 USD
        log: true,
    });
    console.log("ETH Price Feed deployed to:", ethPriceFeed.address);

    // 部署 USDC 价格预言机 - 设置价格为 1 USD
    const usdcPriceFeed = await deploy("UsdcPriceFeed", {
        contract: "AggregatorV3",
        from: (await ethers.getSigners())[0].address,
        args: [ethers.parseEther("1")], // 1 USDC = 1 USD
        log: true,
    });
    console.log("USDC Price Feed deployed to:", usdcPriceFeed.address);
}

module.exports.tags = ["deployOracles"];