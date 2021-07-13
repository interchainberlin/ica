package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

// TrySendCoins builds a banktypes.NewMsgSend and uses the ibc-account module keeper to send the message to another chain
func (keeper Keeper) TrySendCoins(
	ctx sdk.Context,
	owner sdk.AccAddress,
	fromAddr,
	toAddr string,
	amt sdk.Coins,
	connectionId string,
) error {
	msg := &banktypes.MsgSend{FromAddress: fromAddr, ToAddress: toAddr, Amount: amt}
	_, err := keeper.iaKeeper.TrySendTx(ctx, owner, connectionId, msg)
	return err
}
