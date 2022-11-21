package cmd

import (
	"os"
	"time"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "netwaiter",
	Short: "wait for a network resource to become available",
	Long: `
Repeatedly attempt to connect to a network resource and wait until a successful
connection has been established or a timeout has elapsed.

Examples:
wait https://github.com
wait github.com
wait --timeout 10s https://github.com
wait https://github.com https://github.com/merusso/netwaiter`,
	Args: cobra.MinimumNArgs(1),
	RunE: runWait,
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
}
