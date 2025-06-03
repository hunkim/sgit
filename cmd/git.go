package cmd

import (
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

// gitCmd represents passthrough git commands
var gitCmd = &cobra.Command{
	Use:   "git",
	Short: "Passthrough to git commands",
	Long:  `Execute git commands directly. This is a passthrough for convenience.`,
	Run: func(cmd *cobra.Command, args []string) {
		executeGitCommand(args)
	},
	DisableFlagParsing: true,
}

func init() {
	rootCmd.AddCommand(gitCmd)
}

func executeGitCommand(args []string) {
	gitCmd := exec.Command("git", args...)
	gitCmd.Stdin = os.Stdin
	gitCmd.Stdout = os.Stdout
	gitCmd.Stderr = os.Stderr
	
	if err := gitCmd.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			os.Exit(exitError.ExitCode())
		}
		os.Exit(1)
	}
} 