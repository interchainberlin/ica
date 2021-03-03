package keeper

import (
	"github.com/interchainberlin/ica/x/interchainaccount/types"
)

var _ types.QueryServer = Keeper{}
