package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type OpenAIProvider struct {
	BaseURL string
	APIKey  string
	Model   string
}

func NewOpenAIProvider(baseURL, apiKey, model string) *OpenAIProvider {
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}
	return &OpenAIProvider{
		BaseURL: baseURL,
		APIKey:  apiKey,
		Model:   model,
	}
}

type openAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openAIRequest struct {
	Model    string          `json:"model"`
	Messages []openAIMessage `json:"messages"`
}

type openAIChoice struct {
	Message openAIMessage `json:"message"`
}

type openAIResponse struct {
	Choices []openAIChoice `json:"choices"`
}

func (p *OpenAIProvider) GenerateFix(ctx context.Context, errorLog string, codeContext string) (*FixResponse, error) {
	prompt := fmt.Sprintf(`You are an expert software engineer. Analyze the following error log and code context to provide a fix.
Return your response in strict JSON format with the following fields:
"explanation": a brief description of the fix,
"file_path": the relative path to the file that needs changing,
"original_code": the original code snippet being replaced,
"fixed_code": the new code snippet,
"git_patch": a standard git diff patch for the change.

Error Log:
%s

Code Context:
%s`, errorLog, codeContext)

	reqBody := openAIRequest{
		Model: p.Model,
		Messages: []openAIMessage{
			{Role: "system", Content: "You are a helpful assistant that provides bug fixes in JSON format."},
			{Role: "user", Content: prompt},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", p.BaseURL+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	if p.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+p.APIKey)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("openai api error (status %d): %s", resp.StatusCode, string(body))
	}

	var res openAIResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	if len(res.Choices) == 0 {
		return nil, fmt.Errorf("no response from openai")
	}

	var fix FixResponse
	content := res.Choices[0].Message.Content
	// Simple cleanup in case LLM adds markdown blocks
	if start := bytes.Index( []byte(content), []byte("```json")); start != -1 {
		content = content[start+7:]
		if end := bytes.LastIndex( []byte(content), []byte("```")); end != -1 {
			content = content[:end]
		}
	}

	if err := json.Unmarshal([]byte(content), &fix); err != nil {
		return nil, fmt.Errorf("failed to parse AI response: %w. Raw: %s", err, content)
	}

	return &fix, nil
}
