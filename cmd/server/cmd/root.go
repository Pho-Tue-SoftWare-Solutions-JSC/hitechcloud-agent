package cmd

import (
	"github.com/Pho-Tue-SoftWare-Solutions-JSC/hitechcloud-agent/server"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use: "HiTechCloud-agent",
	RunE: func(cmd *cobra.Command, args []string) error {
		server.Start()
		return nil
	},
}
