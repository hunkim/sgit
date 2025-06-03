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

Your task is to analyze the changes and UNDERSTAND THE DEVELOPER'S INTENTION, not just describe what changed.

=== GIT DIFF ===
%s

=== CURRENT BRANCH ===
%s

=== RECENT COMMITS (last 5) ===
%s

=== FILES CHANGED ===
%s

INTENTION ANALYSIS - Consider:
1. **Purpose**: Why was this change made? Look for clues in code patterns, branch name, recent commits
2. **Context**: How does this fit into the larger development story?
3. **Impact**: What problem does this solve or what improvement does it provide?

Generate a commit message that:
1. Follows conventional commit format: type(scope): description
2. CAPTURES THE INTENTION and purpose, not just the mechanics
3. Uses imperative mood ("add" not "added")
4. Includes a brief body explaining the WHY and impact
5. Keep total length between 200-400 characters

Format:
type(scope): intention-focused summary explaining WHY

Brief explanation of the purpose and impact of this change.
Focus on the problem solved or improvement made.

BREAKING CHANGE: description if applicable (only if truly breaking)

Respond with only the commit message, no explanations.`, diff, branch, recentCommits, fileList)

	return c.GenerateResponse(prompt)
}

// GenerateComprehensiveCommitMessageStream generates a commit message with streaming
func (c *Client) GenerateComprehensiveCommitMessageStream(diff, branch, recentCommits, fileList string) (string, error) {
	prompt := fmt.Sprintf(`You are an expert software developer who writes excellent commit messages following the Conventional Commits specification.

Your task is to analyze the changes and UNDERSTAND THE DEVELOPER'S INTENTION, not just describe what changed.

=== GIT DIFF ===
%s

=== CURRENT BRANCH ===
%s

=== RECENT COMMITS (last 5) ===
%s

=== FILES CHANGED ===
%s

INTENTION ANALYSIS - Consider these aspects:
1. **Purpose**: Why was this change made? (bug fix, new feature, improvement, refactor, etc.)
2. **Context Clues**: 
   - Branch name patterns (feature/, fix/, hotfix/, etc.)
   - File patterns (test files = testing, config files = configuration, etc.)
   - Code patterns (adding validation = security/reliability, adding logs = debugging, etc.)
3. **Development Flow**: 
   - How does this fit with recent commits?
   - Is this part of a larger feature or fix?
   - Is this completing something started earlier?
4. **Impact Intent**:
   - Performance improvement? Security enhancement? User experience? Developer experience?
   - Breaking changes? Backward compatibility? API changes?
5. **Technical Intention**:
   - Architecture improvements? Code quality? Maintainability?
   - Integration with external systems? Internal refactoring?

REASONING PATTERNS TO LOOK FOR:
- Adding tests → ensuring reliability/quality
- Adding error handling → improving robustness  
- Adding validation → security/data integrity
- Adding logging → debugging/monitoring
- Changing config → deployment/environment setup
- Updating docs → knowledge sharing/onboarding
- Refactoring → code quality/maintainability
- Adding endpoints → new functionality
- Fixing types → type safety/correctness
- Adding dependencies → leveraging external capabilities

Generate a commit message that:
1. Follows conventional commit format: type(scope): description
2. Types: feat, fix, docs, style, refactor, test, chore, perf, ci, build
3. CAPTURES THE INTENTION, not just the mechanics
4. Uses imperative mood ("add" not "added")
5. Includes a brief body (2-3 lines) explaining:
   - WHY this change was made (the intention/purpose)
   - WHAT problem it solves or improvement it provides
   - HOW it impacts users/developers/system
6. Mentions breaking changes if applicable
7. Keep total length between 200-400 characters

Examples of intention-focused messages:
❌ "feat(api): add new endpoint" (describes mechanics)
✅ "feat(api): enable user profile customization" (describes intention)

❌ "fix(db): change query" (describes mechanics)  
✅ "fix(db): prevent memory leak in long-running queries" (describes intention)

❌ "refactor(auth): update code" (describes mechanics)
✅ "refactor(auth): simplify token validation for better maintainability" (describes intention)

Format:
type(scope): intention-focused summary that explains WHY

Brief explanation of the purpose and impact of this change.
Focus on the problem solved or improvement made, not just what files changed.

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

// AnalyzeLogStream generates insights from the git log with streaming
func (c *Client) AnalyzeLogStream(logOutput, timeframe string) (string, error) {
	prompt := fmt.Sprintf(`Analyze the following git log (%s) and provide detailed insights:

%s

DEVELOPMENT ANALYSIS - Provide comprehensive insights:

1. **📊 Activity Summary**: 
   - Overall development velocity and patterns
   - Peak activity periods and quiet phases
   - Commit frequency and distribution

2. **🚀 Key Features & Improvements**:
   - Major features implemented
   - Significant improvements made
   - New capabilities added

3. **🐛 Bug Fixes & Maintenance**:
   - Critical fixes applied
   - Performance improvements
   - Security enhancements

4. **👥 Contributor Insights**:
   - Active contributors and their focus areas
   - Collaboration patterns
   - Expertise distribution

5. **🔍 Development Patterns**:
   - Coding practices and conventions
   - Testing and documentation habits
   - Release and deployment patterns

6. **💡 Recommendations**:
   - Areas for improvement
   - Suggested next steps
   - Technical debt considerations

Be insightful and actionable. Focus on trends, patterns, and meaningful observations.`, timeframe, logOutput)

	return c.GenerateResponseStream(prompt)
}

// SummarizeDiffStream generates a summary of the git diff with streaming
func (c *Client) SummarizeDiffStream(diff string) (string, error) {
	prompt := fmt.Sprintf(`Analyze the following git diff and provide a comprehensive, structured summary:

%s

CHANGE ANALYSIS - Provide detailed insights:

1. **📋 Summary**: 
   - High-level overview of what changed
   - Primary purpose and intention of changes

2. **📁 Files & Components**:
   - Main files modified, added, or removed
   - Components and modules affected
   - Architecture areas impacted

3. **🔄 Type of Changes**:
   - New features implemented
   - Bug fixes applied  
   - Refactoring and improvements
   - Configuration or documentation updates

4. **⚡ Impact Assessment**:
   - Functional changes and new capabilities
   - Performance implications
   - User experience impacts
   - Developer experience changes

5. **🎯 Technical Details**:
   - Key algorithms or logic changes
   - API modifications
   - Database or schema changes
   - Dependencies added or updated

6. **⚠️ Important Notes**:
   - Breaking changes (if any)
   - Migration requirements
   - Testing considerations
   - Deployment implications

Be thorough yet concise. Focus on what matters most for understanding the change.`, diff)

	return c.GenerateResponseStream(prompt)
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