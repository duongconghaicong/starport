package main

import (
	"os"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	"github.com/hello/planet/app"
	"github.com/hello/planet/cmd/planetd/cmd"
)

func main() {
	rootCmd, _ := cmd.NewRootCmd()
	if err := svrcmd.Execute(rootCmd, app.DefaultNodeHome); err != nil {
		os.Exit(1)
	}
}
