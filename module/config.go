package module

import (
	"fmt"
	"time"

	"github.com/datachainlab/ethereum-ibc-relay-chain/pkg/relay/ethereum"
	"github.com/hyperledger-labs/yui-relayer/core"
)

const (
	QBFTConsensusType  = "qbft"
	IBFT2ConsensusType = "ibft2"
)

var _ core.ProverConfig = (*ProverConfig)(nil)

func (c ProverConfig) Build(chain core.Chain) (core.Prover, error) {
	chain_, ok := chain.(*ethereum.Chain)
	if !ok {
		return nil, fmt.Errorf("chain type must be %T, not %T", &ethereum.Chain{}, chain)
	}
	return NewProver(chain_, c), nil
}

func (c ProverConfig) Validate() error {
	if c.ConsensusType != "" && c.ConsensusType != QBFTConsensusType && c.ConsensusType != IBFT2ConsensusType {
		return fmt.Errorf("invalid consensus type: %s", c.ConsensusType)
	}
	if c.TrustingPeriod != "" {
		if _, err := time.ParseDuration(c.TrustingPeriod); err != nil {
			return fmt.Errorf("invalid trusting period: %s", c.TrustingPeriod)
		}
	}
	if c.MaxClockDrift != "" {
		if _, err := time.ParseDuration(c.MaxClockDrift); err != nil {
			return fmt.Errorf("invalid max clock drift: %s", c.MaxClockDrift)
		}
	}
	return nil
}

func (c ProverConfig) IsIBFT2() bool {
	return c.ConsensusType == IBFT2ConsensusType
}

func (c ProverConfig) GetTrustingPeriod() time.Duration {
	if c.TrustingPeriod == "" {
		return 0
	}
	d, err := time.ParseDuration(c.TrustingPeriod)
	if err != nil {
		panic(err)
	}
	return d
}

func (c ProverConfig) GetMaxClockDrift() time.Duration {
	if c.MaxClockDrift == "" {
		return 0
	}
	d, err := time.ParseDuration(c.MaxClockDrift)
	if err != nil {
		panic(err)
	}
	return d
}
