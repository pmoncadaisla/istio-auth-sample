package cmd

import (
	"github.com/pmoncadaisla/istio-auth-sample/authn/pkg/configuration"
	"github.com/pmoncadaisla/istio-auth-sample/authn/pkg/oauth2"
	"github.com/spf13/cobra"
)

var (
	authnCmd = &cobra.Command{
		Use: "authn",
		Run: authnHandler,
	}
)

func init() {
	MainCmd.AddCommand(authnCmd)
}

func authnHandler(cmd *cobra.Command, args []string) {

	authnInstance := oauth2.New(configuration.Instance.ContextName, configuration.Instance.ServerAddress)
	authnInstance.Run()
}
