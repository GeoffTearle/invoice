package cmd

import (
	"github.com/spf13/cobra"

	"github.com/maaslalani/invoice/cmd/generate"
)

func Root() *cobra.Command {
	rootCmd := cobra.Command{
		Use:   "invoice",
		Short: "Invoice generates invoices from the command line.",
		Long:  `Invoice generates invoices from the command line.`,
	}

	rootCmd.AddCommand(generate.Command())

	return &rootCmd
}
