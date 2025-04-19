package cli

import (
	"fmt"

	"github.com/groundctl/groundctl/internal/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Version(cmd *cobra.Command, args []string) error {
	if viper.GetViper().GetInt("verbose") > 0 {
		// Print verbose version info
		fmt.Printf("groundctl v%s (%s on %s)\n", version.Version, version.Commit, version.Branch)
	} else {
		// Print verbose version info
		fmt.Printf("v%s\n", version.Version)
	}
	return nil
}
