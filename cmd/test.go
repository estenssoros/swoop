package cmd

import "github.com/spf13/cobra"

var testCmd = &cobra.Command{
	Use:     "test",
	Short:   "",
	PreRunE: func(cmd *cobra.Command, args []string) error { return nil },
	RunE:    func(cmd *cobra.Command, args []string) error { return nil },
}
