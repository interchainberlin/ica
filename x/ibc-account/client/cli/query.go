package cli

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"

	"github.com/interchainberlin/ica/x/ibc-account/types"
)

func GetQueryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the ibc account module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(GetIBCAccountCmd())

	return cmd
}

func GetIBCAccountCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "ibcaccount [address_or_data] [port] [channel]",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			address, addrErr := sdk.AccAddressFromBech32(args[0])

			queryClient := types.NewQueryClient(clientCtx)
			if addrErr == nil {
				res, err := queryClient.IBCAccount(context.Background(), &types.QueryIBCAccountRequest{Address: address.String()})
				if err != nil {
					return err
				}

				return clientCtx.PrintProto(res.Account)
			}

			res, err := queryClient.IBCAccountFromData(context.Background(), &types.QueryIBCAccountFromDataRequest{Data: args[0], Port: args[1], Channel: args[2]})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res.Account)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
