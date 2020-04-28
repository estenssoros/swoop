package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
	logrus.SetLevel(logrus.TraceLevel)
}

func init() {
	Cmd.PersistentFlags()
	Cmd.AddCommand(runCmd)
	Cmd.AddCommand(testCmd)
}

var Cmd = &cobra.Command{
	Use:   "swoop",
	Short: "moves data from one source to another",
}

func Execute() error {
	return Cmd.Execute()
}
