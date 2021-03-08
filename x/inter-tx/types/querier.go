package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	QueryRegistered = "registered"
)

type QueryRegisteredParams struct {
	Account       sdk.AccAddress `json:"account"`
	SourcePort    string         `json:"source_port"`
	SourceChannel string         `json:"source_channel"`
}
