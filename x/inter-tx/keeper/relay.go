package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

// TrySendCoins builds a banktypes.NewMsgSend and uses the ibc-account module keeper to send the message to another chain
func (keeper Keeper) TrySendCoins(
	ctx sdk.Context,
	fromAddr sdk.AccAddress,
	toAddr sdk.AccAddress,
	amt sdk.Coins,
	connectionId string,
) error {
	portId := keeper.iaKeeper.GeneratePortId(fromAddr.String(), connectionId)
	interchainAccountAddr, _ := keeper.iaKeeper.GetInterchainAccountAddress(ctx, portId)

	acc, _ := sdk.AccAddressFromBech32(interchainAccountAddr)
	msg := banktypes.NewMsgSend(acc, toAddr, amt)

	_, err := keeper.iaKeeper.TryRunTx(ctx, fromAddr, connectionId, msg)
	return err
}
