package version

import (
	"github.com/groundctl/groundctl/internal/cli"
	"github.com/spf13/cobra"
)

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display version",
	Long:  "Display version and environment info.",
	RunE:  cli.Version,
}
