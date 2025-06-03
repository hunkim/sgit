package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show sgit version information",
	Long:  `Display the current version of sgit.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("sgit version %s\n", version)
		fmt.Println("Solar LLM-powered git wrapper")
		fmt.Println("https://github.com/hunkim/sgit")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
} 