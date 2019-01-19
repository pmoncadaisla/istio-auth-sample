package main

import (
	"github.com/pmoncadaisla/istio-auth-sample/authn/cmd"
	"github.com/pmoncadaisla/istio-auth-sample/authn/pkg/configuration"
)

func main() {
	configuration.New()
	cmd.Execute()

}
