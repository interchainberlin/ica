package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

const (
	ProposalTypeFundInterchainAccount = "FundInterchainAccount"
	ProposalTypeSendInterchainAccount = "SendInterchainAccount"
)

var _ govtypes.Content = &MsgSendProposal{}

func init() {
	govtypes.RegisterProposalType(ProposalTypeSendInterchainAccount)
	govtypes.RegisterProposalType(ProposalTypeFundInterchainAccount)
	govtypes.RegisterProposalTypeCodec(&MsgSendProposal{}, "cosmos-sdk/MsgSendProposal")
	govtypes.RegisterProposalTypeCodec(&MsgFundProposal{}, "cosmos-sdk/MsgFundProposal")
}

func NewMsgSendProposal(title, description, sourcePort, sourceChannel string, toAddress sdk.AccAddress, amount sdk.Coins, coin string) *MsgSendProposal {
	return &MsgSendProposal{title, description, sourcePort, sourceChannel, toAddress, amount, coin}
}

func (csp *MsgSendProposal) ProposalRoute() string { return RouterKey }

func (csp *MsgSendProposal) ProposalType() string { return ProposalTypeSendInterchainAccount }

func (csp *MsgSendProposal) ValidateBasic() error {
	err := govtypes.ValidateAbstract(csp)
	if err != nil {
		return err
	}
	//if !csp.Amount.IsValid() {
	//	return ErrInvalidProposalAmount
	//}
	//if csp.Recipient == "" {
	//	return ErrEmptyProposalRecipient
	//}

	return nil
}

func NewMsgFundProposal(title, description, sourceChannel, coin string) *MsgFundProposal {
	return &MsgFundProposal{title, description, sourceChannel, coin}
}

func (csp *MsgFundProposal) ProposalRoute() string { return RouterKey }

func (csp *MsgFundProposal) ProposalType() string { return ProposalTypeFundInterchainAccount }

func (csp *MsgFundProposal) ValidateBasic() error {
	err := govtypes.ValidateAbstract(csp)
	if err != nil {
		return err
	}
	//if !csp.Amount.IsValid() {
	//	return ErrInvalidProposalAmount
	//}
	//if csp.Recipient == "" {
	//	return ErrEmptyProposalRecipient
	//}

	return nil
}
