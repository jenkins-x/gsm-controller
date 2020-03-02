package pkg

import (
	"github.com/spf13/cobra"
)

type options struct {
	Cmd  *cobra.Command
	Args []string
}

var (
	createLong = `
This is a hello world quickstart for writing CLI's in Go.  This quickstart will setup automatic CI and release
pipelines using Jenkins X, upon release you will get cross platform binaries uploaded as a GitHub release.
`

	createExample = `
# print hello to the terminal
hello user
`
)

// NewCmdHelloWorld creates a command object for the "hello world" command
func NewCmdHelloWorld() *cobra.Command {
	o := &options{}

	cmd := &cobra.Command{
		Use:     "hello",
		Short:   "Hello world command",
		Long:    createLong,
		Example: createExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			o.Args = args
			return o.Run()
		},
	}
	o.Cmd = cmd

	return cmd
}
