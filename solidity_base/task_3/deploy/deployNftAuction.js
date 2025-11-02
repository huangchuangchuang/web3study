const { deployments, upgrades, ethers } = require("hardhat");

const fs = require("fs");
const path = require("path");

// 使用代理的方式部署合约
async function usedProxyDeployContract(name){
  // 1.获取合约工厂，用于部署，确保合约已编译并可访问 ABI 和字节码
  const NftAuction = await ethers.getContractFactory(name);
  // 2.使用工厂通过代理模式部署合约（实际部署），deployProxy: 交易已发送并被接收，不代表部署完成；默认用第一个账号部署
  const nftAuctionProxy = await upgrades.deployProxy(NftAuction, [], {
    initializer: "initialize",
  })
  // 3.等待部署完成；waitForDeployment: 交易已被打包并确认到区块中
  await nftAuctionProxy.waitForDeployment();

  // 获取合约 ABI 和地址信息
  abi = NftAuction.interface.format("json");
  const proxyAddress = await nftAuctionProxy.getAddress();
  const implAddress = await upgrades.erc1967.getImplementationAddress(proxyAddress);

  // 使用 Hardhat Deploy 插件保存部署信息到标准部署目录中
  // 这样可以让其他脚本通过 deployments.get() 方法获取到合约信息
  proxyName = `${name}Proxy`;
  await deployments.save(proxyName, {
    abi: abi,  // 合约 ABI
    address: proxyAddress,   // 合约地址
    // args: [],    // 部署参数（可选）
    // log: true,   // 是否记录日志（可选）
  })
  return { proxyName, proxyAddress, implAddress, abi};
}

function writeToFile(filename, data) {
    // 拼接绝对路径
  const storePath = path.resolve(__dirname, filename);
  // 将对象转换为 JSON 字符串
  jsonString = JSON.stringify(data);
  fs.writeFileSync(storePath, jsonString);
  return jsonString;

}


// Hardhat Deploy 使用依赖注入模式，将需要的工具函数作为参数传入：
// 这是主函数 - module.exports 导出的函数。async 关键字的作用是将函数标记为异步函数，允许在函数内部使用 await 关键字。
module.exports = async ({ getNamedAccounts }) => {
  // 职责导向 - 强调负责什么 deployer admin
  const { deployer } = await getNamedAccounts();
  const { proxyAddress, implAddress, abi, proxyName} = await usedProxyDeployContract("NftAuction");
  console.log("01部署用户地址：", deployer);
  console.log("01代理合约名称：", proxyName);
  console.log("01代理合约地址：", proxyAddress);
  console.log("01实现合约地址：", implAddress);
  console.log("合约ABI：", abi);

  // 将合约信息写入文件
  jsonString = writeToFile("./.cache/proxyNftAuction.json", {proxyAddress, implAddress, abi});
  // console.log("写入文件信息", jsonString);

  // // 直接部署（不支持升级）
  // await deploy("MyContract", {from: deployer, args: ["Hello"],log: true,});
};


module.exports.tags = ["deployNftAuction"];
// Hardhat Deploy插件推荐使用小驼峰命名法来命名部署脚本文件，因为：
// 部署脚本本质上是函数模块
// 与JavaScript社区对函数和变量的命名约定一致
// 便于通过--tags参数调用