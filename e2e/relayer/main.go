package main

import (
	"log"

	"github.com/datachainlab/ethereum-ibc-relay-chain/pkg/relay/ethereum"
	"github.com/datachainlab/ethereum-ibc-relay-chain/pkg/relay/ethereum/signers/hd"
	"github.com/datachainlab/ibft2-relay-prover/module"
	"github.com/hyperledger-labs/yui-relayer/cmd"
)

func main() {
	if err := cmd.Execute(
		ethereum.Module{},
		hd.Module{},
		module.Module{},
	); err != nil {
		log.Fatal(err)
	}
}
