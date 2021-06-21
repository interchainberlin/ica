package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	host "github.com/cosmos/cosmos-sdk/x/ibc/core/24-host"
	"github.com/cosmos/interchain-accounts/x/ibc-account/types"
	"github.com/tendermint/tendermint/crypto/tmhash"
)

func (k Keeper) RegisterIBCAccount(ctx sdk.Context, owner, connectionId, counterPartyChannelId string) error {
	cap := k.portKeeper.BindPort(ctx, "cosmos1mjk79fjjgpplak5wq838w0yd982gzkyfrk07am")
	err := k.ClaimCapability(ctx, cap, host.PortPath("cosmos1mjk79fjjgpplak5wq838w0yd982gzkyfrk07am"))
	if err != nil {
		return err
	}

	return err
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
