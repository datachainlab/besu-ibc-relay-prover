package module

import (
	"log"

	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"
	"github.com/cosmos/ibc-go/v8/modules/core/exported"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/hyperledger-labs/yui-relayer/core"
)

var _ core.Header = (*Header)(nil)

func (Header) ClientType() string {
	return QBFT_CLIENT_TYPE
}

func (h *Header) GetHeight() exported.Height {
	ethHeader, err := h.decodeEthHeader()
	if err != nil {
		log.Panicf("invalid header: %v", h)
	}
	return ethHeightToPB(ethHeader.Number.Uint64())
}

func (h *Header) ValidateBasic() error {
	if _, err := h.decodeEthHeader(); err != nil {
		return err
	}
	if _, err := h.decodeAccountProof(); err != nil {
		return err
	}
	return nil
}

func (h *Header) decodeEthHeader() (*types.Header, error) {
	var ethHeader types.Header
	if err := rlp.DecodeBytes(h.BesuHeaderRlp, &ethHeader); err != nil {
		return nil, err
	}
	return &ethHeader, nil
}

func (h *Header) decodeAccountProof() ([][]byte, error) {
	var decodedProof [][][]byte
	if err := rlp.DecodeBytes(h.AccountStateProof, &decodedProof); err != nil {
		return nil, err
	}
	var accountProof [][]byte
	for i := range decodedProof {
		b, err := rlp.EncodeToBytes(decodedProof[i])
		if err != nil {
			return nil, err
		}
		accountProof = append(accountProof, b)
	}
	return accountProof, nil
}

func ethHeightToPB(height uint64) clienttypes.Height {
	return clienttypes.NewHeight(0, uint64(height))
}
