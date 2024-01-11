/*
 */
package cmd

import (
	"io"
	"os"

	"github.com/spf13/cobra"
)

// completionCmd represents the completion command
var completionCmd = &cobra.Command{
	Use:   "completion",
	Short: "Generate bash completion for the command",
	Long: `To load completions, run:
source < (demo-quickstart completion)

To load completions automaticallon on login, add this line to your .bashrc file:
source < (demo-quickstart completion)
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return completionAction(os.Stdout)
	},
}

func init() {
	RootCmd.AddCommand(completionCmd)
}

func completionAction(out io.Writer) error {
	return RootCmd.GenBashCompletion(out)
}
