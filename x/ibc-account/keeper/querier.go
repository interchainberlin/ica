package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/interchainberlin/ica/x/ibc-account/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewQuerier(k Keeper, legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryIBCAccount:
			return queryIBCAccount(ctx, req, k, legacyQuerierCdc)

		case types.QueryIBCAccountFromData:
			return queryIBCAccountFromData(ctx, req, k, legacyQuerierCdc)

		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown query path: %s", path[0])
		}
	}
}

func queryIBCAccount(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryIBCAccountRequest
	if err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	addr, err := sdk.AccAddressFromBech32(params.Address)
	if err != nil {
		return nil, err
	}

	account, err := k.GetIBCAccount(ctx, addr)
	if err != nil {
		return nil, err
	}

	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, account)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}

func queryIBCAccountFromData(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryIBCAccountFromDataRequest
	if err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	if params.Port == "" {
		return nil, status.Error(codes.InvalidArgument, "port cannot be empty")
	}

	if params.Channel == "" {
		return nil, status.Error(codes.InvalidArgument, "channel cannot be empty")
	}

	identifier := types.GetIdentifier(params.Port, params.Channel)
	address := k.GenerateAddress(identifier, []byte(params.Data))

	account, err := k.GetIBCAccount(ctx, address)
	if err != nil {
		return nil, err
	}

	bz, err := codec.MarshalJSONIndent(legacyQuerierCdc, account)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return bz, nil
}
