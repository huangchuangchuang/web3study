##



```
# 安装
npm install --save-dev hardhat
# 初始化项目
npx hardhat init
# 编译目录下所有合约
npx hardhat compile
# 测试
npx hardhat test

```

## 本地部署合约
```
# 本地启动持久化网络节点
npx hardhat node
# 本地@部署
npx hardhat ignition deploy ignition/modules/Counter.ts --network localhost

```
## 部署到测试网
 - 访问 infura.io 注册，获取 
 - 通过 MetaMost 钱包 获取账号私钥 SEPOLIA_PRIVATE_KEY    


```
# SEPOLIA_RPC_URL



# 创建.env文件配置 SEPOLIA_RPC_URL SEPOLIA_PRIVATE_KEY 
npm install dotenv
在 hardhat.config.ts import 'dotenv/config';


# 本地不需要部署持久化节点
npx hardhat ignition deploy ignition/modules/Counter.ts --network sepolia
# 查询合约 在 https://sepolia.etherscan.io/ 输入地址
0x65a6deb0f37A2c713F99854666e43BA0501201e2


# 部署到Sepolia测试网
npx hardhat run scripts/deploy.js --network sepolia

UUPS 升级的是逻辑合约，而不是代理合约
```

## 验证插件

```
npm install --save-dev @nomicfoundation/hardhat-verify

npx hardhat verify --network sepolia <合约地址> <构造函数参数(如果有)>
```

# 如果你想使用不同的配置文件

```
# 使用特定的配置文件
npx hardhat --config hardhat.config.custom.ts compile
npx hardhat --config hardhat.config.custom.ts test
npx hardhat --config hardhat.config.custom.ts ignition deploy ignition/modules/Counter.ts --network sepolia
```


