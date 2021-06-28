package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/interchain-accounts/x/inter-tx/types"
)

// RegisterIBCAccount uses the ibc-account module keeper to register an account on a target chain
// An address registration queue is used to keep track of registration requests
func (keeper Keeper) RegisterIBCAccount(
	ctx sdk.Context,
	sender sdk.AccAddress,
	sourcePort,
	sourceChannel string,
) error {
	salt := keeper.GetIncrementalSalt(ctx)
	err := keeper.iaKeeper.TryRegisterIBCAccount(ctx, sourcePort, sourceChannel, []byte(salt))
	if err != nil {
		return err
	}

	keeper.PushAddressToRegistrationQueue(ctx, sourcePort, sourceChannel, sender)

	ctx.EventManager().EmitEvent(sdk.NewEvent("register-interchain-account",
		sdk.NewAttribute("salt", salt)))

	return nil
}

// GetIBCAccount returns an interchain account address
func (keeper Keeper) GetIBCAccount(ctx sdk.Context, sourcePort, sourceChannel string, address sdk.AccAddress) (types.QueryIBCAccountFromAddressResponse, error) {
	store := ctx.KVStore(keeper.storeKey)

	key := types.KeyRegisteredAccount(sourcePort, sourceChannel, address)
	if !store.Has(key) {
		return types.QueryIBCAccountFromAddressResponse{}, types.ErrIBCAccountNotExist
	}
	res := types.QueryIBCAccountFromAddressResponse{}
	addr := store.Get(key)

	res.Address = addr

	return res, nil
}

// GetIBCAccount returns an interchain account address
func (keeper Keeper) GetIBCAccountAddr(ctx sdk.Context, sourcePort, sourceChannel string, address sdk.AccAddress) (sdk.AccAddress, error) {
	store := ctx.KVStore(keeper.storeKey)

	key := types.KeyRegisteredAccount(sourcePort, sourceChannel, address)
	if !store.Has(key) {
		return nil, types.ErrIBCAccountNotExist
	}
	addr := store.Get(key)

	return addr, nil
}

// GetIncrementalSalt increments the Salt value by 1 and returns the Salt
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
