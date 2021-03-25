package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/interchainberlin/ica/x/inter-tx/types"
)

// OnAccountCreated is a callback that is fired when an acknowledgement has been recieved from a target chain for registing an interchain account
// This callback is responsible for mapping an account on this chain to a registered interchain account
func (keeper Keeper) OnAccountCreated(ctx sdk.Context, sourcePort, sourceChannel string, address sdk.AccAddress) {
	receiver := keeper.PopAddressFromRegistrationQueue(ctx, sourcePort, sourceChannel)

	if !receiver.Empty() {
		store := ctx.KVStore(keeper.storeKey)

		key := types.KeyRegisteredAccount(sourcePort, sourceChannel, receiver)
		store.Set(key, address)
	}
}
