package module

import (
	"bytes"
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	"github.com/cosmos/ibc-go/v8/modules/core/exported"
	"github.com/datachainlab/ethereum-ibc-relay-chain/pkg/relay/ethereum"
	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/hyperledger-labs/yui-relayer/core"
)

// keccak256(abi.encode(uint256(keccak256("ibc.commitment")) - 1)) & ~bytes32(uint256(0xff))
var IBCCommitmentsSlot = common.HexToHash("1ee222554989dda120e26ecacf756fe1235cd8d726706b57517715dde4f0c900")

type Prover struct {
	chain  *ethereum.Chain
	config ProverConfig
}

var _ core.Prover = (*Prover)(nil)

func NewProver(chain *ethereum.Chain, config ProverConfig) *Prover {
	return &Prover{chain: chain, config: config}
}

// Init implements Prover.Init
func (pr *Prover) Init(homePath string, timeout time.Duration, codec codec.ProtoCodecMarshaler, debug bool) error {
	return nil
}

// SetRelayInfo implements Prover.SetRelayInfo
func (pr *Prover) SetRelayInfo(path *core.PathEnd, counterparty *core.ProvableChain, counterpartyPath *core.PathEnd) error {
	return nil
}

// SetupForRelay implements Prover.SetupForRelay
func (pr *Prover) SetupForRelay(ctx context.Context) error {
	return nil
}

// CreateInitialLightClientState implements Prover.CreateInitialLightClientState
func (pr *Prover) CreateInitialLightClientState(height exported.Height) (exported.ClientState, exported.ConsensusState, error) {
	var blockNumber *big.Int
	if height == nil {
		blockNumber = nil
	} else {
		blockNumber = big.NewInt(int64(height.GetRevisionHeight()))
	}

	header, err := pr.chain.Client().HeaderByNumber(context.Background(), blockNumber)
	if err != nil {
		return nil, nil, err
	}
	extra, err := parseExtraData(header.Extra)
	if err != nil {
		return nil, nil, err
	}
	proof, err := pr.chain.Client().GetProof(pr.chain.Config().IBCAddress(), nil, big.NewInt(int64(header.Number.Int64())))
	if err != nil {
		return nil, nil, err
	}
	var validators [][]byte
	for _, val := range extra.Validators {
		validators = append(validators, val.Bytes())
	}
	var chainIDUint256 [32]byte
	big.NewInt(int64(pr.chain.Config().EthChainId)).FillBytes(chainIDUint256[:])
	clientState := &ClientState{
		ChainId:         chainIDUint256[:],
		IbcStoreAddress: pr.chain.Config().IBCAddress().Bytes(),
		LatestHeight:    clienttypes.NewHeight(0, uint64(header.Number.Int64())),
		TrustingPeriod:  uint64(pr.config.GetTrustingPeriod().Seconds()),
		MaxClockDrift:   uint64(pr.config.GetMaxClockDrift().Seconds()),
	}
	consensusState := &ConsensusState{
		Timestamp:  header.Time,
		Root:       proof.StorageHash[:],
		Validators: validators,
	}
	return clientState, consensusState, nil
}

// GetLatestFinalizedHeader implements Prover.GetLatestFinalizedHeader
func (pr *Prover) GetLatestFinalizedHeader() (latestFinalizedHeader core.Header, err error) {
	return pr.getHeader(context.TODO(), nil)
}

// SetupHeadersForUpdate implements Prover.SetupHeadersForUpdate
func (pr *Prover) SetupHeadersForUpdate(counterparty core.FinalityAwareChain, latestFinalizedHeader core.Header) ([]core.Header, error) {
	header, ok := latestFinalizedHeader.(*Header)
	if !ok {
		return nil, fmt.Errorf("invalid header type: %T", latestFinalizedHeader)
	}
	if err := header.ValidateBasic(); err != nil {
		return nil, err
	}
	latestHeight, err := counterparty.LatestHeight()
	if err != nil {
		return nil, err
	}
	counterpartyClientRes, err := counterparty.QueryClientState(core.NewQueryContext(context.TODO(), latestHeight))
	if err != nil {
		return nil, err
	}
	var cs exported.ClientState
	if err := pr.chain.Codec().UnpackAny(counterpartyClientRes.ClientState, &cs); err != nil {
		return nil, err
	}
	header.TrustedHeight = cs.GetLatestHeight().(clienttypes.Height)
	return []core.Header{header}, nil
}

// ProveState implements Prover.ProveState
func (pr *Prover) ProveState(ctx core.QueryContext, path string, value []byte) ([]byte, clienttypes.Height, error) {
	proofHeight := int64(ctx.Height().GetRevisionHeight())
	height := pr.newHeight(proofHeight)
	proof, err := pr.buildStateProof([]byte(path), proofHeight)
	return proof, height, err
}

// ProveHeader implements Prover.ProveHostConsensusState
func (pr *Prover) ProveHostConsensusState(ctx core.QueryContext, height exported.Height, consensusState exported.ConsensusState) (proof []byte, err error) {
	return clienttypes.MarshalConsensusState(pr.chain.Codec(), consensusState)
}

