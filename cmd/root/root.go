package root

import (
	"fmt"
	"os"

	"github.com/jenkins-x-labs/gsm-controller/pkg"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "gsm",
		Short: "Synchronise values from Google Secret Manager into Kubernetes secrets",
		Long: `Use this command to either watch kubernetes secrets or execute once to update kubernetes secrets
`,
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {

	rootCmd.AddCommand(pkg.NewCmdWatch())
	rootCmd.AddCommand(pkg.NewCmdList())
	rootCmd.AddCommand(pkg.NewCmdSubscribe())

}
