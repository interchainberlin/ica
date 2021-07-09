package cli

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/interchain-accounts/x/inter-tx/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagPacketTimeoutHeight    = "packet-timeout-height"
	flagPacketTimeoutTimestamp = "packet-timeout-timestamp"
	flagAbsoluteTimeouts       = "absolute-timeouts"
)

func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	// this line is used by starport scaffolding # 1
	cmd.AddCommand(
		getRegisterAccountCmd(),
		getSendTxCmd(),
	)

	return cmd
}

func getRegisterAccountCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "register --connection-id",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			connectionId := viper.GetString(FlagConnectionId)

			msg := types.NewMsgRegisterAccount(
				clientCtx.GetFromAddress().String(),
				connectionId,
			)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(fsConnectionId)
	_ = cmd.MarkFlagRequired(FlagConnectionId)

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func getSendTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "send [interchain_account_address] [to_address] [amount] --connection-id",
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			ownerAddr := clientCtx.GetFromAddress()
			interchainAccountAddr := args[0]
			toAddress := args[1]
			if err != nil {
				return err
			}

			amount, err := sdk.ParseCoinsNormalized(args[2])
			if err != nil {
				return err
			}

			connectionId := viper.GetString(FlagConnectionId)

			msg := types.NewMsgSend(
				interchainAccountAddr,
				ownerAddr,
				toAddress,
				amount,
				connectionId,
			)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().AddFlagSet(fsConnectionId)

	_ = cmd.MarkFlagRequired(FlagConnectionId)

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
