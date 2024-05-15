# E2E test

This is an e2e demo and test between two hyperledger besu chains using yui-relayer.

## Requirements

First, need to install go >= v1.20 and node >= v16.

Also, need to install npm dependencies:
```sh
$ npm install
```

## How to run

Just execute the single command:

```sh
$ make test
```

The above commands execute the following in sequence:

1. Launch two HB chains (both chain uses QBFT consensus)
   - `make network-qbft`
2. Deploy the IBC contracts from yui-ibc-solidity to the chains using hardhat
   - `make deploy`
3. Configurate yui-relayer setting with [./relayer/configs](./relayer/configs/)
   - `make init`
4. Perform IBC handshake using yui-relayer
   - `make handshake`
5. Send a packet using [./scripts/sendPacket.js](./scripts/sendPacket.js) and relay it using yui-relayer
   - `make relay`
6. Shutdown two HB chains
   - `make network-down`
