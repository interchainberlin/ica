package types

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgRegisterAccount = "register"
	TypeMsgSend            = "send"
)

var _ sdk.Msg = &MsgRegisterAccount{}

// NewMsgRegisterAccount creates a new MsgRegisterAccount instance
func NewMsgRegisterAccount(
	owner,
	connectionId string,
) *MsgRegisterAccount {
	return &MsgRegisterAccount{
		Owner:        owner,
		ConnectionId: connectionId,
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
	interchainAccountAddr string, owner sdk.AccAddress, toAddress string, amount sdk.Coins, connectionId string,
) *MsgSend {
	return &MsgSend{
		InterchainAccount: interchainAccountAddr,
		Owner:             owner,
		ToAddress:         toAddress,
		Amount:            amount,
		ConnectionId:      connectionId,
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
	return []sdk.AccAddress{msg.Owner}
}

// ValidateBasic performs a basic check of the MsgRegisterAccount fields.
func (msg MsgSend) ValidateBasic() error {
	if strings.TrimSpace(msg.InterchainAccount) == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing sender address")
	}

	if strings.TrimSpace(msg.ToAddress) == "" {
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
