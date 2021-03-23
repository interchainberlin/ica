package cli

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/interchainberlin/ica/x/inter-tx/types"

	"github.com/spf13/cobra"
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
		GetRegisterAccountCmd(),
		GetSendTxCmd(),
	)

	return cmd
}

func GetRegisterAccountCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "register [src-port] [src-channel]",
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			srcPort := args[0]
			srcChannel := args[1]

			msg := types.NewMsgRegisterAccount(
				srcPort,
				srcChannel,
				clientCtx.GetFromAddress().String(),
			)

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func GetSendTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "send [type] [to_address] [from_address] [amount] [source_port], [source_channel]",
		Args: cobra.ExactArgs(6),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			chainType := args[0]
			toAddress := args[1]
			fromAddress := args[2]
			amount, err := sdk.ParseCoinsNormalized(args[3])
			if err != nil {
				return err
			}
			srcPort := args[4]
			srcChannel := args[5]

			msg := types.NewMsgSend(
				chainType,
				srcPort,
				srcChannel,
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

	cmd.Flags().String(flagPacketTimeoutHeight, "0-1000", "Packet timeout block height. The timeout is disabled when set to 0-0.")
	cmd.Flags().Uint64(flagPacketTimeoutTimestamp, uint64((time.Duration(10) * time.Minute).Nanoseconds()), "Packet timeout timestamp in nanoseconds. Default is 10 minutes. The timeout is disabled when set to 0.")
	cmd.Flags().Bool(flagAbsoluteTimeouts, false, "Timeout flags are used as absolute timeouts.")
	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
