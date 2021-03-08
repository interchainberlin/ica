package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	ModuleName = "intertx"

	StoreKey  = ModuleName
	RouterKey = ModuleName

	QuerierRoute = ModuleName
)

func KeyRegistrationQueue(sourcePort, sourceChannel string) []byte {
	return []byte(fmt.Sprintf("registration-queue/%s/%s", sourcePort, sourceChannel))
}

//nolint:interfacer
func KeyRegisteredAccount(sourcePort, sourceChannel string, addr sdk.AccAddress) []byte {
	return []byte(fmt.Sprintf("registered/%s/%s/%s", sourcePort, sourceChannel, addr.String()))
}
