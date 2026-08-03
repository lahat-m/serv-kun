package cmd

import (
	"github.com/spf13/cobra"
)

var springbootLogsCmd = &cobra.Command{
	Use:   "logs [service]",
	Short: "Tail logs for a compose service (defaults to \"app\")",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runSpringbootLogs,
}

func runSpringbootLogs(cmd *cobra.Command, args []string) error {
	requireRoot(cmd)
	proj := projectFromFlags(cmd)

	service := "app"
	if len(args) == 1 {
		service = args[0]
	}

	return proj.Run("logs", "-f", "--tail=200", service)
}
