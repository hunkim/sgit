package cmd

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/hunkim/sgit/pkg/solar"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	commitMessage string
	skipLLM      bool
	interactive  bool
	skipEditor   bool
	useAI        bool
)

// commitCmd represents the commit command
var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Commit changes with optional AI-generated commit message",
	Long: `Commit changes to git. By default, uses AI to generate commit messages,
but supports all git commit options for full compatibility.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := runCommit(cmd, args); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
	DisableFlagParsing: false,
}

func init() {
	rootCmd.AddCommand(commitCmd)
	
	// AI-specific flags
	commitCmd.Flags().BoolVar(&skipLLM, "no-ai", false, "skip AI generation and use standard git commit")
	commitCmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "review and edit AI-generated message in terminal")
	commitCmd.Flags().BoolVar(&skipEditor, "skip-editor", false, "skip editor and use AI message directly")
	commitCmd.Flags().BoolVar(&useAI, "ai", false, "force AI generation even with other git flags")
	
	// Standard git commit flags - we'll pass these through to git
	commitCmd.Flags().StringVarP(&commitMessage, "message", "m", "", "commit message")
	commitCmd.Flags().BoolP("all", "a", false, "automatically stage modified and deleted files")
	commitCmd.Flags().Bool("amend", false, "amend the last commit")
	commitCmd.Flags().BoolP("verbose", "v", false, "show unified diff between HEAD and index")
	commitCmd.Flags().Bool("allow-empty", false, "allow empty commits")
	commitCmd.Flags().Bool("allow-empty-message", false, "allow empty commit message")
	commitCmd.Flags().String("author", "", "override author for commit")
	commitCmd.Flags().String("date", "", "override date for commit")
	commitCmd.Flags().Bool("signoff", false, "add Signed-off-by line")
	commitCmd.Flags().BoolP("patch", "p", false, "use interactive patch selection")
	commitCmd.Flags().String("fixup", "", "create a fixup commit")
	commitCmd.Flags().String("squash", "", "create a squash commit")
	commitCmd.Flags().Bool("reset-author", false, "reset author information")
	commitCmd.Flags().String("file", "", "read commit message from file")
	commitCmd.Flags().String("template", "", "use specified template file")
	commitCmd.Flags().Bool("edit", false, "force edit of commit message")
	commitCmd.Flags().Bool("no-edit", false, "don't edit commit message")
}

func runCommit(cmd *cobra.Command, args []string) error {
	// Check if we're in a git repository
	if !isGitRepository() {
		return fmt.Errorf("not a git repository")
	}

	// Handle -a flag: stage all modified and deleted files first
	if cmd.Flags().Changed("all") {
		allFlag, _ := cmd.Flags().GetBool("all")
		if allFlag {
			fmt.Println("Staging all modified and deleted files...")
			stageCmd := exec.Command("git", "add", "-u")
			if err := stageCmd.Run(); err != nil {
				return fmt.Errorf("error staging files with -a: %v", err)
			}
		}
	}

	// Only bypass AI in these specific cases:
	// 1. User provided explicit message with -m
	// 2. User explicitly disabled AI with --no-ai
	if commitMessage != "" || skipLLM {
		return executeGitCommitPassthrough(cmd, args)
	}

	// AI-enhanced commit logic for ALL other cases
	// Even with flags like --amend, --verbose, --signoff, etc.
	
	// Check for staged changes (required for AI generation)
	hasChanges, err := hasUncommittedChanges()
	if err != nil {
		return fmt.Errorf("error checking for changes: %v", err)
	}
	if !hasChanges {
		fmt.Println("No changes to commit")
		return nil
	}

	// Check configuration and setup if needed
	if err := ensureConfiguration(); err != nil {
		return err
	}

	// Get git diff
	diff, err := getGitDiff()
	if err != nil {
		return fmt.Errorf("error getting git diff: %v", err)
	}

	if strings.TrimSpace(diff) == "" {
		return fmt.Errorf("no diff found - make sure to add files with 'git add' first")
	}

	// Generate commit message using Solar LLM
	apiKey := viper.GetString("upstage_api_key")
	modelName := viper.GetString("upstage_model_name")
	
	client := solar.NewClient(apiKey, modelName, getEffectiveLanguage())
	
	fmt.Println("Generating comprehensive commit message with Solar LLM...")
	
	// Gather additional context for comprehensive commit message
	branch, _ := getCurrentBranch()
	recentCommits, _ := getRecentCommits(5)
	fileList, _ := getEnhancedFileList() // Use enhanced file list with content previews
	
	// Use comprehensive commit message generation with streaming
	generatedMessage, err := client.GenerateComprehensiveCommitMessageStream(diff, branch, recentCommits, fileList)
	
	if err != nil {
		return fmt.Errorf("error generating commit message: %v", err)
	}

	fmt.Println("\n✓ Commit message generated!")

	var finalMessage string

	// Handle different interaction modes
	if interactive {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Edit message (press Enter to use as-is): ")
		userInput, _ := reader.ReadString('\n')
		userInput = strings.TrimSpace(userInput)
		if userInput != "" {
			finalMessage = userInput
		} else {
			finalMessage = generatedMessage
		}
	} else if skipEditor {
		// Ask for confirmation before using AI message directly
		fmt.Print("Use this commit message? (y/n): ")
		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(strings.ToLower(response))
		if response != "y" && response != "yes" {
			fmt.Println("Commit cancelled")
			return nil
		}
		finalMessage = generatedMessage
	} else {
		// Default behavior: open editor with AI-generated message
		editedMessage, editorErr := openEditorWithMessage(generatedMessage)
		if editorErr != nil {
			return fmt.Errorf("error opening editor: %v", editorErr)
		}
		
		if strings.TrimSpace(editedMessage) == "" {
			fmt.Println("Empty commit message, aborting commit")
			return nil
		}
		
		finalMessage = editedMessage
	}

	// Execute git commit with AI message AND any additional flags
	return executeGitCommitWithFlags(finalMessage, cmd)
}

func executeGitCommitPassthrough(cobraCmd *cobra.Command, args []string) error {
	// Build git command with all flags and arguments
	gitArgs := []string{"commit"}
	
	// Add all the flags that were set
	cobraCmd.Flags().Visit(func(flag *pflag.Flag) {
		if flag.Name == "no-ai" || flag.Name == "interactive" || flag.Name == "skip-editor" || flag.Name == "ai" {
			return // Skip our custom flags
		}
		
		value := flag.Value.String()
		if flag.Value.Type() == "bool" && value == "true" {
			gitArgs = append(gitArgs, "--"+flag.Name)
		} else if flag.Value.Type() != "bool" && value != "" {
			gitArgs = append(gitArgs, "--"+flag.Name+"="+value)
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

func getDefaultEditor() string {
	// Try to get editor in order of preference (same as git)
	if editor := os.Getenv("GIT_EDITOR"); editor != "" {
		return editor
	}
	
	// Check git config for core.editor
	cmd := exec.Command("git", "config", "--get", "core.editor")
	if output, err := cmd.Output(); err == nil {
		if editor := strings.TrimSpace(string(output)); editor != "" {
			return editor
		}
	}
	
	if editor := os.Getenv("VISUAL"); editor != "" {
		return editor
	}
	
	if editor := os.Getenv("EDITOR"); editor != "" {
		return editor
	}
	
	// Default editors by platform
	if _, err := exec.LookPath("nano"); err == nil {
		return "nano"
	}
	
	if _, err := exec.LookPath("vim"); err == nil {
		return "vim"
	}
	
	if _, err := exec.LookPath("vi"); err == nil {
		return "vi"
	}
	
	return "nano" // fallback
}

func openEditorWithMessage(message string) (string, error) {
	// Create temporary file
	tmpDir := os.TempDir()
	tmpFile, err := ioutil.TempFile(tmpDir, "sgit-commit-*.txt")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write AI-generated message to temp file with some helpful comments
	content := fmt.Sprintf(`%s

