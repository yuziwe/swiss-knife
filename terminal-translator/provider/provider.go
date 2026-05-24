package provider

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type LLMProvider interface {
    Completions(model string, messages []Message) (*Response, error)
}


