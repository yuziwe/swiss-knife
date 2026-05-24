package provider

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type TTResponse struct {
	ID      string
	Model   string
	Content string
}

const (
	OK                    = 200
	InvalidFormat         = 400
	InvalidAuthentication = 401
	InsufficientBalance   = 402
	PermissionDenied      = 403
	ResourceNotFound      = 404
	RequestTooLarge       = 413
	InvalidParameter      = 422
	ReachedLimit          = 429
	InternalError         = 500
	ServerOverloaded      = 503
	TimeoutError          = 504
	OverloadedError       = 529 // Only use in Anthropic
)

type errorMsg struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

type errResponse struct {
	Error errorMsg `json:"error"`
}

type LLMProvider interface {
	Completions(model string, messages []Message) (*TTResponse, error)
}

func preprocess(resp *http.Response) error {
	if resp.StatusCode != OK {
		resp_bytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("ERROR: read response body failed!: %v", err)
		}

		if len(resp_bytes) == 0 {
			return fmt.Errorf("ERROR: Got empty response")
		}

		erresp := &errResponse{}
		if err := json.Unmarshal(resp_bytes, erresp); err != nil {
			return fmt.Errorf("ERROR: json unserialize for {%s} failed!: %v", string(resp_bytes), err)
		}

		return fmt.Errorf("ERROR: %s", erresp.Error.Message)
	}

	return nil
}
