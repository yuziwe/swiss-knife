package provider

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type nilObject struct{}

type OpenAI struct {
	BaseUrl    string
	ApiKey     string
	HttpClient *http.Client
}

// ===========Request filed===========
type completion struct {
	Model               string      `json:"model"`
	Messages            []Message   `json:"messages"`
	Temperature         float32     `json:"temperature"`
	TopP                float32     `json:"top_p"`
	N                   int         `json:"n"`
	Stream              bool        `json:"stream"`
	StreamOptions       any         `json:"stream_options"`
	Stop                []string    `json:"stop"`
	MaxTokens           int         `json:"max_tokens"`
	MaxCompletionTokens int         `json:"max_completion_tokens"`
	PresendPenalty      int         `json:"presence_penalty"`
	FrequencyPenalty    int         `json:"frequency_penalty"`
	LogitBias           nilObject   `json:"logit_bias"`
	User                string      `json:"user"`
	Tools               []nilObject `json:"tools"`
	ResponseFormat      nilObject   `json:"response_formata"`
	Seed                int         `json:"seed"`
	ReasoningEffort     string      `json:"reasoning_effort"`
	Modalities          []string    `json:"modalities"`
	Audio               nilObject   `json:"audio"`
}

// ===========Response filed===========
type response struct {
	ID      string    `json:"id"`
	Model   string    `json:"model"`
	Usage   usage     `json:"usage"`
	Object  string    `json:"object"`
	Created int       `json:"created"`
	Choices []choices `json:"choices"`
}

type choices struct {
	Index              int         `json:"index"`
	RMessage           respMessage `json:"message"`
	FinishReason       string      `json:"finish_reason"`
	NativeFinishReason string      `json:"native_finish_reason"`
}

type respMessage struct {
	Role             string `json:"role"`
	Content          string `json:"content"`
	ReasoningContent any    `json:"reasoning_content"`
	ToolCalls        any    `json:"tool_calls"`
}

type usage struct {
	CompletionTokens        int                     `json:"completion_tokens"`
	TotalTokens             int                     `json:"total_tokens"`
	PromptTokens            int                     `json:"prompt_tokens"`
	CompletionTokensDetails completionTokensDetails `json:"completion_tokens_details"`
}

type completionTokensDetails struct {
	ReasoningTokens int `json:"reasoning_tokens"`
}

func (m *OpenAI) Completions(model string, messages []Message) (*TTResponse, error) {
	// Fill http request body
	req_completion := &completion{
		Model:           model,
		Messages:        messages,
		Temperature:     0.2,
		TopP:            0.9,
		N:               1,
		Stream:          false,
		StreamOptions:   nil,
		MaxTokens:       8192,
		ReasoningEffort: "low",
	}

	// Serialize
	req_body, err := json.Marshal(req_completion)
	if err != nil {
		return nil, fmt.Errorf("ERROR: json serialize failed: %v", err)
	}

	// Create new request
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/chat/completions", m.BaseUrl), strings.NewReader(string(req_body)))
	if err != nil {
		return nil, fmt.Errorf("ERROR: create https request failed!: %v", err)
	}

	// Add Header
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", m.ApiKey))

	// Do Request
	resp, err := m.HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ERROR: https request failed!: %v", err)
	}
	defer resp.Body.Close()

	if err := preprocess(resp); err != nil {
		return nil, err
	}

	resp_bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ERROR: read response body failed!: %v", err)
	}

	// Parse response
	resp_body := &response{}
	if err := json.Unmarshal(resp_bytes, resp_body); err != nil {
		return nil, fmt.Errorf("ERROR: json unserialize failed!: %v", err)
	}

	if len(resp_body.Choices) == 0 {
		return nil, fmt.Errorf("ERROR: got empty response!")
	}

	return &TTResponse{
		ID:      resp_body.ID,
		Model:   resp_body.Model,
		Content: resp_body.Choices[0].RMessage.Content,
	}, nil
}
