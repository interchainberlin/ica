package keeper

import (
	"encoding/binary"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	clienttypes "github.com/cosmos/cosmos-sdk/x/ibc/core/02-client/types"
	channeltypes "github.com/cosmos/cosmos-sdk/x/ibc/core/04-channel/types"
	host "github.com/cosmos/cosmos-sdk/x/ibc/core/24-host"
	"github.com/interchainberlin/ica/x/ibc-account/types"
	"github.com/tendermint/tendermint/crypto/tmhash"
)

// TryRegisterIBCAccount try to register IBC account to source channel.
// If no source channel exists or doesn't have capability, it will return error.
// Salt is used to generate deterministic address.
func (k Keeper) TryRegisterIBCAccount(ctx sdk.Context, sourcePort, sourceChannel string, salt []byte, timeoutHeight clienttypes.Height, timeoutTimestamp uint64) error {
	sourceChannelEnd, found := k.channelKeeper.GetChannel(ctx, sourcePort, sourceChannel)
	if !found {
		return sdkerrors.Wrap(channeltypes.ErrChannelNotFound, sourceChannel)
	}

	destinationPort := sourceChannelEnd.GetCounterparty().GetPortID()
	destinationChannel := sourceChannelEnd.GetCounterparty().GetChannelID()

	channelCap, ok := k.scopedKeeper.GetCapability(ctx, host.ChannelCapabilityPath(sourcePort, sourceChannel))
	if !ok {
		return sdkerrors.Wrap(channeltypes.ErrChannelCapabilityNotFound, "module does not own channel capability")
	}

	// get the next sequence
	sequence, found := k.channelKeeper.GetNextSequenceSend(ctx, sourcePort, sourceChannel)
	if !found {
		return channeltypes.ErrSequenceSendNotFound
	}

	packetData := types.IBCAccountPacketData{
		Type: types.Type_REGISTER,
		Data: salt,
	}

	// TODO: Add timeout height and timestamp
	packet := channeltypes.NewPacket(
		packetData.GetBytes(),
		sequence,
		sourcePort,
		sourceChannel,
		destinationPort,
		destinationChannel,
		timeoutHeight,
		timeoutTimestamp,
	)

	return k.channelKeeper.SendPacket(ctx, channelCap, packet)
}

// TryRunTx try to send messages to source channel.
func (k Keeper) TryRunTx(ctx sdk.Context, sourcePort, sourceChannel, typ string, data interface{}, timeoutHeight clienttypes.Height, timeoutTimestamp uint64) ([]byte, error) {
	sourceChannelEnd, found := k.channelKeeper.GetChannel(ctx, sourcePort, sourceChannel)
	if !found {
		return []byte{}, sdkerrors.Wrap(channeltypes.ErrChannelNotFound, sourceChannel)
	}

	destinationPort := sourceChannelEnd.GetCounterparty().GetPortID()
	destinationChannel := sourceChannelEnd.GetCounterparty().GetChannelID()

	return k.createOutgoingPacket(ctx, sourcePort, sourceChannel, destinationPort, destinationChannel, typ, data, timeoutHeight, timeoutTimestamp)
}

func (k Keeper) createOutgoingPacket(
	ctx sdk.Context,
	sourcePort,
	sourceChannel,
	destinationPort,
	destinationChannel,
	typ string,
	data interface{},
	timeoutHeight clienttypes.Height,
	timeoutTimestamp uint64,
) ([]byte, error) {
	if data == nil {
		return []byte{}, types.ErrInvalidOutgoingData
	}

	txEncoder, ok := k.GetTxEncoder(typ)
	if !ok {
		return []byte{}, types.ErrUnsupportedChain
	}

	var msgs []sdk.Msg

	switch data := data.(type) {
	case []sdk.Msg:
		msgs = data
	case sdk.Msg:
		msgs = []sdk.Msg{data}
	default:
		return []byte{}, types.ErrInvalidOutgoingData
	}

	txBytes, err := txEncoder(msgs)
	if err != nil {
		return []byte{}, sdkerrors.Wrap(err, "invalid packet data or codec")
	}

	channelCap, ok := k.scopedKeeper.GetCapability(ctx, host.ChannelCapabilityPath(sourcePort, sourceChannel))
	if !ok {
		return []byte{}, sdkerrors.Wrap(channeltypes.ErrChannelCapabilityNotFound, "module does not own channel capability")
	}

	// get the next sequence
	sequence, found := k.channelKeeper.GetNextSequenceSend(ctx, sourcePort, sourceChannel)
	if !found {
		return []byte{}, channeltypes.ErrSequenceSendNotFound
	}

	packetData := types.IBCAccountPacketData{
		Type: types.Type_RUNTX,
		Data: txBytes,
	}

	// TODO: Add timeout height and timestamp
	packet := channeltypes.NewPacket(
		packetData.GetBytes(),
		sequence,
		sourcePort,
		sourceChannel,
		destinationPort,
		destinationChannel,
		timeoutHeight,
		timeoutTimestamp,
	)

	return k.ComputeVirtualTxHash(packetData.Data, packet.Sequence), k.channelKeeper.SendPacket(ctx, channelCap, packet)
}

func (k Keeper) DeserializeTx(_ sdk.Context, txBytes []byte) ([]sdk.Msg, error) {
	var txRaw types.IBCTxRaw

	err := k.cdc.UnmarshalBinaryBare(txBytes, &txRaw)
	if err != nil {
		return nil, err
	}

	var txBody types.IBCTxBody

	err = k.cdc.UnmarshalBinaryBare(txRaw.BodyBytes, &txBody)
	if err != nil {
		return nil, err
	}

	anys := txBody.Messages
	res := make([]sdk.Msg, len(anys))
	for i, any := range anys {
		var msg sdk.Msg
		err := k.cdc.UnpackAny(any, &msg)
		if err != nil {
			return nil, err
		}
		res[i] = msg
	}

	return res, nil
}

