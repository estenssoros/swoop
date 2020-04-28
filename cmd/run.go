package cmd

import (
	"github.com/estenssoros/swoop/config"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var file string

func init() {
	runCmd.Flags().StringVarP(&file, "file", "f", "", "file to load config from")
}

var runCmd = &cobra.Command{
	Use:     "run",
	Short:   "run an etl",
	PreRunE: func(cmd *cobra.Command, args []string) error { return nil },
	RunE: func(cmd *cobra.Command, args []string) error {
		if file != "" {
			return runFile(file)
		}
		return nil
	},
}

func runFile(fileName string) error {
	cfg, err := config.NewFromFile(fileName)
	if err != nil {
		return errors.Wrap(err, "config from file")
	}
	if err := cfg.Connect(); err != nil {
		return errors.Wrap(err, "config connect")
	}
	if err := cfg.Process(); err != nil {
		return errors.Wrap(err, "process")
	}
	return nil
}
