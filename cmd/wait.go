package cmd

import (
	"github.com/guidewire/netwait/wait"
	"github.com/spf13/cobra"
)

// waitCmd represents the wait command
var waitCmd = &cobra.Command{
	Use:   "wait",
	Short: "wait for a network resource to become available",
	Long: `
Repeatedly attempt to connect to a network resource and wait until a successful
connection has been established or a timeout has elapsed.

Examples:
wait https://github.com
wait github.com
wait --timeout 10s https://github.com
wait https://github.com https://github.com/guidewire/netwait
`,
	RunE: runWait,
	Args: cobra.MinimumNArgs(1),
}

func init() {
	rootCmd.AddCommand(waitCmd)
}

func runWait(cmd *cobra.Command, args []string) error {
	timeout, err := cmd.Flags().GetDuration("timeout")
	if err != nil {
		panic(err)
	}
	cmd.SilenceUsage = true

	waiter := wait.CompositeMultiWaiter{}
	return waiter.WaitMulti(args, timeout)
}