func (k Keeper) runTx(ctx sdk.Context, destPort, destChannel string, msgs []sdk.Msg) error {
	identifier := types.GetIdentifier(destPort, destChannel)
	err := k.AuthenticateTx(ctx, msgs, identifier)
	if err != nil {
		return err
	}

	for _, msg := range msgs {
		err := msg.ValidateBasic()
		if err != nil {
			return err
		}
	}

	// Use cache context.
	// Receive packet msg should succeed regardless of the result of logic.
	// But, if we just return the success even though handler is failed,
	// the leftovers of state transition in handler will remain.
	// However, this can make the unexpected error.
	// To solve this problem, use cache context instead of context,
	// and write the state transition if handler succeeds.
	cacheContext, writeFn := ctx.CacheContext()
	err = nil
	for _, msg := range msgs {
		_, msgErr := k.runMsg(cacheContext, msg)
		if msgErr != nil {
			err = msgErr
			break
		}
	}

	if err != nil {
		return err
	}

	// Write the state transitions if all handlers succeed.
	writeFn()

	return nil
}

// AuthenticateTx verifies that the messages have the right permission.
// It will check that the message's signers are the IBC account created by the right chain.
func (k Keeper) AuthenticateTx(ctx sdk.Context, msgs []sdk.Msg, identifier string) error {
	seen := map[string]bool{}
	var signers []sdk.AccAddress
	for _, msg := range msgs {
		for _, addr := range msg.GetSigners() {
			if !seen[addr.String()] {
				signers = append(signers, addr)
				seen[addr.String()] = true
			}
		}
	}

	for _, signer := range signers {
		// Check where the interchain account is made from.
		account := k.accountKeeper.GetAccount(ctx, signer)
		if account == nil {
			return sdkerrors.ErrUnauthorized
		}

		ibcAccount, ok := account.(types.IBCAccountI)
		if !ok {
			return sdkerrors.ErrUnauthorized
		}

		if types.GetIdentifier(ibcAccount.GetDestinationPort(), ibcAccount.GetDestinationChannel()) != identifier {
			return sdkerrors.ErrUnauthorized
		}
	}

	return nil
}

// RunMsg executes the message.
// It tries to get the handler from router. And, if router exites, it will perform message.
func (k Keeper) runMsg(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
	hander := k.router.Route(ctx, msg.Route())
	if hander == nil {
		return nil, types.ErrInvalidRoute
	}

	return hander(ctx, msg)
}

// Compute the virtual tx hash that is used only internally.
func (k Keeper) ComputeVirtualTxHash(txBytes []byte, seq uint64) []byte {
	bz := make([]byte, 8)
	binary.LittleEndian.PutUint64(bz, seq)
	return tmhash.SumTruncated(append(txBytes, bz...))
}

func (k Keeper) OnRecvPacket(ctx sdk.Context, packet channeltypes.Packet) error {
	var data types.IBCAccountPacketData
	// TODO: Remove the usage of global variable "ModuleCdc"
	if err := types.ModuleCdc.UnmarshalJSON(packet.GetData(), &data); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "cannot unmarshal interchain account packet data: %s", err.Error())
	}

	switch data.Type {
	case types.Type_REGISTER:
		_, err := k.registerIBCAccount(ctx, packet.SourcePort, packet.SourceChannel, packet.DestinationPort, packet.DestinationChannel, data.Data)
		if err != nil {
			return err
		}

		return nil
	case types.Type_RUNTX:
		msgs, err := k.DeserializeTx(ctx, data.Data)
		if err != nil {
			return err
		}

		err = k.runTx(ctx, packet.DestinationPort, packet.DestinationChannel, msgs)
		if err != nil {
			return err
		}

		return nil
	default:
		return types.ErrUnknownPacketData
	}
}

func (k Keeper) OnAcknowledgementPacket(ctx sdk.Context, packet channeltypes.Packet, data types.IBCAccountPacketData, ack types.IBCAccountPacketAcknowledgement) error {
	switch ack.Type {
	case types.Type_REGISTER:
		if ack.Code == 0 {
			if k.hook != nil {
				k.hook.OnAccountCreated(ctx, packet.SourcePort, packet.SourceChannel, k.GenerateAddress(types.GetIdentifier(packet.DestinationPort, packet.DestinationChannel), data.Data))
			}
		}
		return nil
	case types.Type_RUNTX:
		if ack.Code == 0 {
			if k.hook != nil {
				k.hook.OnTxSucceeded(ctx, packet.SourcePort, packet.SourceChannel, k.ComputeVirtualTxHash(data.Data, packet.Sequence), data.Data)
			}
		} else {
			if k.hook != nil {
				k.hook.OnTxFailed(ctx, packet.SourcePort, packet.SourceChannel, k.ComputeVirtualTxHash(data.Data, packet.Sequence), data.Data)
			}
		}
		return nil
	default:
		panic("unknown type of acknowledgement")
	}
}

func (k Keeper) OnTimeoutPacket(ctx sdk.Context, packet channeltypes.Packet, data types.IBCAccountPacketData) error {
	if k.hook != nil {
		k.hook.OnTxFailed(ctx, packet.SourcePort, packet.SourceChannel, k.ComputeVirtualTxHash(data.Data, packet.Sequence), data.Data)
	}

	return nil
}
