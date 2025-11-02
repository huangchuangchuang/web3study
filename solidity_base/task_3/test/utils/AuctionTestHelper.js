const { ethers, deployments } = require("hardhat")
const { expect } = require("chai")


class AuctionTestHelper{
    constructor() {
        this.signer = null;
        this.ethBider = null;
        this.usdBider = null;
        this.contracts = {};
        this.priceFeeds = {};
        this.auctionId = 0;
        this.auction = null;
        this.ethRate = 10000; // 1 ETH = 10000 USD
        this.bidByEth = 0;
        this.bidByUsdc = 0;
        this.auctionDuration = 10; // 拍卖持续时间，单位：秒
        this.startPrice = "0.01";
        this.nextNftTokenId = 1;
        this.nftTokenId = 0;
        this.implAddress = {};
    }

    async endAuction() {
        // 结束拍卖; 等待 10 s；为了达到前面创建拍卖时设置的 10 s 拍卖时间
        await new Promise((resolve) => setTimeout(resolve, 10 * 1000));
        await this.contracts.nftAuction.connect(this.signer).endAuction(this.auctionId);
        // 获取结果
        const auctionResult = await this.contracts.nftAuction.auctions(this.auctionId);
        console.log("结束拍卖后读取拍卖成功：：", auctionResult);
        const owner = await this.contracts.testERC721.ownerOf(this.nftTokenId);
        console.log("owner::", owner);
        // 验证结果; 判断 bidEth 和 bidUsdc 哪个更高，决定最终赢家
        const ethBidUsd = parseFloat(this.bidByEth) * this.ethRate;
        const usdcBidUsd = parseFloat(this.bidByUsdc);
        console.log("ethBider:", this.ethBider.address, "ethBidUsd:", ethBidUsd)
        console.log("usdBider:", this.usdBider.address, "usdBidUsd:", usdcBidUsd);
        if (ethBidUsd > usdcBidUsd) {
            expect(auctionResult.highestBidder).to.equal(this.ethBider.address);
            expect(auctionResult.highestBid).to.equal(ethers.parseEther(this.bidByEth));
            // 验证 NFT 所有权
            expect(owner).to.equal(this.ethBider.address);
        } else {
            expect(auctionResult.highestBidder).to.equal(this.usdBider.address);
            expect(auctionResult.highestBid).to.equal(ethers.parseEther(this.bidByUsdc));
            // 验证 NFT 所有权
            expect(owner).to.equal(this.usdBider.address);
        }
    }

    async setup() {
        [this.signer, this.ethBider, this.usdBider] = await ethers.getSigners();
        console.log("部署用户地址：", this.signer.address);
        console.log("eth买家地址：", this.ethBider.address);
        console.log("usd买家地址：", this.usdBider.address);
        await deployments.fixture(["deployNftAuction", "deployTestToken", "deployOracles"]);
 
        this.contracts.nftAuction = await this.getContract("NftAuction", "NftAuctionProxy");
        this.contracts.testERC20 = await this.getContract("TestERC20", "TestERC20");
        this.contracts.testERC721 = await this.getContract("TestERC721", "TestERC721");
        this.contracts.ethPriceFeed = await this.getContract("AggregatorV3", "EthPriceFeed");
        this.contracts.usdcPriceFeed = await this.getContract("AggregatorV3", "UsdcPriceFeed");
        this.priceFeeds.eth = await this.contracts.ethPriceFeed.getAddress();
        this.priceFeeds.usdc = await this.contracts.usdcPriceFeed.getAddress();
        await this.setPriceFeed();  // 设置价格预言机映射关系
        await this.mintNFT(10);     // mint 10个 NFT
        await this.approveNFT(await this.contracts.nftAuction.getAddress());  // 给拍卖合约授权
        await this.transferERC20(this.ethBider.address, ethers.parseEther("1000")); // 给买家1转账一些 USDC 代币
        await this.transferERC20(this.usdBider.address, ethers.parseEther("1000")); // 给买家2转账一些 USDC 代币
        console.log("setup finished.");
        
    }

    async setPriceFeed() {
        console.log("this.priceFeeds.eth: ", this.priceFeeds.eth);
        console.log("this.priceFeeds.usdc: ", this.priceFeeds.usdc);
        const token2Usd = [{
            token: ethers.ZeroAddress,  // ETH 资产地址 (0x000...000)
            priceFeed: this.priceFeeds.eth  // ETH 预言机地址 (对应 10000 USD 的价格)
        }, {
            token: await this.contracts.testERC20.getAddress(),    // USDC 资产地址
            priceFeed: this.priceFeeds.usdc   // USDC 预言机地址 (对应 1 USD 的价格)
        }]

        // 设置代币到价格预言机的映射关系
        for (let i = 0; i < token2Usd.length; i++) {
            const { token, priceFeed } = token2Usd[i];
            await this.contracts.nftAuction.setPriceFeed(token, priceFeed);
        }
        console.log("设置价格预言机映射关系完成。");
    }

