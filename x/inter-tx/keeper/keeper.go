package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	dist "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	transfer "github.com/cosmos/cosmos-sdk/x/ibc/applications/transfer/keeper"
	ibcacckeeper "github.com/cosmos/interchain-accounts/x/ibc-account/keeper"
	"github.com/cosmos/interchain-accounts/x/inter-tx/types"
)

type Keeper struct {
	cdc      codec.Marshaler
	storeKey sdk.StoreKey
	memKey   sdk.StoreKey

	AuthKeeper     types.AccountKeeper
	iaKeeper       ibcacckeeper.Keeper
	distKeeper     dist.Keeper
	TransferKeeper transfer.Keeper
}

func NewKeeper(cdc codec.Marshaler, storeKey sdk.StoreKey, iaKeeper ibcacckeeper.Keeper, distKeeper dist.Keeper, authKeeper types.AccountKeeper, transferKeeper transfer.Keeper) Keeper {
	return Keeper{
		cdc:            cdc,
		storeKey:       storeKey,
		distKeeper:     distKeeper,
		iaKeeper:       iaKeeper,
		AuthKeeper:     authKeeper,
		TransferKeeper: transferKeeper,
	}
}
