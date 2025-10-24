import { network } from "hardhat";

async function main() {
  const { viem } = await network.connect();
  const publicClient = await viem.getPublicClient();
  
  const contractAddress = "0x65a6deb0f37A2c713F99854666e43BA0501201e2"; // 替换为你的合约地址
  
  try {
    // 检查合约代码是否存在
    const code = await publicClient.getCode({ address: contractAddress });
    
    if (code && code !== '0x') {
      console.log("✅ 合约部署成功！");
      console.log("合约地址:", contractAddress);
      
      // 尝试读取合约状态
      const counter = await viem.getContractAt("Counter", contractAddress);
      const value = await counter.read.x();
      console.log("当前计数器值:", value.toString());
    } else {
      console.log("❌ 合约未部署或部署失败");
    }
  } catch (error) {
    console.error("检查部署时出错:", error);
  }
}

main().catch((error) => {
  console.error(error);
  process.exitCode = 1;
});