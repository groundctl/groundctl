package main

import (
	"os"

	"github.com/groundctl/groundctl/pkg/cli"
	"github.com/sirupsen/logrus"
)

func main() {
	err := cli.Execute(os.Args[1:])
	if err != nil {
		logrus.Fatalf("Error: %s", err)
	}
}
