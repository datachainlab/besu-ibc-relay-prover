HARDHAT ?= npx hardhat
RLY_BIN ?= ./build/yrly
CONSENSUS_TYPE ?= qbft

.PHONY: yrly
yrly:
	go build -o ./build/yrly ./relayer

.PHONY: deploy
deploy:
	$(HARDHAT) compile
	$(HARDHAT) --network chain0 deploy
	$(HARDHAT) --network chain1 deploy

.PHONY: network-qbft
network-qbft:
	docker compose up --wait -d qbft-testchain0 qbft-testchain1

.PHONY: network-ibft2
network-ibft2:
	docker compose up --wait -d ibft2-testchain0 ibft2-testchain1

.PHONY: network-down
network-down:
	docker compose down --volumes --remove-orphans

.PHONY: init
init:
	./relayer/scripts/extract-abi.sh
	RLY_BIN=$(RLY_BIN) CONSENSUS_TYPE=$(CONSENSUS_TYPE) ./relayer/scripts/init-rly

.PHONY: handshake
handshake:
	RLY_BIN=$(RLY_BIN) ./relayer/scripts/handshake

.PHONY: relay
relay:
	RLY_BIN=$(RLY_BIN) ./relayer/scripts/relay

.PHONY: test
test: yrly
	$(MAKE) deploy
	$(MAKE) init
	$(MAKE) handshake
	$(MAKE) relay
	$(MAKE) network-down
