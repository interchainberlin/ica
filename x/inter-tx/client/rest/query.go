package rest

import (
//	"fmt"
//	"net/http"
//
//	"github.com/interchainberlin/ica/x/inter-tx/types"
//	"github.com/cosmos/cosmos-sdk/client/context"
//	"github.com/cosmos/cosmos-sdk/codec"
//	sdk "github.com/cosmos/cosmos-sdk/types"
//	"github.com/cosmos/cosmos-sdk/types/rest"
//	"github.com/gorilla/mux"
)

//func QueryRegisteredRequestHandlerFn(ctx context.CLIContext) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		w.Header().Set("Content-Type", "application/json")
//
//		vars := mux.Vars(r)
//		sourcePort := vars["source_port"]
//		sourceChannel := vars["source_channel"]
//		acc := vars["account"]
//
//		account, err := sdk.AccAddressFromBech32(acc)
//		if err != nil {
//			return
//		}
//
//		ctx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, ctx, r)
//		if !ok {
//			return
//		}
//
//		var marshaler codec.JSONMarshaler
//
//		if ctx.Marshaler != nil {
//			marshaler = ctx.Marshaler
//		} else {
//			marshaler = ctx.Codec
//		}
//
//		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryRegistered)
//		params := types.QueryRegisteredParams{Account: account, SourcePort: sourcePort, SourceChannel: sourceChannel}
//
//		bz, err := marshaler.MarshalJSON(params)
//		if rest.CheckBadRequestError(w, err) {
//			return
//		}
//
//		res, height, err := ctx.QueryWithData(route, bz)
//		if rest.CheckInternalServerError(w, err) {
//			return
//		}
//
//		ctx = ctx.WithHeight(height)
//		rest.PostProcessResponse(w, ctx, res)
//	}
//}
