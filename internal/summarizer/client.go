package summarizer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/sandepten/work-obsidian-noter/internal/notes"
)

// Client handles communication with the OpenCode server for AI summaries
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// NewClient creates a new OpenCode API client
func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: strings.TrimSuffix(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// Session represents an OpenCode session
type Session struct {
	ID string `json:"id"`
}

// TextPart represents a text message part
type TextPart struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// PromptRequest represents the request body for sending a message
type PromptRequest struct {
	Model *ModelSpec `json:"model,omitempty"`
	Parts []TextPart `json:"parts"`
}

// ModelSpec specifies which model to use
type ModelSpec struct {
	ProviderID string `json:"providerID"`
	ModelID    string `json:"modelID"`
}

// MessageResponse represents a response from the API
type MessageResponse struct {
	Info  MessageInfo `json:"info"`
	Parts []Part      `json:"parts"`
}

// MessageInfo contains message metadata
type MessageInfo struct {
	ID   string `json:"id"`
	Role string `json:"role"`
}

// Part represents a message part in the response
type Part struct {
	Type string `json:"type"`
	Text string `json:"text,omitempty"`
}

// createSession creates a new session for summarization
func (c *Client) createSession() (*Session, error) {
	req, err := http.NewRequest("POST", c.baseURL+"/session", bytes.NewBuffer([]byte("{}")))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to create session: status %d, body: %s", resp.StatusCode, string(body))
	}

	var session Session
	if err := json.NewDecoder(resp.Body).Decode(&session); err != nil {
		return nil, fmt.Errorf("failed to decode session response: %w", err)
	}

	return &session, nil
}

// sendMessage sends a message to a session and waits for response
func (c *Client) sendMessage(sessionID string, prompt string) (string, error) {
	requestBody := PromptRequest{
		Model: &ModelSpec{
			ProviderID: "google",
			ModelID:    "gemini-2.0-flash",
		},
		Parts: []TextPart{
			{Type: "text", Text: prompt},
		},
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("%s/session/%s/message", c.baseURL, sessionID)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send message: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to send message: status %d, body: %s", resp.StatusCode, string(body))
	}

	var msgResp MessageResponse
	if err := json.NewDecoder(resp.Body).Decode(&msgResp); err != nil {
		return "", fmt.Errorf("failed to decode message response: %w", err)
	}

	// Extract text from parts
	var result strings.Builder
	for _, part := range msgResp.Parts {
		if part.Type == "text" && part.Text != "" {
			result.WriteString(part.Text)
		}
	}

	return strings.TrimSpace(result.String()), nil
}

// SummarizeWorkItems generates an AI summary of completed work items
func (c *Client) SummarizeWorkItems(items []notes.WorkItem) (string, error) {
	if len(items) == 0 {
		return "No work items to summarize.", nil
	}

	// Build the prompt
	var sb strings.Builder
	sb.WriteString("Summarize the following completed work items in 1-2 concise sentences. Focus on the key accomplishments and outcomes. Keep it brief and professional:\n\n")

	for _, item := range items {
		sb.WriteString(fmt.Sprintf("- %s\n", item.Text))
	}

	// Create session
	session, err := c.createSession()
	if err != nil {
		return "", fmt.Errorf("failed to create session: %w", err)
	}

	// Send message and get response
	summary, err := c.sendMessage(session.ID, sb.String())
	if err != nil {
		return "", fmt.Errorf("failed to get summary: %w", err)
	}

	return summary, nil
}

// TestConnection tests if the OpenCode server is reachable
func (c *Client) TestConnection() error {
	req, err := http.NewRequest("GET", c.baseURL+"/global/health", nil)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to connect to OpenCode server: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("OpenCode server returned status %d", resp.StatusCode)
	}

	return nil
}
