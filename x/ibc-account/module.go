package ibc_account

import (
	"encoding/json"
	"fmt"

	"github.com/cosmos/interchain-accounts/x/ibc-account/client/cli"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"

	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	channeltypes "github.com/cosmos/cosmos-sdk/x/ibc/core/04-channel/types"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	porttypes "github.com/cosmos/cosmos-sdk/x/ibc/core/05-port/types"
	"github.com/cosmos/interchain-accounts/x/ibc-account/keeper"
	"github.com/cosmos/interchain-accounts/x/ibc-account/types"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
)

var (
	_ module.AppModule      = AppModule{}
	_ porttypes.IBCModule   = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

type AppModuleBasic struct{}

func (AppModuleBasic) Name() string {
	return types.ModuleName
}

// RegisterLegacyAminoCodec implements AppModuleBasic interface
func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	types.RegisterLegacyAminoCodec(cdc)
}

func (AppModuleBasic) DefaultGenesis(cdc codec.JSONMarshaler) json.RawMessage {
	return cdc.MustMarshalJSON(types.DefaultGenesis())
}

func (AppModuleBasic) ValidateGenesis(cdc codec.JSONMarshaler, config client.TxEncodingConfig, bz json.RawMessage) error {
	return nil
}

func (AppModuleBasic) RegisterRESTRoutes(ctx client.Context, rtr *mux.Router) {
	// noop
}

func (AppModuleBasic) GetTxCmd() *cobra.Command {
	// noop
	return nil
}

func (AppModuleBasic) GetQueryCmd() *cobra.Command {
	return cli.GetQueryCmd()
}

// RegisterInterfaces registers module concrete types into protobuf Any.
func (AppModuleBasic) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	types.RegisterInterfaces(registry)
}

// RegisterGRPCGatewayRoutes registers the gRPC Gateway routes for the ibc-account module.
func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
}

type AppModule struct {
	AppModuleBasic
	keeper keeper.Keeper
}

func NewAppModule(k keeper.Keeper) AppModule {
	return AppModule{
		keeper: k,
	}
}

func (AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {
	// TODO
}

func (AppModule) Route() sdk.Route {
	return sdk.NewRoute(types.RouterKey, nil)
}

func (AppModule) NewHandler() sdk.Handler {
	return nil
}

func (AppModule) QuerierRoute() string {
	return types.QuerierRoute
}

// LegacyQuerierHandler implements the AppModule interface
func (am AppModule) LegacyQuerierHandler(legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return nil
}

// RegisterServices registers a GRPC query service to respond to the
// module-specific GRPC queries.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterQueryServer(cfg.QueryServer(), am.keeper)
}

func (am AppModule) InitGenesis(ctx sdk.Context, cdc codec.JSONMarshaler, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState types.GenesisState
	cdc.MustUnmarshalJSON(data, &genesisState)

	InitGenesis(ctx, am.keeper, genesisState)
	return []abci.ValidatorUpdate{}
}

func (am AppModule) ExportGenesis(ctx sdk.Context, cdc codec.JSONMarshaler) json.RawMessage {
	gs := ExportGenesis(ctx, am.keeper)
	return cdc.MustMarshalJSON(gs)
}

// BeginBlock implements the AppModule interface
func (am AppModule) BeginBlock(ctx sdk.Context, req abci.RequestBeginBlock) {

}

// EndBlock implements the AppModule interface
func (am AppModule) EndBlock(ctx sdk.Context, req abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}

// Implement IBCModule callbacks
func (am AppModule) OnChanOpenInit(
	ctx sdk.Context,
	order channeltypes.Order,
	connectionHops []string,
	portID string,
	channelID string,
	chanCap *capabilitytypes.Capability,
	counterparty channeltypes.Counterparty,
	version string,
) error {
	return am.keeper.OnChanOpenInit(ctx, order, connectionHops, portID, channelID, chanCap, counterparty, version)
}

func (am AppModule) OnChanOpenTry(
	ctx sdk.Context,
	order channeltypes.Order,
	connectionHops []string,
	portID,
	channelID string,
	chanCap *capabilitytypes.Capability,
	counterparty channeltypes.Counterparty,
	version,
	counterpartyVersion string,
) error {
	return am.keeper.OnChanOpenTry(ctx, order, connectionHops, portID, channelID, chanCap, counterparty, version, counterpartyVersion)
}

func (am AppModule) OnChanOpenAck(
	ctx sdk.Context,
	portID,
	channelID string,
	counterpartyVersion string,
) error {
	return am.keeper.OnChanOpenAck(ctx, portID, channelID, counterpartyVersion)
}

func (am AppModule) OnChanOpenConfirm(
	ctx sdk.Context,
	portID,
	channelID string,
) error {
	return am.keeper.OnChanOpenConfirm(ctx, portID, channelID)
}

func (am AppModule) OnChanCloseInit(
	ctx sdk.Context,
	portID,
	channelID string,
) error {
	// Disallow user-initiated channel closing for interchain account channels
	return nil
}

func (am AppModule) OnChanCloseConfirm(
	ctx sdk.Context,
	portID,
	channelID string,
) error {
	return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "user cannot close channel")
}

func (am AppModule) OnRecvPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
) (*sdk.Result, []byte, error) {
	var ack channeltypes.Acknowledgement
	var data types.IBCAccountPacketData
	if err := types.ModuleCdc.UnmarshalJSON(packet.GetData(), &data); err != nil {
		return nil, nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "cannot unmarshal ICS-27 interchain account packet data: %s", err.Error())
	}

	err := am.keeper.OnRecvPacket(ctx, packet)
	if err != nil {
		ack = channeltypes.NewErrorAcknowledgement(fmt.Sprintf("cannot unmarshal interchain account packet data: %s", err.Error()))
	}

	bz := []byte("")
	ack = channeltypes.NewResultAcknowledgement(bz)

	return &sdk.Result{
		Events: ctx.EventManager().Events().ToABCIEvents(),
	}, ack.GetBytes(), nil
}

func (am AppModule) OnAcknowledgementPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	acknowledgement []byte,
) (*sdk.Result, error) {
	var ack channeltypes.Acknowledgement

	if err := types.ModuleCdc.UnmarshalJSON(acknowledgement, &ack); err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "cannot unmarshal ICS-27 interchain account packet acknowledgment: %v", err)
	}
	var data types.IBCAccountPacketData
	if err := types.ModuleCdc.UnmarshalJSON(packet.GetData(), &data); err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "cannot unmarshal ICS-27 interchain account packet data: %s", err.Error())
	}

	if err := am.keeper.OnAcknowledgementPacket(ctx, packet, data, ack); err != nil {
		return nil, err
	}

	return &sdk.Result{
		Events: ctx.EventManager().Events().ToABCIEvents(),
	}, nil
}

func (am AppModule) OnTimeoutPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
) (*sdk.Result, error) {
	// TODO
	return nil, nil
}
