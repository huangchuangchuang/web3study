
const {AuctionTestHelper} = require("./utils/AuctionTestHelper");

describe("Test auction", async function () {

    it("bid two", async function () {
        const cont_used = new AuctionTestHelper();
        await cont_used.setup();
        await cont_used.auctionCreate();
        await cont_used.placeEthBid("0.01");
        await cont_used.placeUsdBid("101");
        await cont_used.endAuction();
    });
    it("bid three", async function () {
        const cont_used = new AuctionTestHelper();
        await cont_used.setup();
        await cont_used.auctionCreate();
        // 竞价必须大于当前最高价；否则合约会 revert
        await cont_used.placeEthBid("0.01");
        await cont_used.placeUsdBid("101");
        await cont_used.placeEthBid("0.011");
        await cont_used.endAuction();
    });
})
