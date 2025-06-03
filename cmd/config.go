package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/term"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure sgit settings",
	Long:  `Configure API key and other settings for sgit.`,
	Run: func(cmd *cobra.Command, args []string) {
		setupConfig()
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}

func readAPIKeyWithVisualFeedback() (string, error) {
	var apiKey []byte
	var char byte
	var err error

	for {
		// Read one character at a time
		if char, err = readSingleChar(); err != nil {
			return "", err
		}

		// Handle Ctrl-C (ASCII 3)
		if char == 3 {
			fmt.Println("\n\n⚠️  Configuration cancelled by user")
			fmt.Println("💡 Run 'sgit config' again anytime to set up your configuration")
			os.Exit(0)
		}

		// Handle Enter (newline)
		if char == '\n' || char == '\r' {
			fmt.Println()
			break
		}

		// Handle Backspace
		if char == 127 || char == 8 {
			if len(apiKey) > 0 {
				apiKey = apiKey[:len(apiKey)-1]
				fmt.Print("\b \b") // Move back, print space, move back again
			}
			continue
		}

		// Add character to password
		apiKey = append(apiKey, char)

		// Display feedback
		if len(apiKey) <= 3 {
			// Show first 3 characters (should be "up_")
			fmt.Print(string(char))
		} else {
			// Show asterisk for characters after the first 3
			fmt.Print("*")
		}
	}

	return string(apiKey), nil
}

func readSingleChar() (byte, error) {
	// Set terminal to raw mode
	oldState, err := term.MakeRaw(int(syscall.Stdin))
	if err != nil {
		return 0, err
	}
	defer term.Restore(int(syscall.Stdin), oldState)

	// Read one byte
	var b [1]byte
	_, err = os.Stdin.Read(b[:])
	return b[0], err
}

