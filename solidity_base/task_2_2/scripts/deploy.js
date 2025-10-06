const { ethers } = require("hardhat");

async function main() {
    const [deployer] = await ethers.getSigners();
    
    console.log("Deploying contracts with the account:", deployer.address);
    console.log("Account balance:", (await deployer.getBalance()).toString());
    
    // 获取合约工厂
    const SimpleNFT = await ethers.getContractFactory("SimpleNFT");
    
    // 部署合约
    console.log("Deploying SimpleNFT...");
    const simpleNFT = await SimpleNFT.deploy("MyDigitalArt", "MDA");
    
    await simpleNFT.deployed();
    
    console.log("SimpleNFT deployed to:", simpleNFT.address);
}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error);
        process.exit(1);
    });