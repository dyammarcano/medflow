package cmd

import (
	"medflow/internal/service"

	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "A brief description of your command",
	RunE:  service.StartOperationRequestHandler,
}

func init() {
	rootCmd.AddCommand(validateCmd)
}
