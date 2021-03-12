package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/interchainberlin/ica/x/inter-tx/types"
)

func (keeper Keeper) pushAddressToRegistrationQueue(ctx sdk.Context, sourcePort, sourceChannel string, address sdk.AccAddress) {
	store := ctx.KVStore(keeper.storeKey)

	queue := types.RegistrationQueue{
		Addresses: make([]sdk.AccAddress, 0),
	}
	bz := store.Get(types.KeyRegistrationQueue(sourcePort, sourceChannel))

	if len(bz) != 0 {
		keeper.cdc.MustUnmarshalBinaryBare(bz, &queue)
	}

	queue.Addresses = append(queue.Addresses, address)

	bz = keeper.cdc.MustMarshalBinaryBare(&queue)

	store.Set(types.KeyRegistrationQueue(sourcePort, sourceChannel), bz)
}

// Pop address from registration queue.
// If queue is empty, it returns []bytes{}.
func (keeper Keeper) popAddressFromRegistrationQueue(ctx sdk.Context, sourcePort, sourceChannel string) sdk.AccAddress {
	store := ctx.KVStore(keeper.storeKey)

	queue := types.RegistrationQueue{
		Addresses: make([]sdk.AccAddress, 0),
	}
	bz := store.Get(types.KeyRegistrationQueue(sourcePort, sourceChannel))

	if len(bz) != 0 {
		keeper.cdc.MustUnmarshalBinaryBare(bz, &queue)
	} else {
		return sdk.AccAddress{}
	}

	if len(queue.Addresses) == 0 {
		return sdk.AccAddress{}
	}

	addr := queue.Addresses[0]

	queue.Addresses = queue.Addresses[1:]

	bz = keeper.cdc.MustMarshalBinaryBare(&queue)
	store.Set(types.KeyRegistrationQueue(sourcePort, sourceChannel), bz)

	return addr
}
