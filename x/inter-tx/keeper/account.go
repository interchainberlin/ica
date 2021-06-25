package keeper

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	icatypes "github.com/cosmos/interchain-accounts/x/ibc-account/types"
)

// RegisterIBCAccount uses the ibc-account module keeper to register an account on a target chain
// An address registration queue is used to keep track of registration requests
func (keeper Keeper) RegisterIBCAccount(
	ctx sdk.Context,
	sender sdk.AccAddress,
) error {
	err := keeper.iaKeeper.InitInterchainAccount(ctx, sender.String())
	if err != nil {
		return err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent("register-interchain-account"))

	return nil
}

// GetIBCAccount returns an interchain account address
func (keeper Keeper) GetIBCAccount(ctx sdk.Context, owner sdk.AccAddress) (string, error) {
	portId := icatypes.IcaPrefix + strings.TrimSpace(owner.String())
	address, err := keeper.iaKeeper.GetInterchainAccountAddress(ctx, portId)

	return address, err
}
