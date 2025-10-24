import { ethers } from "hardhat";

async function main() {
  console.log("🚀 开始部署 Counter 合约...");
  
  try {
    // 获取合约工厂
    const Counter = await ethers.getContractFactory("Counter");
    
    // 部署合约
    const counter = await Counter.deploy();
    await counter.waitForDeployment();
    
    const contractAddress = await counter.getAddress();
    console.log("✅ 合约部署成功！");
    console.log("📋 合约地址:", contractAddress);
    
    // 测试合约功能
    const initialValue = await counter.x();
    console.log("📊 初始值:", initialValue.toString());
    
    const tx = await counter.inc();
    await tx.wait();
    
    const newValue = await counter.x();
    console.log("📊 更新后的值:", newValue.toString());
    
  } catch (error) {
    console.error("❌ 部署过程中出错:", error);
    process.exitCode = 1;
  }
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});