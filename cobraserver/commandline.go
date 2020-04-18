package cobraserver

import (
	"os"

	"github.com/spf13/cobra"
)

func Commandline(commands ...*Command) (cobras []*cobra.Command) {
	for _, commandP := range commands {
		command := *commandP
		cmd := &cobra.Command{
			Use: command.Name,
			RunE: func(c *cobra.Command, args []string) error {
				return command.Handler(os.Stdout, os.Stdin, c.Flags())
			},
		}
		setupFlags(&command, cmd.Flags(), "text")

		cobras = append(cobras, cmd)
	}
	return
}
