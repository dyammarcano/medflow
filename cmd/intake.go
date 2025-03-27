package cmd

import (
	"medflow/internal/service"

	"github.com/spf13/cobra"
)

var intakeCmd = &cobra.Command{
	Use:   "intake",
	Short: "A brief description of your command",
	RunE:  service.StartPatientIntakeService,
}

func init() {
	rootCmd.AddCommand(intakeCmd)
}
