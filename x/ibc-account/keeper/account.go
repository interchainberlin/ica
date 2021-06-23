package keeper

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	host "github.com/cosmos/cosmos-sdk/x/ibc/core/24-host"
	"github.com/cosmos/interchain-accounts/x/ibc-account/types"
	"github.com/tendermint/tendermint/crypto/tmhash"
)

// check if the port is already bound
// if so, do not bind, do not open channel
// the port should be bound with the prefix "ics27-1-"
//TODO: rename this function
func (k Keeper) RegisterIBCAccount(ctx sdk.Context, owner string) error {
	//TODO: type this
	prefix := "ics27-1-"
	portId := prefix + strings.TrimSpace(owner)

	cap := k.portKeeper.BindPort(ctx, portId)
	err := k.ClaimCapability(ctx, cap, host.PortPath(portId))
	if err != nil {
		return err
	}

	return err
}

// RegisterIBCAccount performs registering IBC account
// Here we need to register the account, set the active channel
// TODO we probably need a counter of some sort in order to check for address conflicts when creating an account
// The generated address needs to be deterministic on both sides
func (k Keeper) RegisterBestIBCAccount(ctx sdk.Context, portId string) (types.IBCAccountI, error) {
	address := k.GenerateAddress(portId)

	account := k.accountKeeper.GetAccount(ctx, address)
	if account != nil {
		return nil, sdkerrors.Wrap(types.ErrAccountAlreadyExist, account.String())
	}

	ibcAccount := types.NewIBCAccount(
		authtypes.NewBaseAccountWithAddress(address),
		portId,
	)
	k.accountKeeper.NewAccount(ctx, ibcAccount)
	k.accountKeeper.SetAccount(ctx, ibcAccount)
	return ibcAccount, nil
}

// Determine account's address that will be created.
func (k Keeper) GenerateAddress(identifier string) []byte {
	return tmhash.SumTruncated(append([]byte(identifier)))
}

func (k Keeper) GetIBCAccount(ctx sdk.Context, addr sdk.AccAddress) (types.IBCAccount, error) {
	acc := k.accountKeeper.GetAccount(ctx, addr)
	if acc == nil {
		return types.IBCAccount{}, sdkerrors.Wrap(types.ErrIBCAccountNotFound, "there is no account")
	}

	ibcAcc, ok := acc.(*types.IBCAccount)
	if !ok {
		return types.IBCAccount{}, sdkerrors.Wrap(types.ErrIBCAccountNotFound, "account is not an IBC account")
	}
	return *ibcAcc, nil
}
