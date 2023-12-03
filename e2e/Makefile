HARDHAT ?= npx hardhat
RLY_BIN ?= ./build/yrly

.PHONY: yrly
yrly:
	go build -o ./build/yrly ./relayer

.PHONY: deploy
deploy:
	$(HARDHAT) compile
	$(HARDHAT) --network chain0 deploy
	$(HARDHAT) --network chain1 deploy

.PHONY: network
network:
	docker compose up --wait -d

.PHONY: network-down
network-down:
	docker compose down --volumes --remove-orphans

.PHONY: init
init:
	RLY_BIN=$(RLY_BIN) ./relayer/scripts/init-rly

.PHONY: handshake
handshake:
	RLY_BIN=$(RLY_BIN) ./relayer/scripts/handshake

.PHONY: relay
relay:
	RLY_BIN=$(RLY_BIN) ./relayer/scripts/relay

.PHONY: test
test: yrly
	$(MAKE) network
	sleep 3
	$(MAKE) deploy
	$(MAKE) init
	$(MAKE) handshake
	$(MAKE) relay
	$(MAKE) network-down
