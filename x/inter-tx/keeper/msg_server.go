package keeper

import (
	"context"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/interchainberlin/ica/x/inter-tx/types"
)

type msgServer struct {
	Keeper
}

func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) Register(
	goCtx context.Context,
	msg *types.MsgRegisterAccount,
) (*types.MsgRegisterAccountResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	acc, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return &types.MsgRegisterAccountResponse{}, err
	}

	// check if an account is already registered
	_, err = k.GetIBCAccount(ctx, msg.SourcePort, msg.SourceChannel, acc)
	if err == nil {
		return &types.MsgRegisterAccountResponse{}, fmt.Errorf("Interchain account is already registered for this account")
	}

	err = k.RegisterIBCAccount(
		ctx,
		acc,
		msg.SourcePort,
		msg.SourceChannel,
		msg.TimeoutHeight,
		msg.TimeoutTimestamp,
	)
	if err != nil {
		return &types.MsgRegisterAccountResponse{}, err
	}

	return &types.MsgRegisterAccountResponse{}, nil
}

func (k msgServer) Send(goCtx context.Context, msg *types.MsgSend) (*types.MsgSendResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	acc, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return &types.MsgSendResponse{}, err
	}
	acc2, err := sdk.AccAddressFromBech32(msg.ToAddress)
	if err != nil {
		return &types.MsgSendResponse{}, err
	}

	err = k.TrySendCoins(ctx, msg.SourcePort, msg.SourceChannel, msg.ChainType, acc, acc2, msg.Amount, msg.TimeoutHeight, msg.TimeoutTimestamp)
	if err != nil {
		return nil, err
	}

	return &types.MsgSendResponse{}, nil
}
