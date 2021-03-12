package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/cosmos-sdk/x/ibc/core/02-client/types"
	"github.com/interchainberlin/ica/x/inter-tx/types"
)

func (keeper Keeper) RegisterIBCAccount(
	ctx sdk.Context,
	sender sdk.AccAddress,
	sourcePort,
	sourceChannel string,
	timeoutHeight clienttypes.Height,
	timeoutTimestamp uint64,
) error {
	salt := keeper.GetIncrementalSalt(ctx)
	err := keeper.iaKeeper.TryRegisterIBCAccount(ctx, sourcePort, sourceChannel, []byte(salt), timeoutHeight, timeoutTimestamp)
	if err != nil {
		return err
	}

	keeper.pushAddressToRegistrationQueue(ctx, sourcePort, sourceChannel, sender)

	ctx.EventManager().EmitEvent(sdk.NewEvent("register-interchain-account",
		sdk.NewAttribute("salt", salt)))

	return nil
}

func (keeper Keeper) GetIBCAccount(ctx sdk.Context, sourcePort, sourceChannel string, address sdk.AccAddress) (types.QueryIBCAccountFromAddressResponse, error) {
	store := ctx.KVStore(keeper.storeKey)

	key := types.KeyRegisteredAccount(sourcePort, sourceChannel, address)
	if !store.Has(key) {
		return types.QueryIBCAccountFromAddressResponse{}, types.ErrIAAccountNotExist
	}
	res := types.QueryIBCAccountFromAddressResponse{}
	addr := store.Get(key)

	res.Address = addr

	return res, nil
}

func (keeper Keeper) GetIncrementalSalt(ctx sdk.Context) string {
	kvStore := ctx.KVStore(keeper.storeKey)

	key := []byte("salt")

	salt := types.Salt{
		Salt: 0,
	}
	if kvStore.Has(key) {
		keeper.cdc.MustUnmarshalBinaryBare(kvStore.Get(key), &salt)
		salt.Salt++
	}

	bz := keeper.cdc.MustMarshalBinaryBare(&salt)
	kvStore.Set(key, bz)

	return string(bz)
}
