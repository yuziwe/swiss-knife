package provider

import (
    "fmt"
	"net/http"
	"encoding/json"
	"io"
    "strings"
)

type NilObject struct{}

type OpenAI struct {
	BaseUrl string
	ApiKey  string
}

// ===========Request filed===========
type Completion struct {
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
	LogitBias           NilObject   `json:"logit_bias"`
	User                string      `json:"user"`
	Tools               []NilObject `json:"tools"`
	ResponseFormat      NilObject   `json:"response_formata"`
	Seed                int         `json:"seed"`
	ReasoningEffort     string      `json:"reasoning_effort"`
	Modalities          []string    `json:"modalities"`
	Audio               NilObject   `json:"audio"`
}

// ===========Response filed===========
type Response struct {
	ID      string    `json:"id"`
	Object  string    `json:"object"`
	Created int       `json:"created"`
	Model   string    `json:"model"`
	Choices []Choices `json:"choices"`
	Usage   Usage     `json:"usage"`
}

type RMessage struct {
	Role             string `json:"role"`
	Content          string `json:"content"`
	ReasoningContent any    `json:"reasoning_content"`
	ToolCalls        any    `json:"tool_calls"`
}

type Choices struct {
	Index              int      `json:"index"`
	RMessage           RMessage `json:"message"`
	FinishReason       string   `json:"finish_reason"`
	NativeFinishReason string   `json:"native_finish_reason"`
}

type CompletionTokensDetails struct {
	ReasoningTokens int `json:"reasoning_tokens"`
}

type Usage struct {
	CompletionTokens        int                     `json:"completion_tokens"`
	TotalTokens             int                     `json:"total_tokens"`
	PromptTokens            int                     `json:"prompt_tokens"`
	CompletionTokensDetails CompletionTokensDetails `json:"completion_tokens_details"`
}

func (m *OpenAI) Completions(model string, messages []Message) (*Response, error) {
	// New http client
	client := &http.Client{}

	// Fill http request body
	completion := &Completion{
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
	req_body, err := json.Marshal(completion)
	if err != nil {
		fmt.Println("ERROR: json serialize failed: ", err)
		return nil, nil
	}

	// Create new request
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/chat/completions", m.BaseUrl), strings.NewReader(string(req_body)))
	if err != nil {
		fmt.Println("ERROR: create https request failed!")
		return nil, nil
	}

	// Add Header
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", m.ApiKey))

	// Do Request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("ERROR: https request failed!: ", err)
		return nil, nil
	}
	defer resp.Body.Close()

	resp_bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("ERROR: read response body failed!: ", err)
		return nil, nil
	}

	// Parse response
	resp_body := &Response{}
	if err := json.Unmarshal(resp_bytes, resp_body); err != nil {
		fmt.Println("ERROR: json unserialize failed!: ", err)
		return nil, nil
	}

	return resp_body, nil
}
