package cli

import (
	"context"
	"errors"
	"fmt"

	"github.com/bep/simplecobra"
	"github.com/groundctl/groundctl/pkg/cli/stack"
	"github.com/sirupsen/logrus"
)

func Execute(args []string) error {
	x, err := simplecobra.New(&rootCommand{})
	if err != nil {
		logrus.Fatal(err)
	}
	cd, err := x.Execute(context.Background(), args)
	if err != nil {
		if simplecobra.IsCommandError(err) {
			// Print the help, but also return the error to fail the command.
			cd.CobraCommand.Help()
			fmt.Println()
		}
	}
	return err
}

type rootCommand struct {
	ctx context.Context
}

func (c *rootCommand) Commands() []simplecobra.Commander {
	return []simplecobra.Commander{
		newVersionCommand(),
		stack.NewStackCommand(),
	}
}

func (c *rootCommand) PreRun(this, runner *simplecobra.Commandeer) error {
	return nil
}

func (c *rootCommand) Name() string {
	return "groundctl"
}

func (c *rootCommand) Run(ctx context.Context, cd *simplecobra.Commandeer, args []string) error {
	c.ctx = ctx
	cd.CobraCommand.Help()
	return nil
}

func (c *rootCommand) Init(cd *simplecobra.Commandeer) error {
	return nil
}

type lvl1Command struct {
	name   string
	isInit bool

	aliases []string

	localFlagName  string
	localFlagNameC string

	failInit             bool
	failWithCobraCommand bool
	disableSuggestions   bool

	rootCmd *rootCommand

	commands []simplecobra.Commander

	ctx context.Context
}

func (c *lvl1Command) Commands() []simplecobra.Commander {
	return c.commands
}

func (c *lvl1Command) PreRun(this, runner *simplecobra.Commandeer) error {
	if c.failInit {
		return fmt.Errorf("failInit")
	}
	c.isInit = true
	c.localFlagNameC = c.localFlagName + "_lvl1Command_compiled"
	c.rootCmd = this.Root.Command.(*rootCommand)
	return nil
}

func (c *lvl1Command) Name() string {
	return c.name
}

func (c *lvl1Command) Run(ctx context.Context, cd *simplecobra.Commandeer, args []string) error {
	c.ctx = ctx
	return nil
}

func (c *lvl1Command) Init(cd *simplecobra.Commandeer) error {
	if c.failWithCobraCommand {
		return errors.New("failWithCobraCommand")
	}
	cmd := cd.CobraCommand
	cmd.DisableSuggestions = c.disableSuggestions
	cmd.Aliases = c.aliases
	localFlags := cmd.Flags()
	localFlags.StringVar(&c.localFlagName, "localFlagName", "", "set localFlagName for lvl1Command")
	return nil
}

type lvl2Command struct {
	name          string
	isInit        bool
	localFlagName string

	ctx       context.Context
	rootCmd   *rootCommand
	parentCmd *lvl1Command
}

func (c *lvl2Command) Commands() []simplecobra.Commander {
	return nil
}

func (c *lvl2Command) PreRun(this, runner *simplecobra.Commandeer) error {
	c.isInit = true
	c.rootCmd = this.Root.Command.(*rootCommand)
	c.parentCmd = this.Parent.Command.(*lvl1Command)
	return nil
}

func (c *lvl2Command) Name() string {
	return c.name
}

func (c *lvl2Command) Run(ctx context.Context, cd *simplecobra.Commandeer, args []string) error {
	c.ctx = ctx
	return nil
}

func (c *lvl2Command) Init(cd *simplecobra.Commandeer) error {
	cmd := cd.CobraCommand
	localFlags := cmd.Flags()
	localFlags.StringVar(&c.localFlagName, "localFlagName", "", "set localFlagName for lvl2Command")
	return nil
}

// var rootCmd = &cobra.Command{
// 	Use:   "groundctl",
// 	Short: "groundctl (pronounced \"Ground Control\") is the CLI client for the groundctl server",
// 	Long:  "groundctl (pronounced \"Ground Control\") is the CLI client for the groundctl server",
// }

// func Execute() {
// 	if err := rootCmd.Execute(); err != nil {
// 		// logrus.Error(err)
// 		os.Exit(1)
// 	}
// }

// func init() {
// 	rootCmd.AddCommand(
// 		stack.StackCmd,
// 		versionCmd,
// 	)
// 	rootCmd.PersistentFlags().CountP("verbose", "v", "Enable verbose output")
// 	// Setup logrus
// 	logrus.SetFormatter(&logrus.TextFormatter{
// 		FullTimestamp:          true,
// 		TimestampFormat:        "2006-01-02 15:04:05",
// 		DisableColors:          false,
// 		DisableLevelTruncation: true,
// 		PadLevelText:           true,
// 	})
// }
