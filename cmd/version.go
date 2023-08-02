package main

import (
	"fmt"

	"github.com/soerenschneider/vault-ssh-cli/internal"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version and then exit",

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(internal.BuildVersion)
	},
}
