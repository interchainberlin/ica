package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	ibcacckeeper "github.com/interchainberlin/ica/x/ibc-account/keeper"
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

func (keeper Keeper) TrySendCoins(ctx sdk.Context,
	sourcePort,
	sourceChannel string,
	typ string,
	fromAddr sdk.AccAddress,
	toAddr sdk.AccAddress,
	amt sdk.Coins,
) error {
	ibcAccount, err := keeper.GetIBCAccount(ctx, sourcePort, sourceChannel, fromAddr)
	if err != nil {
		return err
	}

	msg := banktypes.NewMsgSend(ibcAccount.Address, toAddr, amt)

	_, err = keeper.iaKeeper.TryRunTx(ctx, sourcePort, sourceChannel, typ, msg)
	return err
}
