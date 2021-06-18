package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	channeltypes "github.com/cosmos/cosmos-sdk/x/ibc/core/04-channel/types"
	host "github.com/cosmos/cosmos-sdk/x/ibc/core/24-host"
	"github.com/cosmos/interchain-accounts/x/ibc-account/types"
	"github.com/tendermint/tendermint/crypto/tmhash"
)

// This function binds the port for the owner & opens a channel
// TODO:
// When an account is successfully registered you should set the active channel, do this in the onOpenTry
func (k Keeper) RegisterIBCAccount(ctx sdk.Context, owner, connectionId, counterPartyChannelId string) error {
	//TODO:
	// If the port is already bound & the account was created successfully, then exit
	cap := k.portKeeper.BindPort(ctx, owner)
	err := k.ClaimCapability(ctx, cap, host.PortPath(owner))
	if err != nil {
		return err
	}

	counterParty := channeltypes.Counterparty{PortId: "ibcaccount", ChannelId: counterPartyChannelId}
	order := channeltypes.Order(2)
	_, _, err = k.channelKeeper.ChanOpenInit(ctx, order, []string{connectionId}, owner, cap, counterParty, "ics27-1")

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			channeltypes.EventTypeChannelOpenInit,
			sdk.NewAttribute(channeltypes.AttributeKeyPortID, owner),
			sdk.NewAttribute(channeltypes.AttributeKeyChannelID, "channel-1"),
			sdk.NewAttribute(channeltypes.AttributeCounterpartyPortID, "ibcaccount"),
			sdk.NewAttribute(channeltypes.AttributeCounterpartyChannelID, counterPartyChannelId),
			sdk.NewAttribute(channeltypes.AttributeKeyConnectionID, connectionId),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, channeltypes.AttributeValueCategory),
		),
	})

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
