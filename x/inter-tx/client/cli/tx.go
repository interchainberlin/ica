package cli

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/interchainberlin/ica/x/inter-tx/types"
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
		Use: "register",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			sourcePort := viper.GetString(FlagSourcePort)
			sourceChannel := viper.GetString(FlagSourceChannel)

			msg := types.NewMsgRegisterAccount(
				sourcePort,
				sourceChannel,
				clientCtx.GetFromAddress().String(),
			)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	cmd.Flags().AddFlagSet(fsSourcePort)
	cmd.Flags().AddFlagSet(fsSourceChannel)

	_ = cmd.MarkFlagRequired(FlagSourcePort)
	_ = cmd.MarkFlagRequired(FlagSourceChannel)
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func getSendTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "send [type] [to_address] [amount]",
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			chainType := args[0]
			fromAddress := clientCtx.GetFromAddress()
			toAddress, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			amount, err := sdk.ParseCoinsNormalized(args[2])
			if err != nil {
				return err
			}

			sourcePort := viper.GetString(FlagSourcePort)
			sourceChannel := viper.GetString(FlagSourceChannel)

			msg := types.NewMsgSend(
				chainType,
				sourcePort,
				sourceChannel,
				fromAddress,
				toAddress,
				amount,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}
	cmd.Flags().AddFlagSet(fsSourcePort)
	cmd.Flags().AddFlagSet(fsSourceChannel)

	_ = cmd.MarkFlagRequired(FlagSourcePort)
	_ = cmd.MarkFlagRequired(FlagSourceChannel)

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
