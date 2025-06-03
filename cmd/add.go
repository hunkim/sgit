package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/hunkim/sgit/pkg/solar"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	addAll    bool
	addForce  bool
	addDryRun bool
	addAI     bool
)

// addCmd represents the smart add command
var addCmd = &cobra.Command{
	Use:   "add [files...]",
	Short: "Add files to git with optional AI analysis",
	Long: `Add files to git staging area. By default, analyzes untracked files with AI
when no specific files are given. Supports all git add options for full compatibility.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := runSmartAdd(cmd, args); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	
	// AI-specific flags (custom to sgit)
	addCmd.Flags().BoolVar(&addAll, "all-ai", false, "analyze all untracked files with AI")
	addCmd.Flags().BoolVar(&addForce, "force-ai", false, "add files without AI confirmation (smart filtering only)")
	addCmd.Flags().BoolVar(&addDryRun, "dry-run-ai", false, "show what would be added without actually adding")
	addCmd.Flags().BoolVar(&addAI, "ai", false, "force AI analysis even with specific files")
	
	// Standard git add flags - we'll pass these through to git
	addCmd.Flags().BoolP("all", "A", false, "add all changes (git standard)")
	addCmd.Flags().BoolP("update", "u", false, "update tracked files")
	addCmd.Flags().BoolP("patch", "p", false, "interactively choose hunks")
	addCmd.Flags().BoolP("verbose", "v", false, "be verbose")
	addCmd.Flags().BoolP("dry-run", "n", false, "don't actually add files")
	addCmd.Flags().Bool("ignore-errors", false, "continue on add errors")
	addCmd.Flags().Bool("ignore-missing", false, "ignore missing files")
	addCmd.Flags().Bool("no-warn-embedded-repo", false, "don't warn about embedded repositories")
	addCmd.Flags().Bool("renormalize", false, "renormalize EOL of tracked files")
	addCmd.Flags().Bool("chmod", false, "override executable bit")
	addCmd.Flags().BoolP("intent-to-add", "N", false, "record only intent to add")
	addCmd.Flags().Bool("refresh", false, "refresh stat information")
	addCmd.Flags().Bool("ignore-removal", false, "ignore removal of files")
	addCmd.Flags().String("pathspec-from-file", "", "read pathspec from file")
	addCmd.Flags().Bool("pathspec-file-nul", false, "pathspec file is NUL separated")
}

func runSmartAdd(cmd *cobra.Command, args []string) error {
	// Check if we're in a git repository
	if !isGitRepository() {
		return fmt.Errorf("not a git repository")
	}

	// Check if any git-specific flags are set that should bypass AI
	shouldUseGitDirectly := shouldBypassAIForAdd(cmd)
	
	// If specific files are provided or git flags are used, use git behavior
	if (len(args) > 0 && !addAI) || (shouldUseGitDirectly && !addAI) {
		return executeGitAddPassthrough(cmd, args)
	}

	// Only use AI analysis when --all-ai flag is used or no args and no git flags
	if !addAll && len(args) == 0 {
		fmt.Println("Use 'sgit add --all-ai' for AI analysis of untracked files")
		fmt.Println("Use 'sgit add <files>' for standard git add behavior")
		return nil
	}

	// AI-enhanced add logic (only when explicitly requested)
	
	// Check configuration and setup if needed (unless in force mode)
	if !addForce {
		if err := ensureConfiguration(); err != nil {
			return err
		}
	}

	// Get untracked files
	untrackedFiles, err := getUntrackedFiles()
	if err != nil {
		return fmt.Errorf("error getting untracked files: %v", err)
	}

	if len(untrackedFiles) == 0 {
		fmt.Println("No untracked files found")
		return nil
	}

	fmt.Printf("Found %d untracked files. Analyzing with Solar LLM...\n", len(untrackedFiles))

	// Analyze each file
	filesToAdd := []string{}
	for _, file := range untrackedFiles {
		// Skip binary files
		if isBinaryFile(file) {
			fmt.Printf("⏭️  Skipping binary file: %s\n", file)
			continue
		}

		// Skip if file is too large (> 1MB)
		if isLargeFile(file) {
			fmt.Printf("⏭️  Skipping large file: %s\n", file)
			continue
		}

		if addForce {
			filesToAdd = append(filesToAdd, file)
			fmt.Printf("✅ Will add: %s (force mode)\n", file)
			continue
		}

		// Use AI to analyze the file
		shouldAdd, reason, err := analyzeFileWithAI(file)
		if err != nil {
			fmt.Printf("❌ Error analyzing %s: %v\n", file, err)
			continue
		}

		if shouldAdd {
			fmt.Printf("✅ Recommended to add: %s\n   Reason: %s\n", file, reason)
			filesToAdd = append(filesToAdd, file)
		} else {
			fmt.Printf("❌ Recommended to skip: %s\n   Reason: %s\n", file, reason)
		}
	}

	if len(filesToAdd) == 0 {
		fmt.Println("No files recommended for adding")
		return nil
	}

	// Show summary and ask for confirmation
	fmt.Printf("\nFiles recommended for adding:\n")
	for _, file := range filesToAdd {
		fmt.Printf("  - %s\n", file)
	}

	if addDryRun {
		fmt.Println("\n[DRY RUN] No files were actually added")
		return nil
	}

	if !addForce {
		fmt.Print("\nAdd these files? (y/n): ")
		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(strings.ToLower(response))
		if response != "y" && response != "yes" {
			fmt.Println("Add cancelled")
			return nil
		}
	}

	// Add the files
	return executeGitAdd(filesToAdd)
}

func shouldBypassAIForAdd(cmd *cobra.Command) bool {
	// Check for flags that indicate user wants standard git behavior
	flags := []string{
		"all", "update", "patch", "verbose", "dry-run", "ignore-errors",
		"ignore-missing", "no-warn-embedded-repo", "renormalize", "chmod",
		"intent-to-add", "refresh", "ignore-removal", "pathspec-from-file",
		"pathspec-file-nul",
	}
	
	for _, flag := range flags {
		if cmd.Flags().Lookup(flag) != nil && cmd.Flags().Changed(flag) {
			return true
		}
	}
	
	return false
}

func executeGitAddPassthrough(cobraCmd *cobra.Command, args []string) error {
	// Build git command with all flags and arguments
	gitArgs := []string{"add"}
	
	// Add all the flags that were set (excluding our custom AI flags)
	cobraCmd.Flags().Visit(func(flag *pflag.Flag) {
		flagName := flag.Name
		if strings.HasSuffix(flagName, "-ai") || flagName == "ai" {
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
	
	// Add any remaining arguments (files)
	gitArgs = append(gitArgs, args...)
	
	// Execute git command
	gitCmd := exec.Command("git", gitArgs...)
	gitCmd.Stdin = os.Stdin
	gitCmd.Stdout = os.Stdout
	gitCmd.Stderr = os.Stderr
	return gitCmd.Run()
}

func getUntrackedFiles() ([]string, error) {
	cmd := exec.Command("git", "ls-files", "--others", "--exclude-standard")
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

func isBinaryFile(filename string) bool {
	// Check file extension for common binary types
	ext := strings.ToLower(filepath.Ext(filename))
	binaryExts := []string{
		".exe", ".dll", ".so", ".dylib", ".a", ".o", ".obj",
		".jpg", ".jpeg", ".png", ".gif", ".bmp", ".ico", ".tiff", ".svg",
		".mp3", ".mp4", ".avi", ".mov", ".mkv", ".flv", ".wav",
		".zip", ".tar", ".gz", ".bz2", ".xz", ".7z", ".rar",
		".pdf", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx",
		".bin", ".dat", ".db", ".sqlite", ".sqlite3",
	}

	for _, bext := range binaryExts {
		if ext == bext {
			return true
		}
	}

	// Check if file contains null bytes (binary indicator)
	file, err := os.Open(filename)
	if err != nil {
		return false
	}
	defer file.Close()

	// Read first 512 bytes to check for binary content
	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && n == 0 {
		return false
	}

	// Check for null bytes
	for i := 0; i < n; i++ {
		if buffer[i] == 0 {
			return true
		}
	}

	return false
}

func isLargeFile(filename string) bool {
	info, err := os.Stat(filename)
	if err != nil {
		return false
	}
	// Skip files larger than 1MB
	return info.Size() > 1024*1024
}

func analyzeFileWithAI(filename string) (bool, string, error) {
	// Read file content
	content, err := os.ReadFile(filename)
	if err != nil {
		return false, "", fmt.Errorf("error reading file: %v", err)
	}

	// Limit content size for API (max 4KB)
	contentStr := string(content)
	if len(contentStr) > 4096 {
		contentStr = contentStr[:4096] + "\n... [truncated]"
	}

	apiKey := viper.GetString("upstage_api_key")
	modelName := viper.GetString("upstage_model_name")
	
	client := solar.NewClient(apiKey, modelName, getEffectiveLanguage())
	
	prompt := fmt.Sprintf(`You are a helpful assistant that analyzes files in software projects to determine if they should be added to git version control.

Analyze the following file and determine if it should be added to git:

File: %s
Content:
%s

Consider these factors:
1. Is this a source code file, configuration, or documentation that belongs in version control?
2. Is this a temporary file, log file, or build artifact that should be ignored?
3. Does this file contain sensitive information (passwords, keys, tokens)?
4. Is this a generated file that can be recreated from source?

Respond with only:
- "YES: [brief reason]" if the file should be added
- "NO: [brief reason]" if the file should not be added

Keep the reason under 50 characters.`, filename, contentStr)

	response, err := client.GenerateResponse(prompt)
	if err != nil {
		return false, "", err
	}

	response = strings.TrimSpace(response)
	
	if strings.HasPrefix(strings.ToUpper(response), "YES:") {
		reason := strings.TrimSpace(strings.TrimPrefix(response, "YES:"))
		if strings.HasPrefix(strings.ToUpper(reason), "YES:") {
			reason = strings.TrimSpace(strings.TrimPrefix(reason, "YES:"))
		}
		return true, reason, nil
	} else if strings.HasPrefix(strings.ToUpper(response), "NO:") {
		reason := strings.TrimSpace(strings.TrimPrefix(response, "NO:"))
		if strings.HasPrefix(strings.ToUpper(reason), "NO:") {
			reason = strings.TrimSpace(strings.TrimPrefix(reason, "NO:"))
		}
		return false, reason, nil
	}

	// Fallback: if response doesn't match expected format, be conservative
	return false, "AI response unclear, skipping for safety", nil
}

func executeGitAdd(files []string) error {
	if len(files) == 0 {
		return nil
	}

	args := append([]string{"add"}, files...)
	cmd := exec.Command("git", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error adding files: %v", err)
	}

	fmt.Printf("Successfully added %d files\n", len(files))
	return nil
} 