// CheckRefreshRequired implements Prover.CheckRefreshRequired
func (pr *Prover) CheckRefreshRequired(counterparty core.ChainInfoICS02Querier) (bool, error) {
	// TODO implement
	return false, nil
}

func (pr *Prover) newHeight(blockNumber int64) clienttypes.Height {
	return clienttypes.NewHeight(0, uint64(blockNumber))
}

func (pr *Prover) buildStateProof(path []byte, height int64) ([]byte, error) {
	// calculate slot for commitment
	storageKey := crypto.Keccak256Hash(append(
		crypto.Keccak256Hash(path).Bytes(),
		IBCCommitmentsSlot.Bytes()...,
	))
	storageKeyHex, err := storageKey.MarshalText()
	if err != nil {
		return nil, err
	}

	// call eth_getProof
	stateProof, err := pr.chain.Client().GetProof(
		pr.chain.Config().IBCAddress(),
		[][]byte{storageKeyHex},
		big.NewInt(height),
	)
	if err != nil {
		return nil, err
	}
	return stateProof.StorageProofRLP[0], nil
}

func (pr *Prover) getHeader(ctx context.Context, bn *big.Int) (*Header, error) {
	header, err := pr.chain.Client().HeaderByNumber(ctx, bn)
	if err != nil {
		return nil, err
	}
	extra, err := parseExtraData(header.Extra)
	if err != nil {
		return nil, err
	}
	headerBytes, seals, err := pr.validateAndGetOrderedSeals(*header, *extra)
	if err != nil {
		return nil, err
	}
	proof, err := pr.chain.Client().GetProof(pr.chain.Config().IBCAddress(), nil, big.NewInt(int64(header.Number.Int64())))
	if err != nil {
		return nil, err
	}
	return &Header{
		BesuHeaderRlp:     headerBytes,
		Seals:             seals,
		AccountStateProof: proof.AccountProofRLP,
	}, nil
}

type ExtraData struct {
	Vanity     []byte
	Validators []common.Address
	Vote       interface{}
	Round      []byte
	Seals      [][]byte
}

func (pr *Prover) validateAndGetOrderedSeals(header gethtypes.Header, extra ExtraData) ([]byte, [][]byte, error) {
	var (
		extraBytes []byte
		err        error
	)
	if pr.config.IsIBFT2() {
		extraBytes, err = rlp.EncodeToBytes([]interface{}{
			extra.Vanity, extra.Validators, extra.Vote, extra.Round,
		})
	} else {
		extraBytes, err = rlp.EncodeToBytes([]interface{}{
			extra.Vanity, extra.Validators, extra.Vote, extra.Round, [][]byte{},
		})
	}
	if err != nil {
		return nil, nil, err
	}
	header.Extra = extraBytes
	headerBytes, err := rlp.EncodeToBytes(&header)
	if err != nil {
		return nil, nil, err
	}
	vals, err := recoverSeals(headerBytes, extra.Seals)
	if err != nil {
		return nil, nil, err
	}
	var orderedSeals [][]byte
	count := 0
	for _, val := range extra.Validators {
		if seal, ok := vals[val]; ok {
			count++
			orderedSeals = append(orderedSeals, seal)
		} else {
			orderedSeals = append(orderedSeals, nil)
		}
	}
	if threshold := len(extra.Validators) * 2 / 3; count > threshold {
		return headerBytes, orderedSeals, nil
	} else {
		return nil, nil, fmt.Errorf("insufficient voting: %v > %v", count, threshold)
	}
}

func recoverSeals(headerBytes []byte, seals [][]byte) (map[common.Address][]byte, error) {
	headerHash := crypto.Keccak256(headerBytes)
	vals := make(map[common.Address][]byte)
	for _, seal := range seals {
		addr, err := ecrecover(headerHash, seal[:])
		if err != nil {
			return nil, err
		}
		vals[addr] = seal[:]
	}
	return vals, nil
}

func ecrecover(hash, sig []byte) (common.Address, error) {
	pub, err := crypto.SigToPub(hash, sig)
	if err != nil {
		return common.Address{}, err
	}
	return crypto.PubkeyToAddress(*pub), nil
}

func parseExtraData(extraBytes []byte) (*ExtraData, error) {
	var extra ExtraData
	r := bytes.NewReader(extraBytes)
	stream := rlp.NewStream(r, uint64(len(extraBytes)))
	if _, err := stream.List(); err != nil {
		return nil, err
	}
	if err := stream.Decode(&extra.Vanity); err != nil {
		return nil, err
	}
	if err := stream.Decode(&extra.Validators); err != nil {
		return nil, err
	}
	if err := stream.Decode(&extra.Vote); err != nil {
		return nil, err
	}
	if err := stream.Decode(&extra.Round); err != nil {
		return nil, err
	}
	if err := stream.Decode(&extra.Seals); err != nil {
		return nil, err
	}
	if err := stream.ListEnd(); err != nil {
		return nil, err
	}

	return &extra, nil
}
