package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	FlagSourcePort    = "source-port"
	FlagSourceChannel = "source-channel"
)

// common flagsets to add to various functions
var (
	fsSourcePort    = flag.NewFlagSet("", flag.ContinueOnError)
	fsSourceChannel = flag.NewFlagSet("", flag.ContinueOnError)
)

func init() {
	fsSourcePort.String(FlagSourcePort, "", "Source port for ics-27 interchain account")
	fsSourceChannel.String(FlagSourceChannel, "", "Source channel for ics-27 interchain account")
}
