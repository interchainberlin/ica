package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/cosmos-sdk/x/ibc/core/02-client/types"
	ibcacckeeper "github.com/interchainberlin/ica/x/ibc-account/keeper"
	"github.com/interchainberlin/ica/x/inter-tx/types"
)

type Keeper struct {
	cdc      codec.Marshaler
	storeKey sdk.StoreKey
	memKey   sdk.StoreKey

	iaKeeper ibcacckeeper.Keeper
}

func NewKeeper(cdc codec.Marshaler, storeKey sdk.StoreKey, iaKeeper ibcacckeeper.Keeper) Keeper {
	return Keeper{
		cdc:      cdc,
		storeKey: storeKey,

		iaKeeper: iaKeeper,
	}
}

func (keeper Keeper) RegisterInterchainAccount(
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

func (keeper Keeper) GetIBCAccount(ctx sdk.Context, sourcePort, sourceChannel string, address sdk.AccAddress) ([]byte, error) {
	store := ctx.KVStore(keeper.storeKey)

	key := types.KeyRegisteredAccount(sourcePort, sourceChannel, address)
	if !store.Has(key) {
		return []byte{}, types.ErrIAAccountNotExist
	}
	bz := store.Get(key)

	return bz, nil
}

//func (keeper Keeper) TrySendCoins(ctx sdk.Context, sourcePort, sourceChannel string, typ string, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error {
//	ibcAccount, err := keeper.GetIBCAccount(ctx, sourcePort, sourceChannel, fromAddr)
//	if err != nil {
//		return err
//	}
//
//	//fmt.Print(msg)
//
//	//	_, err = keeper.iaKeeper.TryRunTx(ctx, sourcePort, sourceChannel, typ, msg)
//	return err
//}

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

// Push address to registration queue.
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
