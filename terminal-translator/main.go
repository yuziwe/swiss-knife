package main

import (
	"fmt"
	"os"
	"io"
	"net/http"
	"encoding/json"
	"strings"
	tsize "github.com/kopoli/go-terminal-size" // For getting current terminal size
)

const (
	ColorRed    = "\x1b[0;31m"
	ColorGreen  = "\x1b[0;32m"
	ColorYellow = "\x1b[0;33m"
	ColorBlue   = "\x1b[0;34m"
	ColorReset  = "\x1b[0m"
)

const (
	API_URL_KEY = "TERMINAL_TRANSLATOR_API_URL"
	API_KEY_KEY = "TERMINAL_TRANSLATOR_API_KEY"
	SYSTEM_PROMPT = 
	`You are a professional translation assistant that **exclusively performs precise text translation tasks** and strictly adheres to the following rules: 1. **Function Definition**  - Translate Chinese text input into English.  - Translate English text input into Chinese.  - Automatically detect the language of the input text.  2. **Output Rules**  - **Output only the translated text in the target language**, with **no** prefixes, suffixes, explanations, notes, punctuation clarifications, or formatting embellishments.  - Absolutely **do not** output lead-ins such as "Translation:", "Result:", or similar.  - Absolutely **do not** output any characters or line breaks that are not part of the translation itself.  3. **Examples**  - User Input: "你好，世界" → Your Output: "Hello, world"  - User Input: "How are you" → Your Output: "你好吗" Strictly follow these rules to ensure every response contains only the pure translation result.`
)

// ===========Request filed===========
type NilObject struct {}

type Message struct{
	Role	string`json:"role"`
	Content string`json:"content"`
}

type Completion struct {
	Model				string		`json:"model"`
	Messages 			[]Message	`json:"messages"`
	Temperature 		float32		`json:"temperature"`
	TopP				float32		`json:"top_p"`
	N					int			`json:"n"`
	Stream				bool		`json:"stream"`
	StreamOptions 		NilObject	`json:"stream_options"`
	Stop				[]string	`json:"stop"`
	MaxTokens			int			`json:"max_tokens"`
	MaxCompletionTokens int			`json:"max_completion_tokens"`
	PresendPenalty 		int			`json:"presence_penalty"`
	FrequencyPenalty 	int			`json:"frequency_penalty"`
	LogitBias			NilObject	`json:"logit_bias"`
	User				string		`json:"user"`
	Tools				[]NilObject	`json:"tools"`
	ResponseFormat		NilObject	`json:"response_formata"`
	Seed				int			`json:"seed"`
	ReasoningEffort		string		`json:"reasoning_effort"`
	Modalities			[]string	`json:"modalities"`
	Audio				NilObject	`json:"audio"`
}
// ===========Request filed===========

type OpenAI struct {
	BaseUrl		string
	ApiKey		string
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
	Index              int     `json:"index"`
	RMessage           RMessage `json:"message"`
	FinishReason       string  `json:"finish_reason"`
	NativeFinishReason string  `json:"native_finish_reason"`
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
// ===========Response filed===========

func (m *OpenAI)completions(model string, messages []Message) (*Response, error) {
	// New http client
	client := &http.Client{}

	// Fill http request body
	completion := &Completion{
		Model: 			 model,
		Messages: 		 messages,
		Temperature: 	 0.5,
		MaxTokens:		 8192,
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

func generate_separator() {
	s, err := tsize.GetSize()
	if err != nil {
		fmt.Println("ERROR: get window size failed!: ", err)
		os.Exit(1)
	}

	for i := 0; i < s.Width; i++ {
		fmt.Print("=")
	}

	fmt.Println()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage:", os.Args[0], " <English|Chinese> <Model>")
		os.Exit(1)
	}

	model := "gpt-5.1"
	if len(os.Args) > 2 {
		model = os.Args[2]
	}

	generate_separator()
	fmt.Println(ColorRed, "Current model: ", model, ColorReset)
	generate_separator()

	// Read API_URL from environment
	base_api_url := os.Getenv(API_URL_KEY);
	if base_api_url == "" {
		fmt.Println("ERROR: ", API_URL_KEY , " is empty!")
		os.Exit(1)
	}

	// Read API_KEY from environment
	api_key := os.Getenv(API_KEY_KEY)
	if api_key == "" {
		fmt.Println("ERROR: ", API_KEY_KEY, " is empty!")
		os.Exit(1)
	}

	client := &OpenAI{
		BaseUrl:  base_api_url,
		ApiKey:   api_key,
	}

	// Messages
	msgs := []Message{
		{ Role: "system", Content: SYSTEM_PROMPT },
		{ Role: "user"  , Content: os.Args[1] },
	}

	resp, err := client.completions(model, msgs)
	if err != nil {
		fmt.Println("ERROR: create completions failed!: ", err)
		os.Exit(1)
	}

	if len(resp.Choices) == 0 {
		fmt.Println("ERROR: got empty response!")
		os.Exit(1)
	}

	// Output
	fmt.Println(ColorYellow, os.Args[1], ColorReset)
	generate_separator()
	fmt.Println(ColorGreen, resp.Choices[0].RMessage.Content, ColorReset)
	generate_separator()
}

