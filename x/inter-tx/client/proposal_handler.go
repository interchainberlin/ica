package client

import (
	"github.com/cosmos/cosmos-sdk/x/distribution/client/rest"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	"github.com/cosmos/interchain-accounts/x/inter-tx/client/cli"
)

var (
	SendProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitSendProposal, rest.ProposalRESTHandler)
	FundProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitFundProposal, rest.ProposalRESTHandler)
)
