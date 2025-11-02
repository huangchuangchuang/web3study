
const {AuctionTestHelper} = require("./utils/AuctionTestHelper");

describe("Test upgrade", async function () {
  it("Upgrade 1", async function () {
    const cont_used = new AuctionTestHelper();
    await cont_used.setup();
    await cont_used.auctionCreate();
    // await cont_used.placeEthBid("0.01");
    // await cont_used.placeUsdBid("101");
    // await cont_used.endAuction();
    await cont_used.upgradeNftAuction();
    await cont_used.verifyNftAuctionUpgrade();
  });
  // it("Upgrade 2", async function () {
  //   const cont_used = new AuctionTestHelper();
  //   await cont_used.setup();
  //   await cont_used.auctionCreate();
  //   await cont_used.upgradeNftAuction();
  //   await cont_used.verifyNftAuctionUpgrade();
  // });
});
