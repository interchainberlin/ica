package keeper_test

import (
	clienttypes "github.com/cosmos/cosmos-sdk/x/ibc/core/02-client/types"
	channeltypes "github.com/cosmos/cosmos-sdk/x/ibc/core/04-channel/types"
	host "github.com/cosmos/cosmos-sdk/x/ibc/core/24-host"
	"github.com/cosmos/ibc-go/modules/core/exported"
	ibcaccountmocktypes "github.com/cosmos/interchain-accounts/x/ibc-account/mocks/mock-inter-tx/types"
	"github.com/cosmos/interchain-accounts/x/ibc-account/types"
)

func (suite *KeeperTestSuite) TestTryRegisterIBCAccount() {
	// init connection and channel between chain A and chain B
	_, _, connA, connB := suite.coordinator.SetupClientConnections(suite.chainA, suite.chainB, exported.Tendermint)
	channelA, channelB := suite.coordinator.CreateIBCAccountChannels(suite.chainA, suite.chainB, connA, connB, channeltypes.ORDERED)

	// assume that chain A try to register IBC Account to chain B
	msg := &ibcaccountmocktypes.MsgTryRegisterIBCAccount{
		SourcePort:    channelA.PortID,
		SourceChannel: channelA.ID,
		Salt:          []byte("test"),
		TimeoutHeight: clienttypes.Height{
			RevisionNumber: 0,
			RevisionHeight: 100,
		},
		TimeoutTimestamp: 0,
		Sender:           suite.chainA.SenderAccount.GetAddress(),
	}
	err := suite.coordinator.SendMsg(suite.chainA, suite.chainB, channelB.ClientID, msg)
	suite.Require().NoError(err) // message committed

	packetData := types.IBCAccountPacketData{
		Type: types.Type_REGISTER,
		Data: []byte("test"),
	}
	packet := channeltypes.NewPacket(packetData.GetBytes(), 1, channelA.PortID, channelA.ID, channelB.PortID, channelB.ID, clienttypes.NewHeight(0, 100), 0)

	// get proof of packet commitment from chainA
	packetKey := host.PacketCommitmentKey(packet.GetSourcePort(), packet.GetSourceChannel(), packet.GetSequence())
	proof, proofHeight := suite.chainA.QueryProof(packetKey)

	nextAccNum := len(suite.chainB.App.AccountKeeper.GetAllAccounts(suite.chainB.GetContext()))

	recvMsg := channeltypes.NewMsgRecvPacket(packet, proof, proofHeight, suite.chainB.SenderAccount.GetAddress())
	err = suite.coordinator.SendMsg(suite.chainB, suite.chainA, channelA.ClientID, recvMsg)
	suite.Require().NoError(err) // message committed

	acc := suite.chainB.App.AccountKeeper.GetAccount(suite.chainB.GetContext(), suite.chainB.App.IbcAccountKeeper.GenerateAddress(types.GetIdentifier(channelB.PortID, channelB.ID), []byte("test")))
	suite.Require().NotNil(acc)

	suite.Require().Equal(nextAccNum, int(acc.GetAccountNumber()))
}

