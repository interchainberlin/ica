package keeper_test

import (
	"fmt"

	clienttypes "github.com/cosmos/cosmos-sdk/x/ibc/core/02-client/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	channeltypes "github.com/cosmos/cosmos-sdk/x/ibc/core/04-channel/types"
	"github.com/cosmos/interchain-accounts/x/ibc-account/types"
)

func (suite *KeeperTestSuite) TestQueryIBCAccount() {
	var (
		req *types.QueryIBCAccountRequest
	)

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"empty request",
			func() {
				req = &types.QueryIBCAccountRequest{
					Address: "",
				}
			},
			false,
		},
		{
			"invalid bech32 address",
			func() {
				req = &types.QueryIBCAccountRequest{
					Address: "cosmos1ntck6f6534u630q87jpamettes6shwgddag761",
				}
			},
			false,
		},
		{
			"unexist address",
			func() {
				req = &types.QueryIBCAccountRequest{
					Address: "cosmos1ntck6f6534u630q87jpamettes6shwgddag769",
				}
			},
			false,
		},
		{
			"success",
			func() {
				packetData := types.IBCAccountPacketData{
					Type: types.Type_REGISTER,
					Data: []byte{},
				}

				packet := channeltypes.Packet{
					Sequence:           0,
					SourcePort:         "sp",
					SourceChannel:      "sc",
					DestinationPort:    "dp",
					DestinationChannel: "dc",
					Data:               packetData.GetBytes(),
					TimeoutHeight:      clienttypes.Height{},
					TimeoutTimestamp:   0,
				}

				err := suite.chainA.App.IbcAccountKeeper.OnRecvPacket(suite.chainA.GetContext(), packet)
				if err != nil {
					panic(err)
				}

				address := suite.chainA.App.IbcAccountKeeper.GenerateAddress(types.GetIdentifier("dp", "dc"), []byte{})

				req = &types.QueryIBCAccountRequest{
					Address: sdk.AccAddress(address).String(),
				}
			},
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			suite.SetupTest() // reset

			tc.malleate()
			ctx := sdk.WrapSDKContext(suite.chainA.GetContext())

			res, err := suite.queryClient.IBCAccount(ctx, req)

			if tc.expPass {
				suite.Require().NoError(err)
				suite.Require().NotNil(res)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}
