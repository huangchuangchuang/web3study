const { deployments, upgrades, ethers } = require("hardhat");

const fs = require("fs");
const path = require("path");

module.exports = async () => {
    const { deploy } = deployments;
    
    // 部署 TestERC20 合约
    const testERC20 = await deploy("TestERC20", {
        contract: "TestERC20",
        from: (await ethers.getSigners())[0].address,
        args: [],
        log: true,
    });
    console.log("TestERC20 deployed to:", testERC20.address);

    // 部署 TestERC721 合约
    const testERC721 = await deploy("TestERC721", {
        contract: "TestERC721",
        from: (await ethers.getSigners())[0].address,
        args: ["MyNFT", "MNFT"],
        log: true,
    });
    console.log("TestERC721 deployed to:", testERC721.address);
}


module.exports.tags = ["deployTestToken"];