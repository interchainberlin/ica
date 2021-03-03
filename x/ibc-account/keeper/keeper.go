package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	capabilitykeeper "github.com/cosmos/cosmos-sdk/x/capability/keeper"
	capabilitytypes "github.com/cosmos/cosmos-sdk/x/capability/types"
	host "github.com/cosmos/cosmos-sdk/x/ibc/core/24-host"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/interchainberlin/ica/x/ibc-account/types"
)

func SerializeCosmosTx(cdc codec.BinaryMarshaler, registry codectypes.InterfaceRegistry) func(data interface{}) ([]byte, error) {
	return func(data interface{}) ([]byte, error) {
		msgs := make([]sdk.Msg, 0)
		switch data := data.(type) {
		case sdk.Msg:
			msgs = append(msgs, data)
		case []sdk.Msg:
			msgs = append(msgs, data...)
		default:
			return nil, types.ErrInvalidOutgoingData
		}

		msgAnys := make([]*codectypes.Any, len(msgs))

		for i, msg := range msgs {
			var err error
			msgAnys[i], err = codectypes.NewAnyWithValue(msg)
			if err != nil {
				return nil, err
			}
		}

		txBody := &types.IBCTxBody{
			Messages: msgAnys,
		}

		txRaw := &types.IBCTxRaw{
			BodyBytes: cdc.MustMarshalBinaryBare(txBody),
		}

		bz, err := cdc.MarshalBinaryBare(txRaw)
		if err != nil {
			return nil, err
		}

		return bz, nil
	}
}

// Keeper defines the IBC transfer keeper
type Keeper struct {
	storeKey sdk.StoreKey
	cdc      codec.BinaryMarshaler

	// Key can be chain type which means what blockchain framework the host chain was built on or just direct chain id.
	txEncoders map[string]types.TxEncoder

	hook types.IBCAccountHooks

	channelKeeper types.ChannelKeeper
	portKeeper    types.PortKeeper
	accountKeeper types.AccountKeeper

	scopedKeeper capabilitykeeper.ScopedKeeper

	router types.Router
	memKey sdk.StoreKey
}

// NewKeeper creates a new IBC account Keeper instance
func NewKeeper(
	memKey sdk.StoreKey,
	cdc codec.BinaryMarshaler, key sdk.StoreKey,
	txEncoders map[string]types.TxEncoder, channelKeeper types.ChannelKeeper, portKeeper types.PortKeeper,
	accountKeeper types.AccountKeeper, scopedKeeper capabilitykeeper.ScopedKeeper, router types.Router,
) Keeper {
	return Keeper{
		storeKey:      key,
		cdc:           cdc,
		txEncoders:    txEncoders,
		channelKeeper: channelKeeper,
		portKeeper:    portKeeper,
		accountKeeper: accountKeeper,
		scopedKeeper:  scopedKeeper,
		router:        router,
		memKey:        memKey,
	}
}

func (k Keeper) AddTxEncoder(typ string, txEncoder types.TxEncoder) error {
	_, ok := k.txEncoders[typ]
	if ok {
		return types.ErrTxEncoderAlreadyRegistered
	}
	k.txEncoders[typ] = txEncoder
	return nil
}

func (k Keeper) GetTxEncoder(typ string) (types.TxEncoder, bool) {
	info, ok := k.txEncoders[typ]
	return info, ok
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s-%s", host.ModuleName, types.ModuleName))
}

// IsBound checks if the interchain account module is already bound to the desired port
func (k Keeper) IsBound(ctx sdk.Context, portID string) bool {
	_, ok := k.scopedKeeper.GetCapability(ctx, host.PortPath(portID))
	return ok
}

// BindPort defines a wrapper function for the ort Keeper's function in
// order to expose it to module's InitGenesis function
func (k Keeper) BindPort(ctx sdk.Context, portID string) error {
	// Set the portID into our store so we can retrieve it later
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(types.PortKey), []byte(portID))

	cap := k.portKeeper.BindPort(ctx, portID)
	return k.ClaimCapability(ctx, cap, host.PortPath(portID))
}

// GetPort returns the portID for the ibc account module. Used in ExportGenesis
func (k Keeper) GetPort(ctx sdk.Context) string {
	store := ctx.KVStore(k.storeKey)
	return string(store.Get([]byte(types.PortKey)))
}

// ClaimCapability allows the transfer module that can claim a capability that IBC module
// passes to it
func (k Keeper) ClaimCapability(ctx sdk.Context, cap *capabilitytypes.Capability, name string) error {
	return k.scopedKeeper.ClaimCapability(ctx, cap, name)
}
