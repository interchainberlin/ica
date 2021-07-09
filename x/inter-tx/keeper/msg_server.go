package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/interchain-accounts/x/inter-tx/types"
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

	err = k.RegisterInterchainAccount(
		ctx,
		acc,
		msg.ConnectionId,
	)
	if err != nil {
		return &types.MsgRegisterAccountResponse{}, err
	}

	return &types.MsgRegisterAccountResponse{}, nil
}

// Send is used to send tokens from an interchain account to another account on a target chain
// The inter-tx module keeper uses the ibc-account module keeper to build and send an IBC packet with the RUNTX type
func (k msgServer) Send(goCtx context.Context, msg *types.MsgSend) (*types.MsgSendResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	err := k.TrySendCoins(ctx, msg.Owner, msg.InterchainAccount, msg.ToAddress, msg.Amount, msg.ConnectionId)
	if err != nil {
		return nil, err
	}

	return &types.MsgSendResponse{}, nil
}
