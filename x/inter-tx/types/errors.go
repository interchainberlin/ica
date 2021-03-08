package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrIAAccountAlreadyExist = sdkerrors.Register(ModuleName, 2, "interchain account already registered")
	ErrIAAccountNotExist     = sdkerrors.Register(ModuleName, 3, "interchain account not exist")
)