    async getContract(contractName, deployName) {
        const deployment = await deployments.get(deployName);
        return ethers.getContractAt(contractName, deployment.address);
    }
    async mintNFT(number) {
        for (let i = 0; i < number; i++) {
            await this.contracts.testERC721.mint(this.signer.address, i + 1);
        }
        console.log(`Mint ${number} NFTs to ${this.signer.address} successfully.`);
    }

    async approveNFT(toAddress) {
        await this.contracts.testERC721.connect(this.signer).setApprovalForAll(toAddress, true);
        console.log(`Approved NFT to ${toAddress} successfully.`);
    }

    async transferERC20(toAddress, amount) {
        let tx = await this.contracts.testERC20.connect(this.signer).transfer(toAddress, amount);
        await tx.wait();
        console.log(`Transferred ${ethers.formatEther(amount)} ERC20 to ${toAddress} successfully.`);
        return tx;
    }
    
    async auctionCreate() {
        this.nftTokenId = this.nextNftTokenId++; // 使用并递增
        const tx = await this.contracts.nftAuction.connect(this.signer).createAuction(
            this.auctionDuration,
            ethers.parseEther(this.startPrice),
            await this.contracts.testERC721.getAddress(),
            this.nftTokenId 
        );
        await tx.wait();
        this.auctionId = (await this.contracts.nftAuction.nextAuctionId()) - 1n;
        this.auction = await this.contracts.nftAuction.auctions(this.auctionId);
        console.log(`auctionId=${this.auctionId} 创建拍卖成功`, this.auction);
        return this.auction;
    }

    async placeEthBid(bidAmount) {
        if (bidAmount >= this.bidByEth) {
            this.bidByEth = bidAmount;
        }
        // ethers.parseEther 的作用是将字符串形式的以太币数量转换为 wei 单位的 BigNumber。
        let bidValue = ethers.parseEther(bidAmount);
        // ETH参与竞价；ethers.ZeroAddress, 表示使用 ETH 参与竞价; 
        let tx = await this.contracts.nftAuction.connect(this.ethBider).placeBid(
            this.auctionId, 0, ethers.ZeroAddress, { value: bidValue }
        );
        await tx.wait();
        console.log(`ethBider ${this.ethBider.address} placed ETH bid of ${ethers.formatEther(bidValue)} successfully.`);

    }
    
    async placeUsdBid(bidAmount) {
        if (bidAmount >= this.bidByUsdc) {
            this.bidByUsdc = bidAmount;
        }
        let bidValue = ethers.parseEther(bidAmount);
        // USDC参与竞价； ethers.MaxUint256 无限授权：授权最大数量
        let tx = await this.contracts.testERC20.connect(this.usdBider).approve(
            await this.contracts.nftAuction.getAddress(), ethers.MaxUint256
        );
        await tx.wait();
        tx = await this.contracts.nftAuction.connect(this.usdBider).placeBid(
            this.auctionId, bidValue, await this.contracts.testERC20.getAddress()
        );
        await tx.wait();
        console.log(`usdBider ${this.usdBider.address} placed USDC bid of ${ethers.formatEther(bidValue)} successfully.`);
    }

    async getImplAddress() {
        // getAddress() 返回的是代理地址，这是用户与之交互的地址。
        const implAddress = await upgrades.erc1967.getImplementationAddress(
            await this.contracts.nftAuction.getAddress()
        );
        return implAddress;
    }

    async upgradeNftAuction() {
        this.implAddress.before = await this.getImplAddress();
        console.log("升级前 NftAuction 实现合约地址：", this.implAddress.before);
        await deployments.fixture(["upgradeNftAuction"]);
        this.implAddress.after = await this.getImplAddress();
        console.log("升级后 NftAuction 实现合约地址：", this.implAddress.after);

    }

    async getNewNftAuctionContract() {
        // 代理合约地址不变，但实现合约地址变了；所以需要使用新的 ABI 来获取合约实例
        this.contracts.nftAuctionV2 = await ethers.getContractAt(
            "NftAuctionV2",
            await this.contracts.nftAuction.getAddress()
        );
        return this.contracts.nftAuctionV2;
    }

    async verifyNftAuctionUpgrade() {
        // 验证升级前后，实现合约地址不一样
        expect(this.implAddress.before).to.not.equal(this.implAddress.after);
        // 验证新方法可用
        const nftAuctionV2 = await this.getNewNftAuctionContract();
        const hello =  await nftAuctionV2.testHello();
        console.log("hello::", hello);
        // 验证历史数据未丢失
        const auction2 = await nftAuctionV2.auctions(this.auctionId);
        console.log("升级后读取拍卖成功：：", auction2);
        expect(auction2.startTime).to.equal(this.auction.startTime);
        return hello;
    }
}

module.exports = { AuctionTestHelper };

