package solar

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"bufio"
	"time"
	"sync"
	"os"
)

// Client represents the Solar LLM API client
type Client struct {
	apiKey    string
	modelName string
	baseURL   string
}

// Message represents a chat message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest represents the request structure for Solar LLM API
type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

// ChatResponse represents the response structure from Solar LLM API
type ChatResponse struct {
	Choices []Choice `json:"choices"`
}

// StreamResponse represents a streaming response chunk
type StreamResponse struct {
	Choices []StreamChoice `json:"choices"`
}

// StreamChoice represents a choice in streaming response
type StreamChoice struct {
	Delta StreamDelta `json:"delta"`
}

// StreamDelta represents the delta content in streaming
type StreamDelta struct {
	Content string `json:"content"`
}

// Choice represents a choice in the response
type Choice struct {
	Message Message `json:"message"`
}

// Options for the request (commented out as it might not be supported)
// type Options struct {
// 	ReasoningEffort string `json:"reasoning_effort"`
// }

// Spinner represents a loading spinner
type Spinner struct {
	chars    []string
	delay    time.Duration
	active   bool
	mu       sync.RWMutex
	stopChan chan bool
}

// NewSpinner creates a new spinner with default settings
func NewSpinner() *Spinner {
	// Try to use Unicode spinner if available, fall back to ASCII
	unicodeSpinner := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	
	// ASCII fallback spinner
	asciiSpinner := []string{"|", "/", "-", "\\"}
	
	// Check if we're likely in a terminal that supports Unicode
	termType := os.Getenv("TERM")
	if strings.Contains(termType, "xterm") || strings.Contains(termType, "screen") || 
	   termType == "" || strings.Contains(termType, "color") {
		return &Spinner{
			chars:    unicodeSpinner,
			delay:    100 * time.Millisecond,
			stopChan: make(chan bool, 1),
		}
	}
	
	// Fall back to ASCII for older/simpler terminals
	return &Spinner{
		chars:    asciiSpinner,
		delay:    100 * time.Millisecond,
		stopChan: make(chan bool, 1),
	}
}

// Start begins the spinner animation
func (s *Spinner) Start(message string) {
	s.mu.Lock()
	s.active = true
	s.mu.Unlock()

	go func() {
		i := 0
		for {
			select {
			case <-s.stopChan:
				return
			default:
				s.mu.RLock()
				if !s.active {
					s.mu.RUnlock()
					return
				}
				s.mu.RUnlock()

				fmt.Printf("\r%s %s", s.chars[i%len(s.chars)], message)
				i++
				time.Sleep(s.delay)
			}
		}
	}()
}

// Stop ends the spinner animation and clears the line
func (s *Spinner) Stop() {
	s.mu.Lock()
	s.active = false
	s.mu.Unlock()

	select {
	case s.stopChan <- true:
	default:
	}

	// Clear the spinner line
	fmt.Print("\r" + strings.Repeat(" ", 60) + "\r")
}

// NewClient creates a new Solar LLM client
func NewClient(apiKey, modelName string) *Client {
	if modelName == "" {
		modelName = "solar-pro2-preview"
	}
	return &Client{
		apiKey:    apiKey,
		modelName: modelName,
		baseURL:   "https://api.upstage.ai/v1/chat/completions",
	}
}

// GenerateCommitMessage generates a commit message based on the git diff
func (c *Client) GenerateCommitMessage(diff string) (string, error) {
	prompt := fmt.Sprintf(`You are an expert software developer who writes excellent commit messages following the Conventional Commits specification.

Analyze the following git diff and generate a concise, descriptive commit message:

%s

Guidelines:
1. Use conventional commit format: type(scope): description
2. Types: feat, fix, docs, style, refactor, test, chore, perf, ci, build
3. Include scope if relevant (e.g., auth, api, ui, db)
4. Description should be imperative mood ("add" not "added")
5. Keep first line under 50 characters if possible
6. If changes are complex, add a brief body explaining the what and why

Examples:
- feat(auth): add OAuth2 integration
- fix(api): handle null pointer in user service
- docs: update installation instructions
- refactor(db): optimize query performance

Respond with only the commit message, no explanations.`, diff)

	return c.GenerateResponse(prompt)
}

