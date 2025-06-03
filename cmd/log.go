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
	logAIAnalysis bool
	logTimeframe  string
)

// logCmd represents the log command
var logCmd = &cobra.Command{
	Use:   "log [options]",
	Short: "Show commit logs with optional AI analysis",
	Long: `Show commit logs with optional AI-powered analysis of development patterns.
Supports all git log options for full compatibility.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := runLog(cmd, args); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(logCmd)
	
	// AI-specific flags
	logCmd.Flags().BoolVar(&logAIAnalysis, "ai-analysis", false, "provide AI-powered analysis of commit history")
	logCmd.Flags().StringVar(&logTimeframe, "ai-timeframe", "last 20 commits", "timeframe description for AI analysis")
	
	// Standard git log flags - we'll pass these through to git
	logCmd.Flags().Bool("oneline", false, "show commits in one line")
	logCmd.Flags().StringP("pretty", "p", "", "pretty-print format")
	logCmd.Flags().String("format", "", "output format")
	logCmd.Flags().Bool("graph", false, "show commit graph")
	logCmd.Flags().Bool("decorate", false, "show ref names")
	logCmd.Flags().Bool("all", false, "show all branches")
	logCmd.Flags().StringP("since", "s", "", "show commits since date")
	logCmd.Flags().String("until", "", "show commits until date")
	logCmd.Flags().String("after", "", "show commits after date")
	logCmd.Flags().String("before", "", "show commits before date")
	logCmd.Flags().String("author", "", "filter by author")
	logCmd.Flags().String("committer", "", "filter by committer")
	logCmd.Flags().String("grep", "", "search commit messages")
	logCmd.Flags().StringP("number", "n", "", "limit number of commits")
	logCmd.Flags().String("skip", "", "skip number of commits")
	logCmd.Flags().Bool("reverse", false, "show commits in reverse order")
	logCmd.Flags().Bool("merges", false, "show only merge commits")
	logCmd.Flags().Bool("no-merges", false, "exclude merge commits")
	logCmd.Flags().Bool("first-parent", false, "follow only first parent")
	logCmd.Flags().Bool("stat", false, "show diffstat")
	logCmd.Flags().Bool("shortstat", false, "show short diffstat")
	logCmd.Flags().Bool("name-only", false, "show only file names")
	logCmd.Flags().Bool("name-status", false, "show file names and status")
	logCmd.Flags().String("abbrev-commit", "", "abbreviate commit hashes")
}

func runLog(cmd *cobra.Command, args []string) error {
	// Check if we're in a git repository
	if !isGitRepository() {
		return fmt.Errorf("not a git repository")
	}

	// If AI analysis is requested, we need to get the log first
	if logAIAnalysis {
		return runLogWithAIAnalysis(cmd, args)
	}

	// Otherwise, pass through to git
	return executeGitLogPassthrough(cmd, args)
}

func runLogWithAIAnalysis(cmd *cobra.Command, args []string) error {
	// Check configuration and setup if needed
	if err := ensureConfiguration(); err != nil {
		return err
	}

	// First, get the git log output
	logOutput, err := getGitLogOutput(cmd, args)
	if err != nil {
		return fmt.Errorf("error getting git log: %v", err)
	}

	if strings.TrimSpace(logOutput) == "" {
		fmt.Println("No commits found")
		return nil
	}

	// Show the regular log first
	fmt.Println("=== GIT LOG ===")
	fmt.Println(logOutput)
	fmt.Println()

	// Generate AI analysis
	apiKey := viper.GetString("upstage_api_key")
	modelName := viper.GetString("upstage_model_name")
	
	client := solar.NewClient(apiKey, modelName)
	
	fmt.Println("Generating AI analysis...")
	analysis, err := client.AnalyzeLog(logOutput, logTimeframe)
	if err != nil {
		return fmt.Errorf("error generating log analysis: %v", err)
	}

	fmt.Println("=== AI ANALYSIS ===")
	fmt.Println(analysis)

	return nil
}

func executeGitLogPassthrough(cobraCmd *cobra.Command, args []string) error {
	// Build git command with all flags and arguments
	gitArgs := []string{"log"}
	
	// Add all the flags that were set (excluding our custom AI flags)
	cobraCmd.Flags().Visit(func(flag *pflag.Flag) {
		flagName := flag.Name
		if flagName == "ai-analysis" || flagName == "ai-timeframe" {
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
			if flag.Shorthand != "" && len(flag.Shorthand) == 1 {
				gitArgs = append(gitArgs, "-"+flag.Shorthand, value)
			} else {
				gitArgs = append(gitArgs, "--"+flagName+"="+value)
			}
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

func getGitLogOutput(cmd *cobra.Command, args []string) (string, error) {
	// Build git command with all flags and arguments (excluding AI flags)
	gitArgs := []string{"log"}
	
	// Add all the flags that were set (excluding our custom AI flags)
	cmd.Flags().Visit(func(flag *pflag.Flag) {
		flagName := flag.Name
		if flagName == "ai-analysis" || flagName == "ai-timeframe" {
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
			if flag.Shorthand != "" && len(flag.Shorthand) == 1 {
				gitArgs = append(gitArgs, "-"+flag.Shorthand, value)
			} else {
				gitArgs = append(gitArgs, "--"+flagName+"="+value)
			}
		}
	})
	
	// Add any remaining arguments
	gitArgs = append(gitArgs, args...)
	
	// If no number limit is specified, default to last 20 commits for AI analysis
	hasNumberLimit := false
	for _, arg := range gitArgs {
		if strings.HasPrefix(arg, "-n") || strings.HasPrefix(arg, "--number") || strings.HasPrefix(arg, "-") && len(arg) > 1 && arg[1] >= '0' && arg[1] <= '9' {
			hasNumberLimit = true
			break
		}
	}
	
	if !hasNumberLimit {
		gitArgs = append(gitArgs, "-20")
	}
	
	// Execute git command and capture output
	gitCmd := exec.Command("git", gitArgs...)
	output, err := gitCmd.Output()
	if err != nil {
		return "", err
	}
	
	return string(output), nil
} 