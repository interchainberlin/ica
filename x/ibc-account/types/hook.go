package types

import sdk "github.com/cosmos/cosmos-sdk/types"

type IBCAccountHooks interface {
	OnAccountCreated(ctx sdk.Context, sourcePort, sourceChannel string, address sdk.AccAddress)
	OnTxSucceeded(ctx sdk.Context, sourcePort, sourceChannel string, txHash []byte, txBytes []byte)
	OnTxFailed(ctx sdk.Context, sourcePort, sourceChannel string, txHash []byte, txBytes []byte)
}

type MultiIBCAccountHooks []IBCAccountHooks

var (
	_ IBCAccountHooks = MultiIBCAccountHooks{}
)

func NewMultiIBCAccountHooks(hooks ...IBCAccountHooks) MultiIBCAccountHooks {
	return hooks
}

func (h MultiIBCAccountHooks) OnAccountCreated(ctx sdk.Context, sourcePort, sourceChannel string, address sdk.AccAddress) {
	for i := range h {
		h[i].OnAccountCreated(ctx, sourcePort, sourceChannel, address)
	}
}

func (h MultiIBCAccountHooks) OnTxSucceeded(ctx sdk.Context, sourcePort, sourceChannel string, txHash []byte, txBytes []byte) {
	for i := range h {
		h[i].OnTxSucceeded(ctx, sourcePort, sourceChannel, txHash, txBytes)
	}
}

func (h MultiIBCAccountHooks) OnTxFailed(ctx sdk.Context, sourcePort, sourceChannel string, txHash []byte, txBytes []byte) {
	for i := range h {
		h[i].OnTxFailed(ctx, sourcePort, sourceChannel, txHash, txBytes)
	}
}
