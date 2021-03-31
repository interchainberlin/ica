package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ibcacckeeper "github.com/cosmos/interchain-accounts/x/ibc-account/keeper"
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
