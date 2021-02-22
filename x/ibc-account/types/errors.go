package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrUnknownPacketData          = sdkerrors.Register(ModuleName, 1, "unknown packet data")
	ErrAccountAlreadyExist        = sdkerrors.Register(ModuleName, 2, "account already exist")
	ErrUnsupportedChain           = sdkerrors.Register(ModuleName, 3, "unsupported chain")
	ErrInvalidOutgoingData        = sdkerrors.Register(ModuleName, 4, "invalid outgoing data")
	ErrInvalidRoute               = sdkerrors.Register(ModuleName, 5, "invalid route")
	ErrTxEncoderAlreadyRegistered = sdkerrors.Register(ModuleName, 6, "tx encoder already registered")
	ErrIBCAccountNotFound         = sdkerrors.Register(ModuleName, 7, "ibc account not found")
)
