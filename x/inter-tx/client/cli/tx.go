package cli

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
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
		GetCmdSubmitSendProposal(),
		GetCmdSubmitFundProposal(),
		getRegisterCommunityAccountCmd(),
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

func getRegisterCommunityAccountCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "register-community-pool",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			sourcePort := viper.GetString(FlagSourcePort)
			sourceChannel := viper.GetString(FlagSourceChannel)

			msg := types.NewMsgRegisterCommunityAccount(
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
		Use:  "send [type] [to_address] [amount] [coin]",
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

			coin := args[3]
			sourcePort := viper.GetString(FlagSourcePort)
			sourceChannel := viper.GetString(FlagSourceChannel)

			msg := types.NewMsgSend(
				chainType,
				sourcePort,
				sourceChannel,
				fromAddress,
				toAddress,
				amount,
				coin,
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

const SEND_HELPER_TEXT = `Submit a community pool send proposal along with an initial deposit.
The proposal details must be supplied via a JSON file.
Example:
$ %s tx gov submit-proposal community-pool-send <path/to/proposal.json> --from=<key_or_address>
Where proposal.json contains:
{
  "title": "I",
  "description": "Send tokens to account on host chain",
  "sourcePort": "ibcaccount",
  "sourceChannel": "channel-0",
  "deposit": "1000stake",
	"toAddress": "cosmos1mjk79fjjgpplak5wq838w0yd982gzkyfrk07am",
	"amount": "5",
	"coin": "stake"
}
`

func GetCmdSubmitSendProposal() *cobra.Command {
	bech32PrefixAccAddr := sdk.GetConfig().GetBech32AccountAddrPrefix()

	cmd := &cobra.Command{
		Use:   "community-pool-send [proposal-file]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit a community pool send proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(SEND_HELPER_TEXT, version.AppName, bech32PrefixAccAddr),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			proposal, err := ParseCommunityPoolSendProposalWithDeposit(clientCtx.JSONMarshaler, args[0])
			if err != nil {
				return err
			}

			from := clientCtx.GetFromAddress()
			deposit, err := sdk.ParseCoinsNormalized("50stake")

			content := types.NewMsgSendProposal(proposal.Title, proposal.Description, "ibcaccount", proposal.SourceChannel, proposal.ToAddress, proposal.Amount, proposal.Coin)

			msg, err := govtypes.NewMsgSubmitProposal(content, deposit, from)
			if err != nil {
				return err
			}
			if err != nil {
				return err
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	return cmd
}

func ParseCommunityPoolSendProposalWithDeposit(cdc codec.JSONMarshaler, proposalFile string) (types.MsgSendProposal, error) {
	proposal := types.MsgSendProposal{}

	contents, err := ioutil.ReadFile(proposalFile)
	if err != nil {
		return proposal, err
	}

	if err = cdc.UnmarshalJSON(contents, &proposal); err != nil {
		return proposal, err
	}

	return proposal, nil
}

const FUND_HELPER_TEXT = `Submit a community pool fund proposal along with an initial deposit.
The proposal details must be supplied via a JSON file.
Example:
$ %s tx gov submit-proposal community-pool-fund <path/to/proposal.json> --from=<key_or_address>
Where proposal.json contains:
{
  "title": "I",
  "description": "Funding community pool interchain account",
  "sourceChannel": "channel-1",
	"coin": "stake",
}
`

func GetCmdSubmitFundProposal() *cobra.Command {
	bech32PrefixAccAddr := sdk.GetConfig().GetBech32AccountAddrPrefix()

	cmd := &cobra.Command{
		Use:   "community-pool-fund [proposal-file]",
		Args:  cobra.ExactArgs(1),
		Short: "Submit a community pool fund proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(FUND_HELPER_TEXT, version.AppName, bech32PrefixAccAddr),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			proposal, err := ParseCommunityPoolFundProposalWithDeposit(clientCtx.JSONMarshaler, args[0])
			if err != nil {
				return err
			}

			from := clientCtx.GetFromAddress()
			deposit, err := sdk.ParseCoinsNormalized("1000000000stake")

			content := types.NewMsgFundProposal(proposal.Title, proposal.Description, proposal.SourceChannel, proposal.Coin)

			msg, err := govtypes.NewMsgSubmitProposal(content, deposit, from)
			if err != nil {
				return err
			}
			if err != nil {
				return err
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	return cmd
}

func ParseCommunityPoolFundProposalWithDeposit(cdc codec.JSONMarshaler, proposalFile string) (types.MsgFundProposal, error) {
	proposal := types.MsgFundProposal{}

	contents, err := ioutil.ReadFile(proposalFile)
	if err != nil {
		return proposal, err
	}

	if err = cdc.UnmarshalJSON(contents, &proposal); err != nil {
		return proposal, err
	}

	return proposal, nil
}
