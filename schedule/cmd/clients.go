package cmd

import (
	"github.com/spf13/cobra"
	"github.com/thomas-fossati/trafic/runner"
)

var clientsCmd = &cobra.Command{
	Use:   "clients",
	Short: "Run the client side of a traffic mix",
	Run:   clients,
}

func init() {
	rootCmd.AddCommand(clientsCmd)
}

func clients(cmd *cobra.Command, args []string) {
	run(runner.RoleClient)
}
