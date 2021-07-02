package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/interchain-accounts/x/inter-tx/types"
)

// IBCAccountFromAddress implements the Query/IBCAccount gRPC method
func (k Keeper) IBCAccountFromAddress(ctx context.Context, req *types.QueryIBCAccountFromAddressRequest) (*types.QueryIBCAccountFromAddressResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	portId := k.iaKeeper.GeneratePortId(req.Address.String(), req.ConnectionId)
	addr, err := k.iaKeeper.GetInterchainAccountAddress(sdkCtx, portId)
	if err != nil {
		return nil, err
	}

	ibcAccount := types.QueryIBCAccountFromAddressResponse{Address: addr}

	return &ibcAccount, nil
}
