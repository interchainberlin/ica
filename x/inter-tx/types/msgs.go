package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/cosmos-sdk/x/ibc/core/02-client/types"
)

const (
	TypeMsgRegisterAccount = "test-send"
)

var _ sdk.Msg = &MsgRegisterAccount{}

// NewMsgRegisterAccount creates a new MsgRegisterAccount instance
func NewMsgRegisterAccount(
	port, channel string, height clienttypes.Height, timestamp uint64, owner string,
) *MsgRegisterAccount {
	return &MsgRegisterAccount{
		SourcePort:       port,
		SourceChannel:    channel,
		TimeoutHeight:    height,
		TimeoutTimestamp: timestamp,
		Owner:            owner,
	}
}

// Route implements sdk.Msg
func (MsgRegisterAccount) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (MsgRegisterAccount) Type() string {
	return TypeMsgRegisterAccount
}

// ValidateBasic performs a basic check of the MsgRegisterAccount fields.
func (msg MsgRegisterAccount) ValidateBasic() error {
	return nil

}

func (msg MsgRegisterAccount) GetSignBytes() []byte {
	panic("IBC messages do not support amino")
}

// GetSigners implements sdk.Msg
func (msg MsgRegisterAccount) GetSigners() []sdk.AccAddress {
	accAddr, err := sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{accAddr}
}
