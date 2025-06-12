package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var langFlag string
var version = "dev" // Will be set during build with -ldflags

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "sgit",
	Short: "Solar LLM-powered git wrapper",
	Long: `sgit is a git wrapper that uses Solar LLM to automatically generate
commit messages based on your code changes.`,
	SilenceUsage:  true,
	SilenceErrors: true,
	Version:       version, // Will be set during build
}

// executeGitPassthrough passes commands directly to git
func executeGitPassthrough(args []string) error {
	gitArgs := append([]string{}, args...)

	gitCmd := exec.Command("git", gitArgs...)
	gitCmd.Stdin = os.Stdin
	gitCmd.Stdout = os.Stdout
	gitCmd.Stderr = os.Stderr

	if err := gitCmd.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			os.Exit(exitError.ExitCode())
		}
		os.Exit(1)
	}

	return nil
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	err := rootCmd.Execute()

	// If it's an unknown command error, try to pass it through to git
	if err != nil && strings.Contains(err.Error(), "unknown command") {
		// Get the original args
		args := os.Args[1:] // Skip the program name
		if len(args) > 0 {
			// Execute git command and exit with its status
			gitCmd := exec.Command("git", args...)
			gitCmd.Stdin = os.Stdin
			gitCmd.Stdout = os.Stdout
			gitCmd.Stderr = os.Stderr

			if gitErr := gitCmd.Run(); gitErr != nil {
				if exitError, ok := gitErr.(*exec.ExitError); ok {
					os.Exit(exitError.ExitCode())
				}
				os.Exit(1)
			}
			return // Success
		}
	}

	// Handle other errors
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// getEffectiveLanguage returns the language to use, considering both config and flag
func getEffectiveLanguage() string {
	// Command-line flag takes precedence
	if langFlag != "" {
		lang := strings.ToLower(strings.TrimSpace(langFlag))
		if isValidLanguageCode(lang) {
			return lang
		}
		fmt.Fprintf(os.Stderr, "Warning: Invalid language code '%s'. Using default 'en'.\n", langFlag)
		return "en"
	}

	// Fall back to config file setting
	configLang := viper.GetString("language")
	if configLang != "" {
		lang := strings.ToLower(strings.TrimSpace(configLang))
		if isValidLanguageCode(lang) {
			return lang
		}
		return "en"
	}

	// Default to English
	return "en"
}

// isValidLanguageCode checks if the provided language code is supported
func isValidLanguageCode(code string) bool {
	validCodes := map[string]bool{
		"en": true,
		"ko": true,
		"ja": true,
		"zh": true,
		"es": true,
		"fr": true,
		"de": true,
	}
	return validCodes[code]
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/sgit/config.yaml)")
	rootCmd.PersistentFlags().StringVar(&langFlag, "lang", "", "language for AI responses (en|ko|ja|zh|es|fr|de, overrides config setting)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error finding home directory: %v\n", err)
			os.Exit(1)
		}

		// Search config in home directory with name ".sgitrc" (without extension).
		configDir := filepath.Join(home, ".config", "sgit")
		viper.AddConfigPath(configDir)
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")

		// Create config directory if it doesn't exist
		if err := os.MkdirAll(configDir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating config directory: %v\n", err)
		}
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		// Config file loaded successfully
	}
}
