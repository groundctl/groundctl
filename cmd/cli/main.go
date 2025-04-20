package main

import (
	"github.com/sirupsen/logrus"
)

func main() {
	if err := RootCmd.Execute(); err != nil {
		logrus.Fatalf("Error: %s", err)
	}
}
