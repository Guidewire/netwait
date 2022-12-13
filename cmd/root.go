package cmd

import (
	"github.com/guidewire/netwait/wait"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "netwait",
	Short: "wait for a network resource to become available",
	Long: `
Repeatedly attempt to connect to a network resource and wait until a successful
connection has been established or a timeout has elapsed.

Examples:
netwait https://github.com
netwait github.com
netwait --timeout 10s https://github.com
netwait https://github.com https://github.com/guidewire/netwait`,
	Args: cobra.MinimumNArgs(1),
	RunE: runWait,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		silentEnabled, _ := cmd.Flags().GetBool("silent")
		if silentEnabled {
			wait.CurrentOutputLevel = wait.SILENT
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().DurationP("timeout", "t", 1*time.Minute,
		"timeout to abort connection attempts")
	rootCmd.PersistentFlags().BoolP("silent", "s", false,
		"do not print to standard out")
}
