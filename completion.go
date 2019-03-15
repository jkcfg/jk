package main

import (
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:       "completion [shell]",
	Args:      cobra.OnlyValidArgs,
	ValidArgs: []string{"bash", "zsh"},
	Short:     "Generate shell completions",
	Long: `To load completion run

. <(jk completion)

To configure your bash shell to load completions for each session add to your bashrc

# ~/.bashrc or ~/.profile
. <(jk completion)

If you want to use zsh instead, do the following:

$ jk completion zsh > _jk
Then move _jk into $fpath and run compinit.
`,
	Run: func(cmd *cobra.Command, args []string) {
		shell := "bash"
		if len(args) > 0 {
			shell = args[0]
		}

		switch shell {
		case "zsh":
			jk.GenZshCompletion(os.Stdout)
		default:
			jk.GenBashCompletion(os.Stdout)
		}
	},
}

func init() {
	jk.AddCommand(completionCmd)
}
