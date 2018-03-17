package cmd

import (
	"github.com/spf13/cobra"
	"github.com/thomas-fossati/trafic/runner"
)

var serversCmd = &cobra.Command{
	Use:   "servers",
	Short: "Run the server side of a traffic mix",
	Run:   servers,
}

func init() {
	rootCmd.AddCommand(serversCmd)
}

func servers(cmd *cobra.Command, args []string) {
	run(runner.RoleServer)
}
