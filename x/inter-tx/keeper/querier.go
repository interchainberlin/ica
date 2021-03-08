package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/interchainberlin/ica/x/inter-tx/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		//case types.QueryRegistered:
		//return queryRegistered(ctx, req, k)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown %s query endpoint: %s", types.ModuleName, path[0])
		}
	}
}

/*func queryRegistered(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var params types.QueryRegisteredParams

	if err := types.ModuleCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	registered, err := k.GetIBCAccount(ctx, params.SourcePort, params.SourceChannel, params.Account)
	if err != nil {
		return nil, err
	}

	bz, err := codec.MarshalJSONIndent(types.ModuleCdc, registered)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}*/
