package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"strings"
)

const (
	TypeMsgRegisterAccount = "register"
	TypeMsgSend            = "send"
)

var _ sdk.Msg = &MsgRegisterAccount{}

// NewMsgRegisterAccount creates a new MsgRegisterAccount instance
func NewMsgRegisterAccount(
	port, channel string, owner string,
) *MsgRegisterAccount {
	return &MsgRegisterAccount{
		SourcePort:    port,
		SourceChannel: channel,
		Owner:         owner,
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

func (msg MsgRegisterAccount) ValidateBasic() error {
	if strings.TrimSpace(msg.Owner) == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing sender address")
	}

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

var _ sdk.Msg = &MsgSend{}

// NewMsgSend creates a new MsgSend instance
func NewMsgSend(
	chainType, port, channel string, sender, toAddress sdk.AccAddress, amount sdk.Coins,
) *MsgSend {
	return &MsgSend{
		ChainType:     chainType,
		SourcePort:    port,
		SourceChannel: channel,
		Sender:        sender,
		ToAddress:     toAddress,
		Amount:        amount,
	}
}

// Route implements sdk.Msg
func (MsgSend) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (MsgSend) Type() string {
	return TypeMsgSend
}

// GetSigners implements sdk.Msg
func (msg MsgSend) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

// ValidateBasic performs a basic check of the MsgRegisterAccount fields.
func (msg MsgSend) ValidateBasic() error {
	if strings.TrimSpace(msg.Sender.String()) == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing sender address")
	}

	if strings.TrimSpace(msg.ToAddress.String()) == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing recipient address")
	}

	if !msg.Amount.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, msg.Amount.String())
	}
	return nil
}

func (msg MsgSend) GetSignBytes() []byte {
	panic("IBC messages do not support amino")
}
