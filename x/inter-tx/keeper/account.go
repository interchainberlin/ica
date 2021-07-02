package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Calls the InitInterchainAccount fn which binds a new port for the account owner and opens a new ics27 channel
func (keeper Keeper) RegisterInterchainAccount(
	ctx sdk.Context,
	owner sdk.AccAddress,
	connectionId string,
) error {
	err := keeper.iaKeeper.InitInterchainAccount(ctx, connectionId, owner.String())
	if err != nil {
		return err
	}

	return nil
}
