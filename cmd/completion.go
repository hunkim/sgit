package cmd

import (
	"fmt"
	
	"github.com/spf13/cobra"
)

func init() {
	// Add custom completion for language flag
	rootCmd.RegisterFlagCompletionFunc("lang", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{
			"en\tEnglish",
			"ko\tKorean (한국어)", 
			"ja\tJapanese (日本語)",
			"zh\tChinese (中文)",
			"es\tSpanish (Español)",
			"fr\tFrench (Français)",
			"de\tGerman (Deutsch)",
		}, cobra.ShellCompDirectiveNoFileComp
	})
	
	// Add custom completion help
	if completionCmd := rootCmd.Commands(); len(completionCmd) > 0 {
		for _, cmd := range completionCmd {
			if cmd.Name() == "completion" {
				cmd.Long = `Generate the autocompletion script for sgit for the specified shell.

QUICK SETUP:
Run the setup script for automatic installation:
  ./scripts/setup-completion.sh

MANUAL SETUP:
Choose your shell and follow the instructions below.

The script will setup completion for sgit native commands:
  • add, commit, diff, log, merge, config
  • Global flags: --lang, --config  
  • Language codes: en, ko, ja, zh, es, fr, de

See each sub-command's help for details on manual installation.`
				
				// Add a pre-run hook to show helpful information
				originalRun := cmd.Run
				cmd.Run = func(cmd *cobra.Command, args []string) {
					if len(args) == 0 {
						fmt.Println("💡 TIP: Use './scripts/setup-completion.sh' for automatic setup!")
						fmt.Println("")
					}
					if originalRun != nil {
						originalRun(cmd, args)
					}
				}
				break
			}
		}
	}
} 