// GenerateComprehensiveCommitMessage generates a comprehensive commit message based on the git diff, branch, recent commits, and file list
func (c *Client) GenerateComprehensiveCommitMessage(diff, branch, recentCommits, fileList string) (string, error) {
	prompt := fmt.Sprintf(`You are an expert software developer who writes excellent commit messages following the Conventional Commits specification.

Analyze the following information to generate a detailed but concise commit message (200-400 characters):

=== GIT DIFF ===
%s

=== CURRENT BRANCH ===
%s

=== RECENT COMMITS (last 5) ===
%s

=== FILES CHANGED ===
%s

Generate a commit message that:
1. Follows conventional commit format: type(scope): description
2. Types: feat, fix, docs, style, refactor, test, chore, perf, ci, build
3. Considers the context of recent commits and branch name
4. Reflects the scope and impact of changes
5. Uses imperative mood ("add" not "added")
6. Includes a brief body (2-3 lines) explaining:
   - What specifically was changed and why
   - Key files/components affected
   - Any important technical details or impacts
7. Mentions breaking changes if applicable
8. Keep total length between 200-400 characters

Format:
type(scope): concise but descriptive summary

Brief explanation of what was implemented, modified, or fixed.
Mention key files affected and any important technical considerations.

BREAKING CHANGE: description if applicable (only if truly breaking)

Respond with only the commit message, no explanations.`, diff, branch, recentCommits, fileList)

	return c.GenerateResponse(prompt)
}

// GenerateComprehensiveCommitMessageStream generates a commit message with streaming
func (c *Client) GenerateComprehensiveCommitMessageStream(diff, branch, recentCommits, fileList string) (string, error) {
	prompt := fmt.Sprintf(`You are an expert software developer who writes excellent commit messages following the Conventional Commits specification.

Analyze the following information to generate a detailed but concise commit message (200-400 characters):

=== GIT DIFF ===
%s

=== CURRENT BRANCH ===
%s

=== RECENT COMMITS (last 5) ===
%s

=== FILES CHANGED ===
%s

Generate a commit message that:
1. Follows conventional commit format: type(scope): description
2. Types: feat, fix, docs, style, refactor, test, chore, perf, ci, build
3. Considers the context of recent commits and branch name
4. Reflects the scope and impact of changes
5. Uses imperative mood ("add" not "added")
6. Includes a brief body (2-3 lines) explaining:
   - What specifically was changed and why
   - Key files/components affected
   - Any important technical details or impacts
7. Mentions breaking changes if applicable
8. Keep total length between 200-400 characters

Format:
type(scope): concise but descriptive summary

Brief explanation of what was implemented, modified, or fixed.
Mention key files affected and any important technical considerations.

BREAKING CHANGE: description if applicable (only if truly breaking)

Respond with only the commit message, no explanations.`, diff, branch, recentCommits, fileList)

	return c.GenerateResponseStream(prompt)
}

// SummarizeDiff generates a summary of the git diff
func (c *Client) SummarizeDiff(diff string) (string, error) {
	prompt := fmt.Sprintf(`Analyze the following git diff and provide a clear, concise summary of the changes:

%s

Provide:
1. **Summary**: One-line overview of what changed
2. **Files Modified**: List of main files/components affected
3. **Type of Changes**: New features, bug fixes, refactoring, etc.
4. **Impact**: Potential effects on functionality
5. **Notable**: Any important details (breaking changes, performance impacts, etc.)

Keep it concise but informative.`, diff)

	return c.GenerateResponse(prompt)
}

// AnalyzeLog generates insights from the git log
func (c *Client) AnalyzeLog(logOutput, timeframe string) (string, error) {
	prompt := fmt.Sprintf(`Analyze the following git log (%s) and provide insights:

%s

Provide:
1. **Activity Summary**: Overall development activity
2. **Key Features**: Major features or changes
3. **Bug Fixes**: Important fixes
4. **Contributors**: Active contributors and their focus areas
5. **Patterns**: Development patterns, frequency, focus areas
6. **Recommendations**: Suggestions for the project

Be concise but insightful.`, timeframe, logOutput)

	return c.GenerateResponse(prompt)
}

