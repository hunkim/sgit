package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/hunkim/sgit/pkg/solar"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	diffNoAI bool
)

// diffCmd represents the diff command
var diffCmd = &cobra.Command{
	Use:   "diff [files...]",
	Short: "Show changes with AI summary (default)",
	Long: `Show changes between commits, commit and working tree, etc. with AI-powered summaries by default.
Supports all git diff options for full compatibility. Use --no-ai to disable AI analysis.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := runDiff(cmd, args); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(diffCmd)
	
	// AI-specific flags
	diffCmd.Flags().BoolVar(&diffNoAI, "no-ai", false, "disable AI summary and use standard git diff")
	
	// Standard git diff flags - we'll pass these through to git
	diffCmd.Flags().Bool("cached", false, "show diff of staged changes")
	diffCmd.Flags().Bool("staged", false, "show diff of staged changes (same as --cached)")
	diffCmd.Flags().BoolP("patch", "p", false, "generate patch")
	diffCmd.Flags().Bool("stat", false, "show diffstat")
	diffCmd.Flags().Bool("numstat", false, "show numeric diffstat")
	diffCmd.Flags().Bool("shortstat", false, "show only summary line")
	diffCmd.Flags().Bool("name-only", false, "show only file names")
	diffCmd.Flags().Bool("name-status", false, "show only file names and status")
	diffCmd.Flags().StringP("unified", "U", "", "number of context lines")
	diffCmd.Flags().Bool("no-index", false, "compare files outside repository")
	diffCmd.Flags().Bool("ignore-space-change", false, "ignore whitespace changes")
	diffCmd.Flags().Bool("ignore-all-space", false, "ignore all whitespace")
	diffCmd.Flags().Bool("ignore-blank-lines", false, "ignore blank line changes")
	diffCmd.Flags().String("word-diff", "", "show word-level diff")
	diffCmd.Flags().String("color", "", "use colored output")
	diffCmd.Flags().Bool("no-color", false, "disable colored output")
	diffCmd.Flags().String("color-words", "", "highlight changed words")
	diffCmd.Flags().Bool("check", false, "check for whitespace errors")
	diffCmd.Flags().String("ws-error-highlight", "", "highlight whitespace errors")
}

func runDiff(cmd *cobra.Command, args []string) error {
	// Check if we're in a git repository
	if !isGitRepository() {
		return fmt.Errorf("not a git repository")
	}

	// Use AI summary by default, unless --no-ai is specified
	if !diffNoAI {
		return runDiffWithAISummary(cmd, args)
	}

	// Otherwise, pass through to git
	return executeGitDiffPassthrough(cmd, args)
}

func runDiffWithAISummary(cmd *cobra.Command, args []string) error {
	// Check configuration and setup if needed
	if err := ensureConfiguration(); err != nil {
		return err
	}

	// First, get the git diff output
	diff, err := getGitDiffOutput(cmd, args)
	if err != nil {
		return fmt.Errorf("error getting git diff: %v", err)
	}

	if strings.TrimSpace(diff) == "" {
		fmt.Println("No changes found")
		return nil
	}

	// Show the regular diff first
	fmt.Println("=== GIT DIFF ===")
	fmt.Println(diff)
	fmt.Println()

	// Generate AI summary with streaming
	apiKey := viper.GetString("upstage_api_key")
	modelName := viper.GetString("upstage_model_name")
	
	client := solar.NewClient(apiKey, modelName, getEffectiveLanguage())
	
	fmt.Println("=== AI SUMMARY ===")
	_, err = client.SummarizeDiffStream(diff)
	if err != nil {
		return fmt.Errorf("error generating diff summary: %v", err)
	}

	fmt.Println() // Add newline after streaming output
	return nil
}

func executeGitDiffPassthrough(cobraCmd *cobra.Command, args []string) error {
	// Build git command with all flags and arguments
	gitArgs := []string{"diff"}
	
	// Add all the flags that were set (excluding our custom AI flags)
	cobraCmd.Flags().Visit(func(flag *pflag.Flag) {
		flagName := flag.Name
		if flagName == "no-ai" {
			return // Skip our custom AI flags
		}
		
		value := flag.Value.String()
		if flag.Value.Type() == "bool" && value == "true" {
			if flag.Shorthand != "" && len(flag.Shorthand) == 1 {
				gitArgs = append(gitArgs, "-"+flag.Shorthand)
			} else {
				gitArgs = append(gitArgs, "--"+flagName)
			}
		} else if flag.Value.Type() != "bool" && value != "" {
			gitArgs = append(gitArgs, "--"+flagName+"="+value)
		}
	})
	
	// Add any remaining arguments
	gitArgs = append(gitArgs, args...)
	
	// Execute git command
	gitCmd := exec.Command("git", gitArgs...)
	gitCmd.Stdin = os.Stdin
	gitCmd.Stdout = os.Stdout
	gitCmd.Stderr = os.Stderr
	return gitCmd.Run()
}

func getGitDiffOutput(cmd *cobra.Command, args []string) (string, error) {
	// Build git command with all flags and arguments (excluding AI flags)
	gitArgs := []string{"diff"}
	
	// Add all the flags that were set (excluding our custom AI flags)
	cmd.Flags().Visit(func(flag *pflag.Flag) {
		flagName := flag.Name
		if flagName == "no-ai" {
			return // Skip our custom AI flags
		}
		
		value := flag.Value.String()
		if flag.Value.Type() == "bool" && value == "true" {
			if flag.Shorthand != "" && len(flag.Shorthand) == 1 {
				gitArgs = append(gitArgs, "-"+flag.Shorthand)
			} else {
				gitArgs = append(gitArgs, "--"+flagName)
			}
		} else if flag.Value.Type() != "bool" && value != "" {
			gitArgs = append(gitArgs, "--"+flagName+"="+value)
		}
	})
	
	// Add any remaining arguments
	gitArgs = append(gitArgs, args...)
	
	// Execute git command and capture output
	gitCmd := exec.Command("git", gitArgs...)
	output, err := gitCmd.Output()
	if err != nil {
		return "", err
	}
	
	return string(output), nil
} 