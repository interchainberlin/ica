package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/seantking/interchain-account/x/ibc-account/types"
)

var _ types.QueryServer = Keeper{}

// IBCAccount implements the Query/IBCAccount gRPC method
func (k Keeper) IBCAccount(ctx context.Context, req *types.QueryIBCAccountRequest) (*types.QueryIBCAccountResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.Address == "" {
		return nil, status.Error(codes.InvalidArgument, "address cannot be empty")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	address, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, err
	}

	ibcAccount, err := k.GetIBCAccount(sdkCtx, address)
	if err != nil {
		return nil, err
	}

	return &types.QueryIBCAccountResponse{Account: &ibcAccount}, nil
}

// IBCAccountFromData implements the Query/IBCAccount gRPC method
func (k Keeper) IBCAccountFromData(ctx context.Context, req *types.QueryIBCAccountFromDataRequest) (*types.QueryIBCAccountResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.Port == "" {
		return nil, status.Error(codes.InvalidArgument, "port cannot be empty")
	}

	if req.Channel == "" {
		return nil, status.Error(codes.InvalidArgument, "channel cannot be empty")
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	identifier := types.GetIdentifier(req.Port, req.Channel)
	address := k.GenerateAddress(identifier, []byte(req.Data))

	ibcAccount, err := k.GetIBCAccount(sdkCtx, address)
	if err != nil {
		return nil, err
	}

	return &types.QueryIBCAccountResponse{Account: &ibcAccount}, nil
}