// AnalyzeMergeConflicts provides guidance for resolving merge conflicts
func (c *Client) AnalyzeMergeConflicts(conflictFiles string) (string, error) {
	prompt := fmt.Sprintf(`Analyze the following merge conflict information and provide resolution guidance:

%s

Provide:
1. **Conflict Summary**: What files have conflicts and why
2. **Resolution Strategy**: Recommended approach for resolving
3. **Risk Assessment**: Potential risks of different resolution approaches
4. **Testing Recommendations**: What to test after resolution
5. **Prevention**: How to avoid similar conflicts in the future

Be practical and actionable.`, conflictFiles)

	return c.GenerateResponse(prompt)
}

// GenerateMergeCommitMessage generates a comprehensive merge commit message
func (c *Client) GenerateMergeCommitMessage(sourceBranch, targetBranch, changes string) (string, error) {
	prompt := fmt.Sprintf(`Generate a comprehensive merge commit message for merging '%s' into '%s'.

Changes being merged:
%s

Create a merge commit message that:
1. Clearly states what is being merged
2. Summarizes the key changes/features
3. Follows conventional commit format if appropriate
4. Mentions any important notes about the merge

Format as a proper merge commit message.`, sourceBranch, targetBranch, changes)

	return c.GenerateResponse(prompt)
}

// GenerateResponse sends a prompt to Solar LLM and returns the response
func (c *Client) GenerateResponse(prompt string) (string, error) {
	request := ChatRequest{
		Model: c.modelName,
		Messages: []Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Stream: false,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("error marshaling request: %v", err)
	}

	req, err := http.NewRequest("POST", c.baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var response ChatResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("error unmarshaling response: %v", err)
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no response choices returned")
	}

	content := response.Choices[0].Message.Content
	
	// Clean up the response by removing any <think>...</think> tags
	content = cleanResponse(content)
	
	return strings.TrimSpace(content), nil
}

// GenerateResponseStream sends a prompt to Solar LLM and returns the streaming response
func (c *Client) GenerateResponseStream(prompt string) (string, error) {
	request := ChatRequest{
		Model: c.modelName,
		Messages: []Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Stream: true,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("error marshaling request: %v", err)
	}

	req, err := http.NewRequest("POST", c.baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	// Start spinner while waiting for response
	spinner := NewSpinner()
	spinner.Start("Thinking...")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		spinner.Stop()
		return "", fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		spinner.Stop()
		body, _ := ioutil.ReadAll(resp.Body)
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var fullContent strings.Builder
	scanner := bufio.NewScanner(resp.Body)
	firstChunk := true
	
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" || !strings.HasPrefix(line, "data: ") {
			continue
		}
		
		// Remove "data: " prefix
		data := strings.TrimPrefix(line, "data: ")
		if data == "[DONE]" {
			break
		}
		
		var streamResp StreamResponse
		if err := json.Unmarshal([]byte(data), &streamResp); err != nil {
			continue // Skip invalid JSON lines
		}
		
		if len(streamResp.Choices) > 0 && streamResp.Choices[0].Delta.Content != "" {
			// Stop spinner on first content chunk and start printing
			if firstChunk {
				spinner.Stop()
				fmt.Print("Generated commit message: ")
				firstChunk = false
			}
			
			content := streamResp.Choices[0].Delta.Content
			fmt.Print(content) // Print streaming content immediately
			fullContent.WriteString(content)
		}
	}
	
	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error reading stream: %v", err)
	}
	
	fmt.Println() // Add newline after streaming
	
	finalContent := fullContent.String()
	// Clean up the response by removing any <think>...</think> tags
	finalContent = cleanResponse(finalContent)
	
	return strings.TrimSpace(finalContent), nil
}

// cleanResponse removes <think>...</think> blocks from the AI response.
func cleanResponse(content string) string {
	// Remove <think>...</think> blocks
	for {
		start := strings.Index(content, "<think>")
		if start == -1 {
			break
		}
		end := strings.Index(content[start:], "</think>")
		if end == -1 {
			break
		}
		end += start + len("</think>")
		content = content[:start] + content[end:]
	}
	
	return strings.TrimSpace(content)
} 