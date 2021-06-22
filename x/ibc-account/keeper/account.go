package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
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

// RegisterIBCAccount performs registering IBC account.
// It will generate the deterministic address by hashing {sourcePort}/{sourceChannel}/{salt}.
func (k Keeper) RegisterBestIBCAccount(ctx sdk.Context, sourcePort, sourceChannel, destPort, destChannel string) (types.IBCAccountI, error) {
	identifier := types.GetIdentifier(destPort, destChannel)
	address := k.GenerateAddress(identifier)

	account := k.accountKeeper.GetAccount(ctx, address)
	// TODO: Discuss the vulnerabilities when creating a new account only if the old account does not exist
	// Attackers can interrupt creating accounts by sending some assets before the packet is delivered.
	// So it is needed to check that the account is not created from users.
	// Returns an error only if the account was created by other chain.
	// We need to discuss how we can judge this case.
	if account != nil {
		return nil, sdkerrors.Wrap(types.ErrAccountAlreadyExist, account.String())
	}

	ibcAccount := types.NewIBCAccount(
		authtypes.NewBaseAccountWithAddress(address),
		sourcePort, sourceChannel, destPort, destChannel,
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
