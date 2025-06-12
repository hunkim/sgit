package solar

import (
	"strings"
)

const (
	// Maximum tokens allowed for input (40K as requested by user)
	MaxInputTokens = 40000
	// Maximum words to stay under 40K tokens (40K / 1.5 = ~27K words)
	MaxInputWords = 27000
	// Solar Pro model's actual context limit
	ModelContextLimit = 65536
)

// TokenCounter provides functionality to count tokens using Solar Pro tokenizer logic
type TokenCounter struct {
	vocabSize int
}

// NewTokenCounter creates a new token counter
func NewTokenCounter() *TokenCounter {
	return &TokenCounter{
		vocabSize: 32128, // Solar Pro tokenizer vocab size
	}
}

// EstimateTokens provides a simple word-based token estimation
func (tc *TokenCounter) EstimateTokens(text string) int {
	if text == "" {
		return 0
	}

	// Simple word count
	words := strings.Fields(text)
	wordCount := len(words)

	// For code/diff content, assume 1.5 tokens per word (conservative)
	// This accounts for special characters, variable names, operators, etc.
	estimatedTokens := int(float64(wordCount) * 1.5)

	return estimatedTokens
}

// TruncateToWordLimit truncates text to fit within the specified word limit
func (tc *TokenCounter) TruncateToWordLimit(text string, maxWords int) (string, int) {
	if text == "" {
		return "", 0
	}

	words := strings.Fields(text)
	if len(words) <= maxWords {
		return text, len(words)
	}

	// Take the first N words and add truncation notice
	truncatedWords := words[:maxWords]
	truncatedText := strings.Join(truncatedWords, " ")
	truncatedText += "\n\n[... truncated to stay within token limit ...]"

	return truncatedText, maxWords
}

// CountWords returns the number of words in the text
func (tc *TokenCounter) CountWords(text string) int {
	return len(strings.Fields(text))
}

// SplitContent intelligently splits content into sections for better truncation
func (tc *TokenCounter) SplitContent(diff, branch, recentCommits, fileList string) (string, string, string, string, int) {
	// Calculate words for each section
	diffWords := tc.CountWords(diff)
	branchWords := tc.CountWords(branch)
	recentCommitsWords := tc.CountWords(recentCommits)
	fileListWords := tc.CountWords(fileList)

	totalWords := diffWords + branchWords + recentCommitsWords + fileListWords

	// If total is within limit, return as-is
	if totalWords <= MaxInputWords {
		return diff, branch, recentCommits, fileList, totalWords
	}

	// Priority order: diff (most important), fileList, recentCommits, branch
	remainingWords := MaxInputWords

	// Always preserve branch info (small)
	if branchWords < remainingWords {
		remainingWords -= branchWords
	} else {
		branch, _ = tc.TruncateToWordLimit(branch, remainingWords/4)
		remainingWords -= remainingWords / 4
	}

	// Allocate words proportionally, with diff getting priority
	diffAllocation := int(float64(remainingWords) * 0.6)                            // 60% for diff
	fileListAllocation := int(float64(remainingWords) * 0.25)                       // 25% for file list
	recentCommitsAllocation := remainingWords - diffAllocation - fileListAllocation // remainder for recent commits

	// Truncate each section
	truncatedDiff, actualDiffWords := tc.TruncateToWordLimit(diff, diffAllocation)
	truncatedFileList, actualFileListWords := tc.TruncateToWordLimit(fileList, fileListAllocation)
	truncatedRecentCommits, actualRecentCommitsWords := tc.TruncateToWordLimit(recentCommits, recentCommitsAllocation)

	actualTotal := actualDiffWords + branchWords + actualRecentCommitsWords + actualFileListWords

	return truncatedDiff, branch, truncatedRecentCommits, truncatedFileList, actualTotal
}

// TruncateContent truncates a single content input to fit within word limits
func (tc *TokenCounter) TruncateContent(content string) (string, int, bool) {
	words := tc.CountWords(content)
	if words <= MaxInputWords {
		return content, words, false
	}

	truncated, actualWords := tc.TruncateToWordLimit(content, MaxInputWords)
	return truncated, actualWords, true
}