func setupConfig() {
	// Set up signal handling for Ctrl-C
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	
	// Handle interrupt in a goroutine
	go func() {
		<-sigChan
		fmt.Println("\n\n⚠️  Configuration cancelled by user")
		fmt.Println("💡 Run 'sgit config' again anytime to set up your configuration")
		os.Exit(0)
	}()
	
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("🔧 sgit Configuration Setup")
	fmt.Println("Your API key will be stored locally and securely in ~/.config/sgit/config.yaml")
	fmt.Println("💡 Press Ctrl-C anytime to cancel")
	fmt.Println()

	// Check existing configuration
	existingAPIKey := viper.GetString("upstage_api_key")
	existingModelName := viper.GetString("upstage_model_name")
	existingLanguage := viper.GetString("language")

	var apiKeyStr string
	var err error

	// Get API key
	fmt.Println("(get one at https://console.upstage.ai/)")
	if existingAPIKey != "" {
		// Show masked existing API key
		maskedKey := ""
		if len(existingAPIKey) >= 3 {
			maskedKey = existingAPIKey[:3] + strings.Repeat("*", len(existingAPIKey)-3)
		} else {
			maskedKey = strings.Repeat("*", len(existingAPIKey))
		}
		fmt.Printf("Enter your Upstage API key (current: %s, press Enter to keep): ", maskedKey)
		
		// For existing keys, use simple input to allow easy Enter-to-keep
		input, err := reader.ReadString('\n')
		if err != nil {
			// Check if it's an interrupt
			if strings.Contains(err.Error(), "interrupt") {
				fmt.Println("\n\n⚠️  Configuration cancelled by user")
				fmt.Println("💡 Run 'sgit config' again anytime to set up your configuration")
				os.Exit(0)
			}
			fmt.Printf("Error reading API key: %v\n", err)
			return
		}
		input = strings.TrimSpace(input)
		
		if input == "" {
			apiKeyStr = existingAPIKey
			fmt.Println("✓ Keeping existing API key")
		} else {
			apiKeyStr = input
		}
	} else {
		fmt.Print("Enter your Upstage API key: ")
		apiKeyStr, err = readAPIKeyWithVisualFeedback()
		if err != nil {
			fmt.Printf("\nError reading API key: %v\n", err)
			return
		}
		
		if apiKeyStr == "" {
			fmt.Println("API key cannot be empty")
			return
		}
	}

	// Get model name with existing value
	defaultModel := "solar-pro2-preview"
	if existingModelName != "" {
		fmt.Printf("Enter model name (current: %s, press Enter to keep): ", existingModelName)
	} else {
		fmt.Printf("Enter model name (default: %s): ", defaultModel)
	}
	
	modelName, err := reader.ReadString('\n')
	if err != nil {
		// Check if it's an interrupt
		if strings.Contains(err.Error(), "interrupt") {
			fmt.Println("\n\n⚠️  Configuration cancelled by user")
			fmt.Println("💡 Run 'sgit config' again anytime to set up your configuration")
			os.Exit(0)
		}
		fmt.Printf("Error reading model name: %v\n", err)
		return
	}
	modelName = strings.TrimSpace(modelName)
	
	// Use existing value if empty, otherwise use default
	if modelName == "" {
		if existingModelName != "" {
			modelName = existingModelName
			fmt.Printf("✓ Keeping existing model: %s\n", modelName)
		} else {
			modelName = defaultModel
		}
	}

	// Get language preference with existing value
	fmt.Println("\nAvailable languages:")
	fmt.Println("  en - English")
	fmt.Println("  ko - Korean (한국어)")
	fmt.Println("  ja - Japanese (日本語)")
	fmt.Println("  zh - Chinese (中文)")
	fmt.Println("  es - Spanish (Español)")
	fmt.Println("  fr - French (Français)")
	fmt.Println("  de - German (Deutsch)")
	
	if existingLanguage != "" {
		validLanguages := map[string]string{
			"en": "English",
			"ko": "Korean",
			"ja": "Japanese",
			"zh": "Chinese",
			"es": "Spanish",
			"fr": "French",
			"de": "German",
		}
		currentLangName := validLanguages[existingLanguage]
		if currentLangName == "" {
			currentLangName = existingLanguage
		}
		fmt.Printf("Enter language code (current: %s - %s, press Enter to keep): ", existingLanguage, currentLangName)
	} else {
		fmt.Print("Enter language code (default: en): ")
	}
	
	language, err := reader.ReadString('\n')
	if err != nil {
		// Check if it's an interrupt
		if strings.Contains(err.Error(), "interrupt") {
			fmt.Println("\n\n⚠️  Configuration cancelled by user")
			fmt.Println("💡 Run 'sgit config' again anytime to set up your configuration")
			os.Exit(0)
		}
		fmt.Printf("Error reading language: %v\n", err)
		return
	}
	language = strings.TrimSpace(strings.ToLower(language))
	
	// Use existing value if empty, otherwise use default
	if language == "" {
		if existingLanguage != "" {
			language = existingLanguage
			fmt.Printf("✓ Keeping existing language: %s\n", language)
		} else {
			language = "en"
		}
	}
	
	// Validate language code
	validLanguages := map[string]string{
		"en": "English",
		"ko": "Korean",
		"ja": "Japanese",
		"zh": "Chinese",
		"es": "Spanish",
		"fr": "French",
		"de": "German",
	}
	
	if _, valid := validLanguages[language]; !valid {
		fmt.Printf("Invalid language code '%s'. Defaulting to 'en' (English)\n", language)
		language = "en"
	} else {
		fmt.Printf("Selected language: %s (%s)\n", language, validLanguages[language])
	}

	// Save configuration
	viper.Set("upstage_api_key", apiKeyStr)
	viper.Set("upstage_model_name", modelName)
	viper.Set("language", language)

	// Get config file path
	configDir := filepath.Join(os.Getenv("HOME"), ".config", "sgit")
	configFile := filepath.Join(configDir, "config.yaml")

	if err := viper.WriteConfigAs(configFile); err != nil {
		fmt.Printf("Error saving configuration: %v\n", err)
		return
	}

	fmt.Printf("\n✅ Configuration saved securely to %s\n", configFile)
	
	// Stop listening for signals since we're done
	signal.Stop(sigChan)
}

// ensureConfiguration checks if configuration exists and runs setup if needed
func ensureConfiguration() error {
	apiKey := viper.GetString("upstage_api_key")
	if apiKey == "" {
		fmt.Println("No API key configured. Running setup...")
		fmt.Println()
		setupConfig()
		
		// Re-read configuration after setup
		apiKey = viper.GetString("upstage_api_key")
		if apiKey == "" {
			return fmt.Errorf("configuration setup failed or was cancelled")
		}
		
		fmt.Println()
		fmt.Println("Configuration complete! Continuing...")
	}
	return nil
} 