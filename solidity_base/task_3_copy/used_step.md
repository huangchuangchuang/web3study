# ready

```
# 安装 hardhat
npm install --save-dev hardhat

# 初始化项目,直接回车，安装相关吗默认组件
npx hardhat@2.26.3 init

# 安装依赖1
npm install @openzeppelin/contracts

# 安装依赖2，安装后在 hardhat.config.js 中添加：require("hardhat-deploy");
npm install --save-dev hardhat-deploy

# 安装依赖3，安装后在 hardhat.config.js 中添加：require("@openzeppelin/hardhat-upgrades")
npm install --save-dev @openzeppelin/hardhat-upgrades

# 安装依赖4 比较慢|预言机集成｜实际测试是本地模拟未使用上
npm install @chainlink/contracts
```

# run
```
# 创建目录用于升级合约，本地数据读取
mkdir -p deploy/.cache

# 执行所有测试用例
npx hardhat test

# 执行单个用例
npx hardhat test test/auction.js 
```
