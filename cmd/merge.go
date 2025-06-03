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
	mergeAIHelp    bool
	mergeAIMessage bool
)

// mergeCmd represents the merge command
var mergeCmd = &cobra.Command{
	Use:   "merge [branch]",
	Short: "Join development histories with optional AI assistance",
	Long: `Join two or more development histories together with optional AI assistance
for conflict resolution and merge message generation. Supports all git merge options.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := runMerge(cmd, args); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(mergeCmd)
	
	// AI-specific flags
	mergeCmd.Flags().BoolVar(&mergeAIHelp, "ai-help", false, "provide AI assistance for merge conflicts")
	mergeCmd.Flags().BoolVar(&mergeAIMessage, "ai-message", false, "generate AI merge commit message")
	
	// Standard git merge flags - we'll pass these through to git
	mergeCmd.Flags().Bool("commit", false, "perform the merge and commit the result")
	mergeCmd.Flags().Bool("no-commit", false, "perform merge but don't commit")
	mergeCmd.Flags().Bool("edit", false, "edit merge commit message")
	mergeCmd.Flags().Bool("no-edit", false, "don't edit merge commit message")
	mergeCmd.Flags().Bool("ff", false, "fast-forward if possible")
	mergeCmd.Flags().Bool("no-ff", false, "create merge commit even if fast-forward is possible")
	mergeCmd.Flags().Bool("ff-only", false, "abort if fast-forward is not possible")
	mergeCmd.Flags().Bool("log", false, "populate log message with one-line descriptions")
	mergeCmd.Flags().Bool("no-log", false, "don't populate log message with one-line descriptions")
	mergeCmd.Flags().Bool("stat", false, "show diffstat after merge")
	mergeCmd.Flags().Bool("no-stat", false, "don't show diffstat after merge")
	mergeCmd.Flags().Bool("squash", false, "create single commit instead of merge")
	mergeCmd.Flags().Bool("no-squash", false, "don't squash commits")
	mergeCmd.Flags().StringP("strategy", "s", "", "merge strategy")
	mergeCmd.Flags().StringP("strategy-option", "X", "", "merge strategy option")
	mergeCmd.Flags().Bool("verify-signatures", false, "verify signatures on merge commits")
	mergeCmd.Flags().Bool("no-verify-signatures", false, "don't verify signatures")
	mergeCmd.Flags().Bool("summary", false, "show summary of changes")
	mergeCmd.Flags().Bool("no-summary", false, "don't show summary of changes")
	mergeCmd.Flags().StringP("message", "m", "", "merge commit message")
	mergeCmd.Flags().Bool("quiet", false, "operate quietly")
	mergeCmd.Flags().Bool("verbose", false, "be verbose")
	mergeCmd.Flags().Bool("progress", false, "show progress")
	mergeCmd.Flags().Bool("no-progress", false, "don't show progress")
	mergeCmd.Flags().Bool("allow-unrelated-histories", false, "allow merging unrelated histories")
}

func runMerge(cmd *cobra.Command, args []string) error {
	// Check if we're in a git repository
	if !isGitRepository() {
		return fmt.Errorf("not a git repository")
	}

	// If AI assistance is requested, we handle it specially
	if mergeAIHelp || mergeAIMessage {
		return runMergeWithAI(cmd, args)
	}

	// Otherwise, pass through to git
	return executeGitMergePassthrough(cmd, args)
}

func runMergeWithAI(cmd *cobra.Command, args []string) error {
	// Check configuration and setup if needed
	if err := ensureConfiguration(); err != nil {
		return err
	}

	if len(args) == 0 {
		return fmt.Errorf("no branch specified for merge")
	}

	sourceBranch := args[0]
	targetBranch, _ := getCurrentBranch()

	// First, try the merge to see if there are conflicts
	fmt.Printf("Attempting to merge %s into %s...\n", sourceBranch, targetBranch)
	
	// Execute the merge with --no-commit first to check for conflicts
	mergeArgs := buildMergeArgs(cmd, args)
	mergeArgs = append(mergeArgs, "--no-commit")
	
	gitCmd := exec.Command("git", mergeArgs...)
	gitCmd.Stdin = os.Stdin
	gitCmd.Stdout = os.Stdout
	gitCmd.Stderr = os.Stderr
	
	err := gitCmd.Run()
	if err != nil {
		// Check if there are merge conflicts
		conflictFiles, conflictErr := getMergeConflicts()
		if conflictErr == nil && len(conflictFiles) > 0 {
			fmt.Println("\n🚨 Merge conflicts detected!")
			
			if mergeAIHelp {
				fmt.Println("Getting AI assistance for conflict resolution...")
				if aiErr := provideMergeConflictHelp(conflictFiles); aiErr != nil {
					fmt.Printf("Warning: Could not get AI assistance: %v\n", aiErr)
				}
			}
			
			fmt.Println("\nPlease resolve conflicts manually and then:")
			fmt.Println("  git add <resolved-files>")
			fmt.Println("  sgit merge --continue")
			return nil
		}
		return fmt.Errorf("merge failed: %v", err)
	}

	// No conflicts, proceed with commit
	if mergeAIMessage {
		return commitMergeWithAIMessage(sourceBranch, targetBranch)
	}

	// Complete the merge with regular commit
	return exec.Command("git", "commit").Run()
}

func provideMergeConflictHelp(conflictFiles []string) error {
	apiKey := viper.GetString("upstage_api_key")
	modelName := viper.GetString("upstage_model_name")
	
	client := solar.NewClient(apiKey, modelName, getEffectiveLanguage())
	
	conflictInfo := strings.Join(conflictFiles, "\n")
	
	help, err := client.AnalyzeMergeConflicts(conflictInfo)
	if err != nil {
		return err
	}

	fmt.Println("\n=== AI MERGE CONFLICT ASSISTANCE ===")
	fmt.Println(help)
	fmt.Println()

	return nil
}

func commitMergeWithAIMessage(sourceBranch, targetBranch string) error {
	// Get information about the changes being merged
	changesCmd := exec.Command("git", "log", "--oneline", "--no-merges", fmt.Sprintf("%s..%s", targetBranch, sourceBranch))
	changesOutput, err := changesCmd.Output()
	if err != nil {
		changesOutput = []byte("Unable to get merge changes")
	}

	apiKey := viper.GetString("upstage_api_key")
	modelName := viper.GetString("upstage_model_name")
	
	client := solar.NewClient(apiKey, modelName, getEffectiveLanguage())
	
	fmt.Println("Generating AI merge commit message...")
	message, err := client.GenerateMergeCommitMessage(sourceBranch, targetBranch, string(changesOutput))
	if err != nil {
		return fmt.Errorf("error generating merge message: %v", err)
	}

	fmt.Printf("Generated merge message:\n%s\n", message)

	// Complete the merge with the AI-generated message
	return exec.Command("git", "commit", "-m", message).Run()
}

func getMergeConflicts() ([]string, error) {
	cmd := exec.Command("git", "diff", "--name-only", "--diff-filter=U")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	files := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(files) == 1 && files[0] == "" {
		return []string{}, nil
	}

	return files, nil
}

func buildMergeArgs(cmd *cobra.Command, args []string) []string {
	gitArgs := []string{"merge"}
	
	// Add all the flags that were set (excluding our custom AI flags)
	cmd.Flags().Visit(func(flag *pflag.Flag) {
		flagName := flag.Name
		if flagName == "ai-help" || flagName == "ai-message" {
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
	
	return gitArgs
}

func executeGitMergePassthrough(cobraCmd *cobra.Command, args []string) error {
	gitArgs := buildMergeArgs(cobraCmd, args)
	
	// Execute git command
	gitCmd := exec.Command("git", gitArgs...)
	gitCmd.Stdin = os.Stdin
	gitCmd.Stdout = os.Stdout
	gitCmd.Stderr = os.Stderr
	return gitCmd.Run()
} 