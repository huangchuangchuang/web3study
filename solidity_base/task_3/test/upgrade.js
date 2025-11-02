
const {AuctionTestHelper} = require("./utils/AuctionTestHelper");

describe("Test upgrade", async function () {
  it("Upgrade 1", async function () {
    const used = new AuctionTestHelper();
    await used.setup();
    await used.auctionCreate();
    // await used.placeEthBid("0.01");
    // await used.placeUsdBid("101");
    // await used.endAuction();
    await used.upgradeNftAuction();
    await used.verifyNftAuctionUpgrade();
  });
  // it("Upgrade 2", async function () {
  //   const used = new AuctionTestHelper();
  //   await used.setup();
  //   await used.auctionCreate();
  //   await used.upgradeNftAuction();
  //   await used.verifyNftAuctionUpgrade();
  // });
});
