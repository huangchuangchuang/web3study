
const {AuctionTestHelper} = require("./utils/AuctionTestHelper");

describe("Test auction", async function () {

    it("bid two", async function () {
        const used = new AuctionTestHelper();
        await used.setup();
        await used.auctionCreate();
        await used.placeEthBid("0.01");
        await used.placeUsdBid("101");
        await used.endAuction();
    });
    it("bid three", async function () {
        const used = new AuctionTestHelper();
        await used.setup();
        await used.auctionCreate();
        // 竞价必须大于当前最高价；否则合约会 revert
        await used.placeEthBid("0.01");
        await used.placeUsdBid("101");
        await used.placeEthBid("0.011");
        await used.endAuction();
    });
})
