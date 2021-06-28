package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"

	// this line is used by starport scaffolding # 1
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	//	cdc.RegisterConcrete(MsgRegister{}, "intertx/MsgRegister", nil)
	//	cdc.RegisterConcrete(MsgSend{}, "intertx/MsgSend", nil)
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgSend{},
		&MsgRegisterAccount{},
		&MsgRegisterCommunityAccount{},
	)
	registry.RegisterImplementations(
		(*govtypes.Content)(nil),
		&MsgSendProposal{},
		&MsgFundProposal{},
	)

}

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)
