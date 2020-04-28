package main

import (
	"log"
	"time"

	"github.com/estenssoros/swoop/cmd"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
	logrus.SetLevel(logrus.TraceLevel)
}

func run() error {
	return cmd.Execute()
}

func main() {
	start := time.Now()
	if err := run(); err != nil {
		log.Fatal(err)
	}
	logrus.Infof("process took: %v", time.Since(start))
}
