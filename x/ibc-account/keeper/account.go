package keeper

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	channeltypes "github.com/cosmos/cosmos-sdk/x/ibc/core/04-channel/types"
	host "github.com/cosmos/cosmos-sdk/x/ibc/core/24-host"
	"github.com/cosmos/interchain-accounts/x/ibc-account/types"
	"github.com/tendermint/tendermint/crypto/tmhash"
)

// The first step in registering an interchain account
// Binds a new port & calls OnChanOpenInit
func (k Keeper) InitInterchainAccount(ctx sdk.Context, owner string) error {
	portId := types.IcaPrefix + strings.TrimSpace(owner)

	// Check if the port is already bound
	isBound := k.IsBound(ctx, portId)
	if isBound == true {
		return sdkerrors.Wrap(types.ErrPortAlreadyBound, portId)
	}

	cap := k.portKeeper.BindPort(ctx, portId)
	err := k.ClaimCapability(ctx, cap, host.PortPath(portId))
	if err != nil {
		return err
	}
	connectionId := "connection-0"
	counterParty := channeltypes.Counterparty{PortId: "ibcaccount", ChannelId: ""}
	order := channeltypes.Order(2)
	_, _, err = k.channelKeeper.ChanOpenInit(ctx, order, []string{connectionId}, portId, cap, counterParty, "ics27-1")

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			channeltypes.EventTypeChannelOpenInit,
			sdk.NewAttribute(channeltypes.AttributeKeyPortID, portId),
			sdk.NewAttribute(channeltypes.AttributeKeyChannelID, "channel-0"),
			sdk.NewAttribute(channeltypes.AttributeCounterpartyPortID, "ibcaccount"),
			sdk.NewAttribute(channeltypes.AttributeCounterpartyChannelID, ""),
			sdk.NewAttribute(channeltypes.AttributeKeyConnectionID, connectionId),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, channeltypes.AttributeValueCategory),
		),
	})

	return err
}

// Register interchain account if it has not already been created
// TODO: if the address is already taken we need a mechanism to work around. The address needs to be deterministic on both sending/recieving chain
func (k Keeper) RegisterInterchainAccount(ctx sdk.Context, portId string) (types.IBCAccountI, error) {
	address := k.GenerateAddress(portId)
	account := k.accountKeeper.GetAccount(ctx, address)

	if account != nil {
		return nil, sdkerrors.Wrap(types.ErrAccountAlreadyExist, account.String())
	}

	interchainAccount := types.NewIBCAccount(
		authtypes.NewBaseAccountWithAddress(address),
		portId,
	)
	k.accountKeeper.NewAccount(ctx, interchainAccount)
	k.accountKeeper.SetAccount(ctx, interchainAccount)
	return interchainAccount, nil
}

func (k Keeper) SetActiveChannel(ctx sdk.Context, portId, channelId string) error {
	store := ctx.KVStore(k.storeKey)

	key := types.KeyActiveChannel(portId)
	store.Set(key, []byte(channelId))
	return nil
}

func (k Keeper) GetActiveChannel(ctx sdk.Context, portId string) (string, error) {
	store := ctx.KVStore(k.storeKey)
	key := types.KeyActiveChannel(portId)
	if !store.Has(key) {
		return "", sdkerrors.Wrap(types.ErrActiveChannelNotFound, portId)
	}

	activeChannel := string(store.Get(key))
	return activeChannel, nil
}

func (k Keeper) SetInterchainAccountAddress(ctx sdk.Context, portId string) sdk.AccAddress {
	store := ctx.KVStore(k.storeKey)
	address := sdk.AccAddress(k.GenerateAddress(portId))
	key := types.KeyOwnerAccount(portId)
	store.Set(key, address)
	return address
}

func (k Keeper) GetInterchainAccountAddress(ctx sdk.Context, portId string) (string, error) {
	store := ctx.KVStore(k.storeKey)
	key := types.KeyOwnerAccount(portId)
	if !store.Has(key) {
		return "", sdkerrors.Wrap(types.ErrIBCAccountNotFound, portId)
	}

	interchainAccountAddr := sdk.AccAddress(store.Get(key))
	return interchainAccountAddr.String(), nil
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