# Please edit the commit message above.
# Lines starting with '#' will be ignored.
# An empty message aborts the commit.
#
# AI-generated message based on your changes.
# You can edit, replace, or completely rewrite it.
`, message)

	if _, err := tmpFile.WriteString(content); err != nil {
		tmpFile.Close()
		return "", fmt.Errorf("failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	// Get the editor
	editor := getDefaultEditor()
	
	// Split editor command (handle cases like "code --wait")
	editorParts := strings.Fields(editor)
	if len(editorParts) == 0 {
		return "", fmt.Errorf("no editor found")
	}

	// Run editor
	cmd := exec.Command(editorParts[0], append(editorParts[1:], tmpFile.Name())...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("editor exited with error: %v", err)
	}

	// Read the edited content
	editedBytes, err := ioutil.ReadFile(tmpFile.Name())
	if err != nil {
		return "", fmt.Errorf("failed to read edited file: %v", err)
	}

	// Process the content (remove comment lines and trim)
	lines := strings.Split(string(editedBytes), "\n")
	var resultLines []string
	
	for _, line := range lines {
		// Skip comment lines and empty lines at the end
		if !strings.HasPrefix(strings.TrimSpace(line), "#") {
			resultLines = append(resultLines, line)
		}
	}

	// Join and trim the result
	result := strings.TrimSpace(strings.Join(resultLines, "\n"))
	return result, nil
}

func isGitRepository() bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	err := cmd.Run()
	return err == nil
}

func hasUncommittedChanges() (bool, error) {
	cmd := exec.Command("git", "diff", "--cached", "--quiet")
	err := cmd.Run()
	if err != nil {
		// If git diff --cached --quiet fails, there are staged changes
		if exitError, ok := err.(*exec.ExitError); ok {
			return exitError.ExitCode() != 0, nil
		}
		return false, err
	}
	return false, nil
}

func getGitDiff() (string, error) {
	cmd := exec.Command("git", "diff", "--cached")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func executeGitCommit(message string) error {
	cmd := exec.Command("git", "commit", "-m", message)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// executeGitCommitWithFlags commits with AI message and preserves user's git flags
func executeGitCommitWithFlags(message string, cobraCmd *cobra.Command) error {
	// Build git command starting with commit and the AI message
	gitArgs := []string{"commit", "-m", message}
	
	// Add all the git flags that were set (excluding our custom AI flags)
	cobraCmd.Flags().Visit(func(flag *pflag.Flag) {
		// Skip our custom sgit flags
		if flag.Name == "no-ai" || flag.Name == "interactive" || flag.Name == "skip-editor" || flag.Name == "ai" {
			return
		}
		
		// Skip message flag since we're using the AI-generated message
		if flag.Name == "message" {
			return
		}
		
		// Add the flag to git command
		value := flag.Value.String()
		if flag.Value.Type() == "bool" && value == "true" {
			gitArgs = append(gitArgs, "--"+flag.Name)
		} else if flag.Value.Type() != "bool" && value != "" {
			gitArgs = append(gitArgs, "--"+flag.Name+"="+value)
		}
	})
	
	// Execute git command with AI message and all user flags
	gitCmd := exec.Command("git", gitArgs...)
	gitCmd.Stdin = os.Stdin
	gitCmd.Stdout = os.Stdout
	gitCmd.Stderr = os.Stderr
	return gitCmd.Run()
}

func executeInteractiveGitCommit() error {
	cmd := exec.Command("git", "commit")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func getCurrentBranch() (string, error) {
	cmd := exec.Command("git", "branch", "--show-current")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func getRecentCommits(count int) (string, error) {
	cmd := exec.Command("git", "log", fmt.Sprintf("-%d", count), "--oneline", "--no-merges")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func getEnhancedFileList() (string, error) {
	// Get list of staged files
	stagedCmd := exec.Command("git", "diff", "--cached", "--name-status")
	stagedOutput, err := stagedCmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get staged files: %w", err)
	}

	if len(stagedOutput) == 0 {
		return "No files staged for commit", nil
	}

	var fileInfo []string
	lines := strings.Split(strings.TrimSpace(string(stagedOutput)), "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}
		
		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}
		
		status := parts[0]
		filePath := parts[1]
		
		// Get file size
		fileSize := "unknown"
		if stat, err := os.Stat(filePath); err == nil {
			fileSize = fmt.Sprintf("%d bytes", stat.Size())
		}
		
		fileDesc := fmt.Sprintf("- %s %s (%s)", status, filePath, fileSize)
		
		// For new files (A = Added), include content preview
		if status == "A" && !isBinaryFile(filePath) {
			if stat, err := os.Stat(filePath); err == nil && stat.Size() <= 50*1024 { // Only for files <= 50KB
				contentPreview := getFileContentPreview(filePath, 20) // First 20 lines
				fileDesc += fmt.Sprintf("\n  Content preview:\n%s", 
					strings.ReplaceAll(contentPreview, "\n", "\n  "))
			}
		}
		
		fileInfo = append(fileInfo, fileDesc)
	}

	return strings.Join(fileInfo, "\n"), nil
}

// Helper function to get file content preview for new files
func getFileContentPreview(filePath string, maxLines int) string {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Sprintf("Error reading file: %v", err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	lineCount := 0
	
	for scanner.Scan() && lineCount < maxLines {
		lines = append(lines, scanner.Text())
		lineCount++
	}
	
	if err := scanner.Err(); err != nil {
		return fmt.Sprintf("Error scanning file: %v", err)
	}
	
	content := strings.Join(lines, "\n")
	if lineCount == maxLines {
		content += "\n... (file continues)"
	}
	
	return content
}

