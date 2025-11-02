const { ethers, deployments } = require("hardhat")
const { expect } = require("chai")
const path = require('path');

describe("Test auction", async function () {
    it("Should be ok", async function () {
        await main();
    });
})

async function main() {
    // 角色导向 - 强调能做什么 signer, buyer
    // 获取两个测试账户
    const [signer, buyer] = await ethers.getSigners()
    console.log(path.basename(__filename), "部署用户地址：", signer.address);
    console.log(path.basename(__filename), "买家用户地址：", buyer.address);

    // 部署合约 加载预定义的部署状态；deployments.fixture 会执行带有指定标签的部署脚本(在deploy目录下)
    await deployments.fixture(["depolyNftAuction"]);
    // 获取已部署合约的信息
    const nftAuctionProxy = await deployments.get("NftAuctionProxy"); // save 时的名称
    // 获取合约实例, 支持 合约方法调用
    const nftAuction = await ethers.getContractAt("NftAuction",nftAuctionProxy.address);


    // 1.获取合约工厂，用于部署，确保合约已编译并可访问 ABI 和字节码；
    const TestERC20 = await ethers.getContractFactory("TestERC20");
    // 2.部署测试用的 ERC20 代币合约（USDC）
    const testERC20 = await TestERC20.deploy();
    // 3.等待部署完成
    await testERC20.waitForDeployment(); 

    // 获取部署后的合约地址
    const UsdcAddress = await testERC20.getAddress();

    // 给买家转账一些 USDC 代币｜ethers.parseEther 将 "1000" ETH 转换为 wei 单位
    let tx = await testERC20.connect(signer).transfer(buyer, ethers.parseEther("1000"))
    await tx.wait()

    // 准备价格预言机合约
    const aggreagatorV3 = await ethers.getContractFactory("AggreagatorV3")
    // 部署 ETH 价格预言机 - 设置价格为 10000 used
    const priceFeedEthDeploy = await aggreagatorV3.deploy(ethers.parseEther("10000"))
    const priceFeedEth = await priceFeedEthDeploy.waitForDeployment()

    // 输出价格预言机地址
    const priceFeedEthAddress = await priceFeedEth.getAddress()
    console.log("ethFeed: ", priceFeedEthAddress)

    // 部署 USDC 价格预言机 - 设置价格为 1 usdc
    const priceFeedUSDCDeploy = await aggreagatorV3.deploy(ethers.parseEther("1"))
    const priceFeedUSDC = await priceFeedUSDCDeploy.waitForDeployment()

    // 输出价格预言机地址
    const priceFeedUSDCAddress = await priceFeedUSDC.getAddress()
    console.log("usdcFeed: ", priceFeedUSDCAddress)

    const token2Usd = [{
        token: ethers.ZeroAddress,  // ETH 资产地址 (0x000...000)
        priceFeed: priceFeedEthAddress  // ETH 预言机地址 (对应 10000 USD 的价格)
    }, {
        token: UsdcAddress,    // USDC 资产地址
        priceFeed: priceFeedUSDCAddress   // USDC 预言机地址 (对应 1 USD 的价格)
    }]

    // 设置代币到价格预言机的映射关系
    for (let i = 0; i < token2Usd.length; i++) {
        const { token, priceFeed } = token2Usd[i];
        await nftAuction.setPriceFeed(token, priceFeed);
    }


    // 1. 部署 ERC721 合约 NFT
    const TestERC721 = await ethers.getContractFactory("TestERC721");
    const testERC721 = await TestERC721.deploy();
    await testERC721.waitForDeployment();
    const testERC721Address = await testERC721.getAddress();
    console.log("testERC721Address::", testERC721Address);

    // mint 10个 NFT
    for (let i = 0; i < 10; i++) {
        await testERC721.mint(signer.address, i + 1);
    }

    const tokenId = 1;    

    // 给代理合约授权
    await testERC721.connect(signer).setApprovalForAll(nftAuctionProxy.address, true);

    await nftAuction.createAuction(
        10,
        ethers.parseEther("0.01"),
        testERC721Address,
        tokenId
    );

    const auction = await nftAuction.auctions(0);

    console.log("创建拍卖成功：：", auction);

    // 3. 购买者参与拍卖
    // await testERC721.connect(buyer).approve(nftAuctionProxy.address, tokenId);
    // ETH参与竞价；ethers.ZeroAddress, 表示使用 ETH 参与竞价
    tx = await nftAuction.connect(buyer).placeBid(0, 0, ethers.ZeroAddress, { value: ethers.parseEther("0.01") });
    await tx.wait()

    // USDC参与竞价； ethers.MaxUint256 无限授权：授权最大数量
    tx = await testERC20.connect(buyer).approve(nftAuctionProxy.address, ethers.MaxUint256)
    await tx.wait()
    tx = await nftAuction.connect(buyer).placeBid(0, ethers.parseEther("101"), UsdcAddress);
    await tx.wait()

    // 4. 结束拍卖
    // 等待 10 s；为了达到前面创建拍卖时设置的 10 s 拍卖时间
    await new Promise((resolve) => setTimeout(resolve, 10 * 1000));

    await nftAuction.connect(signer).endAuction(0);

    // 验证结果
    const auctionResult = await nftAuction.auctions(0);
    console.log("结束拍卖后读取拍卖成功：：", auctionResult);
    expect(auctionResult.highestBidder).to.equal(buyer.address);
    expect(auctionResult.highestBid).to.equal(ethers.parseEther("101"));

    // 验证 NFT 所有权
    const owner = await testERC721.ownerOf(tokenId);
    console.log("owner::", owner);
    expect(owner).to.equal(buyer.address);
}