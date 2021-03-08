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
	fmt.Println(ctx)
	acc, _ := sdk.AccAddressFromBech32(msg.Owner)
	fmt.Println(acc)

	_ = k.RegisterInterchainAccount(
		ctx,
		acc,
		msg.SourcePort,
		msg.SourceChannel,
		msg.TimeoutHeight,
		msg.TimeoutTimestamp,
	)

	return &types.MsgRegisterAccountResponse{}, nil
}
