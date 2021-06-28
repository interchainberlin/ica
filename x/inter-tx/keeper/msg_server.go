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

// Register checks if an interchain account has account is already registered and if so returns an error.
// If no account has been registered we call RegisterIBCAccount which uses the ibc-account module keeper to send an outgoing IBC packet with a REGISTER message type.
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
		return &types.MsgRegisterAccountResponse{}, types.ErrIBCAccountAlreadyExist
	}

	err = k.RegisterIBCAccount(
		ctx,
		acc,
		msg.SourcePort,
		msg.SourceChannel,
	)

	if err != nil {
		return &types.MsgRegisterAccountResponse{}, err
	}

	return &types.MsgRegisterAccountResponse{}, nil
}

// Register checks if an interchain account has account is already registered and if so returns an error.
// If no account has been registered we call RegisterIBCAccount which uses the ibc-account module keeper to send an outgoing IBC packet with a REGISTER message type.
func (k msgServer) RegisterCommunityPool(
	goCtx context.Context,
	msg *types.MsgRegisterCommunityAccount,
) (*types.MsgRegisterCommunityAccountResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	acc := k.AuthKeeper.GetModuleAddress("distribution")
	if acc == nil {
		return &types.MsgRegisterCommunityAccountResponse{}, nil
	}

	// check if an account is already registered
	_, err := k.GetIBCAccount(ctx, msg.SourcePort, msg.SourceChannel, acc)
	if err == nil {
		return &types.MsgRegisterCommunityAccountResponse{}, types.ErrIBCAccountAlreadyExist
	}

	//TODO: no need to use salt with registering account. only one account should be allowed for the distrbution module
	err = k.RegisterIBCAccount(
		ctx,
		acc,
		msg.SourcePort,
		msg.SourceChannel,
	)
	if err != nil {
		return &types.MsgRegisterCommunityAccountResponse{}, err
	}

	return &types.MsgRegisterCommunityAccountResponse{}, nil
}

// Send is used to send tokens from an interchain account to another account on a target chain
// The inter-tx module keeper uses the ibc-account module keeper to build and send an IBC packet with the RUNTX type
func (k msgServer) Send(goCtx context.Context, msg *types.MsgSend) (*types.MsgSendResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	err := k.TrySendCoins(ctx, msg.SourcePort, msg.SourceChannel, msg.ChainType, msg.Sender, msg.ToAddress, msg.Amount)
	if err != nil {
		return nil, err
	}

	return &types.MsgSendResponse{}, nil
}
