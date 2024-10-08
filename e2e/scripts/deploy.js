const qbftClientType = "hb-qbft";
const mockAppPortId = "mockapp";

async function deploy(deployer, contractName, args = []) {
  const factory = await hre.ethers.getContractFactory(contractName);
  const contract = await factory.connect(deployer).deploy(...args);
  await contract.waitForDeployment();
  return contract;
}

async function deployIBC(deployer) {
  const logicNames = [
    "IBCClient",
    "IBCConnectionSelfStateNoValidation",
    "IBCChannelHandshake",
    "IBCChannelPacketSendRecv",
    "IBCChannelPacketTimeout",
    "IBCChannelUpgradeInitTryAck",
    "IBCChannelUpgradeConfirmOpenTimeoutCancel"
  ];
  const logics = [];
  for (const name of logicNames) {
    const logic = await deploy(deployer, name);
    logics.push(logic);
  }
  return deploy(deployer, "OwnableIBCHandler", logics.map(l => l.target));
}

function saveContractAddresses(addresses) {
  const path = require("path");
  const fs = require("fs");
  const envFile = path.join(__dirname, "..", network.name + ".env.sh");
  let content = "";
  for (const [key, value] of Object.entries(addresses)) {
    content += `export ${key}=${value}\n`;
  }
  fs.writeFileSync(envFile, content);
}

task("deploy", "Deploy the contracts")
  .setAction(async (taskArgs, hre) => {
    // This is just a convenience check
    if (network.name === "hardhat") {
      console.warn(
        "You are trying to deploy a contract to the Hardhat Network, which" +
          "gets automatically created and destroyed every time. Use the Hardhat" +
          " option '--network localhost'"
      );
    }

    // ethers is available in the global scope
    const [deployer] = await hre.ethers.getSigners();
    console.log(
      "Deploying the contracts with the account:",
      await deployer.getAddress()
    );
    console.log("Account balance:", (await hre.ethers.provider.getBalance(deployer.getAddress())).toString());

    const ibcHandler = await deployIBC(deployer);
    console.log("IBCHandler address:", ibcHandler.target);

    const qbftClient = await deploy(deployer, "QBFTClient", [ibcHandler.target]);
    console.log("QBFTClient address:", qbftClient.target);

    const ibcMockApp = await deploy(deployer, "IBCMockApp", [ibcHandler.target]);
    console.log("IBCMockApp address:", ibcMockApp.target);

    await ibcHandler.registerClient(qbftClientType, qbftClient.target);
    await ibcHandler.bindPort(mockAppPortId, ibcMockApp.target);

    saveContractAddresses({
      IBC_HANDLER: ibcHandler.target,
      QBFT_CLIENT: qbftClient.target,
      IBC_MOCKAPP: ibcMockApp.target
    });
  });
