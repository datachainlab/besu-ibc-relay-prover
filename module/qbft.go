package module

import (
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/ibc-go/v8/modules/core/exported"
)

const QBFT_CLIENT_TYPE = "hb-qbft"

var _ exported.ClientState = (*ClientState)(nil)

func (cs *ClientState) ClientType() string {
	return QBFT_CLIENT_TYPE
}

func (cs *ClientState) GetLatestHeight() exported.Height {
	return cs.LatestHeight
}

func (cs *ClientState) Validate() error {
	return nil
}

func (cs *ClientState) Status(ctx sdk.Context, clientStore storetypes.KVStore, cdc codec.BinaryCodec) exported.Status {
	panic("not implemented")
}

func (cs *ClientState) ExportMetadata(clientStore storetypes.KVStore) []exported.GenesisMetadata {
	panic("not implemented")
}

func (cs *ClientState) ZeroCustomFields() exported.ClientState {
	panic("not implemented")
}

func (cs *ClientState) GetTimestampAtHeight(ctx sdk.Context, clientStore storetypes.KVStore, cdc codec.BinaryCodec, height exported.Height) (uint64, error) {
	panic("not implemented")
}

func (cs *ClientState) Initialize(ctx sdk.Context, cdc codec.BinaryCodec, clientStore storetypes.KVStore, consensusState exported.ConsensusState) error {
	panic("not implemented")
}

func (cs *ClientState) VerifyMembership(ctx sdk.Context, clientStore storetypes.KVStore, cdc codec.BinaryCodec, height exported.Height, delayTimePeriod uint64, delayBlockPeriod uint64, proof []byte, path exported.Path, value []byte) error {
	panic("not implemented")
}

func (cs *ClientState) VerifyNonMembership(ctx sdk.Context, clientStore storetypes.KVStore, cdc codec.BinaryCodec, height exported.Height, delayTimePeriod uint64, delayBlockPeriod uint64, proof []byte, path exported.Path) error {
	panic("not implemented")
}

func (cs *ClientState) VerifyClientMessage(ctx sdk.Context, cdc codec.BinaryCodec, clientStore storetypes.KVStore, clientMsg exported.ClientMessage) error {
	panic("not implemented")
}

func (cs *ClientState) CheckForMisbehaviour(ctx sdk.Context, cdc codec.BinaryCodec, clientStore storetypes.KVStore, clientMsg exported.ClientMessage) bool {
	panic("not implemented")
}

func (cs *ClientState) UpdateStateOnMisbehaviour(ctx sdk.Context, cdc codec.BinaryCodec, clientStore storetypes.KVStore, clientMsg exported.ClientMessage) {
	panic("not implemented")
}

func (cs *ClientState) UpdateState(ctx sdk.Context, cdc codec.BinaryCodec, clientStore storetypes.KVStore, clientMsg exported.ClientMessage) []exported.Height {
	panic("not implemented")
}

func (cs *ClientState) CheckSubstituteAndUpdateState(ctx sdk.Context, cdc codec.BinaryCodec, subjectClientStore, substituteClientStore storetypes.KVStore, substituteClient exported.ClientState) error {
	panic("not implemented")
}

func (cs *ClientState) VerifyUpgradeAndUpdateState(ctx sdk.Context, cdc codec.BinaryCodec, store storetypes.KVStore, newClient exported.ClientState, newConsState exported.ConsensusState, proofUpgradeClient, proofUpgradeConsState []byte) error {
	panic("not implemented")
}

var _ exported.ConsensusState = (*ConsensusState)(nil)

func (cs *ConsensusState) ClientType() string {
	panic("not implemented")
}

func (cs *ConsensusState) GetTimestamp() uint64 {
	panic("not implemented")
}

func (cs *ConsensusState) ValidateBasic() error {
	panic("not implemented")
}
