package cmd

import (
	"medflow/internal/service"

	"github.com/spf13/cobra"
)

// monitorCmd represents the monitor command
var monitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "A brief description of your command",
	RunE:  service.ClinicalMonitorService,
}

func init() {
	rootCmd.AddCommand(monitorCmd)
}