/*func (suite *KeeperTestSuite) TestRunTx() {
	// init connection and channel between chain A and chain B
	_, _, connA, connB := suite.coordinator.SetupClientConnections(suite.chainA, suite.chainB, exported.Tendermint)
	channelA, channelB := suite.coordinator.CreateIBCAccountChannels(suite.chainA, suite.chainB, connA, connB, channeltypes.ORDERED)

	func() {
		// assume that chain A try to register IBC Account to chain B
		msg := &ibcaccountmocktypes.MsgTryRegisterIBCAccount{
			SourcePort:    channelA.PortID,
			SourceChannel: channelA.ID,
			Salt:          []byte("test"),
			TimeoutHeight: clienttypes.Height{
				RevisionNumber: 0,
				RevisionHeight: 100,
			},
			TimeoutTimestamp: 0,
			Sender:           suite.chainA.SenderAccount.GetAddress(),
		}
		err := suite.coordinator.SendMsg(suite.chainA, suite.chainB, channelB.ClientID, msg)
		suite.Require().NoError(err) // message committed

		packetData := types.IBCAccountPacketData{
			Type: types.Type_REGISTER,
			Data: []byte("test"),
		}
		packet := channeltypes.NewPacket(packetData.GetBytes(), 1, channelA.PortID, channelA.ID, channelB.PortID, channelB.ID, clienttypes.NewHeight(0, 100), 0)

		// get proof of packet commitment from chainA
		packetKey := host.PacketCommitmentKey(packet.GetSourcePort(), packet.GetSourceChannel(), packet.GetSequence())
		proof, proofHeight := suite.chainA.QueryProof(packetKey)

		recvMsg := channeltypes.NewMsgRecvPacket(packet, proof, proofHeight, suite.chainB.SenderAccount.GetAddress())
		err = suite.coordinator.SendMsg(suite.chainB, suite.chainA, channelA.ClientID, recvMsg)
		suite.Require().NoError(err) // message committed

		acc := suite.chainB.App.AccountKeeper.GetAccount(suite.chainB.GetContext(), suite.chainB.App.IbcAccountKeeper.GenerateAddress(types.GetIdentifier(channelB.PortID, channelB.ID), []byte("test")))
		suite.Require().NotNil(acc)
	}()

	acc := suite.chainB.App.AccountKeeper.GetAccount(suite.chainB.GetContext(), suite.chainB.App.IbcAccountKeeper.GenerateAddress(types.GetIdentifier(channelB.PortID, channelB.ID), []byte("test")))

	// mint the token to IBC account on Chain B
	err := suite.chainB.App.BankKeeper.AddCoins(suite.chainB.GetContext(), acc.GetAddress(), sdk.Coins{sdk.Coin{
		Denom:  "test",
		Amount: sdk.NewInt(10000),
	}})
	suite.Require().NoError(err)

	toAddress, err := sdk.AccAddressFromHex("0000000000000000000000000000000000000000")
	suite.Require().NoError(err)
	// try to run tx that sends the token from IBC account to other account
	msg := &ibcaccountmocktypes.MsgTryRunTxMsgSend{
		SourcePort:    channelA.PortID,
		SourceChannel: channelA.ID,
		TimeoutHeight: clienttypes.Height{
			RevisionNumber: 0,
			RevisionHeight: 100,
		},
		TimeoutTimestamp: 0,
		FromAddress:      acc.GetAddress(),
		ToAddress:        toAddress,
		Amount: sdk.Coins{sdk.Coin{
			Denom:  "test",
			Amount: sdk.NewInt(5000),
		}},
		Sender: suite.chainA.SenderAccount.GetAddress(),
	}
	err = suite.coordinator.SendMsg(suite.chainA, suite.chainB, channelB.ClientID, msg)
	suite.Require().NoError(err) // message committed

	txBytes, err := ibcaccountkeeper.SerializeCosmosTx(suite.chainB.App.AppCodec(), suite.chainB.App.InterfaceRegistry())([]sdk.Msg{
		banktypes.NewMsgSend(acc.GetAddress(), toAddress, sdk.Coins{sdk.Coin{
			Denom:  "test",
			Amount: sdk.NewInt(5000),
		}}),
	})
	packetData := types.IBCAccountPacketData{
		Type: types.Type_RUNTX,
		Data: txBytes,
	}
	packet := channeltypes.NewPacket(packetData.GetBytes(), 2, channelA.PortID, channelA.ID, channelB.PortID, channelB.ID, clienttypes.NewHeight(0, 100), 0)

	// get proof of packet commitment from chainA
	packetKey := host.PacketCommitmentKey(packet.GetSourcePort(), packet.GetSourceChannel(), packet.GetSequence())
	proof, proofHeight := suite.chainA.QueryProof(packetKey)

	recvMsg := channeltypes.NewMsgRecvPacket(packet, proof, proofHeight, suite.chainB.SenderAccount.GetAddress())
	err = suite.coordinator.SendMsg(suite.chainB, suite.chainA, channelA.ClientID, recvMsg)
	suite.Require().NoError(err) // message committed

	// check that the balance has been transfered
	bal := suite.chainB.App.BankKeeper.GetBalance(suite.chainB.GetContext(), acc.GetAddress(), "test")
	suite.Equal("5000", bal.Amount.String())
	bal = suite.chainB.App.BankKeeper.GetBalance(suite.chainB.GetContext(), toAddress, "test")
	suite.Equal("5000", bal.Amount.String())
}
*/
