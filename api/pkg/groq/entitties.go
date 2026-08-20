package groq

import "net/http"

type Client struct {
	apiKey            string
	chatCompletionURL string
	httpClient        *http.Client
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Option func(*requestBody)

type ClientOption func(*Client)

// func WithAPIKey(apiKey string) ClientOption {
// 	return func(c *Client) {
// 		c.apiKey = apiKey
// 	}
// }

type ResponseFormat struct {
	Type string `json:"type,omitempty"`
}

type requestBody struct {
	Messages        []Message       `json:"messages"`
	Model           string          `json:"model"`
	MaxTokens       int             `json:"max_tokens"`
	ResponseFormat  *ResponseFormat `json:"response_format,omitempty"`
	ReasoningFormat string          `json:"reasoning_format,omitempty"`
	Seed            int             `json:"seed,omitempty"`
	Stream          bool            `json:"stream"`
	Stop            *string         `json:"stop,omitempty"`
	Temperature     float64         `json:"temperature"`
	TopP            float64         `json:"top_p"`
}

type ChatCompletionResponse struct {
	ID      string `json:"id,omitempty"`
	Object  string `json:"object,omitempty"`
	Created int    `json:"created,omitempty"`
	Model   string `json:"model,omitempty"`
	Choices []struct {
		Index        int         `json:"index,omitempty"`
		Message      Message     `json:"message,omitempty"`
		Logprobs     interface{} `json:"logprobs,omitempty"`
		FinishReason string      `json:"finish_reason,omitempty"`
	} `json:"choices,omitempty"`
	Usage struct {
		QueueTime        float64 `json:"queue_time,omitempty"`
		PromptTokens     int     `json:"prompt_tokens,omitempty"`
		PromptTime       float64 `json:"prompt_time,omitempty"`
		CompletionTokens int     `json:"completion_tokens,omitempty"`
		CompletionTime   float64 `json:"completion_time,omitempty"`
		TotalTokens      int     `json:"total_tokens,omitempty"`
		TotalTime        float64 `json:"total_time,omitempty"`
	} `json:"usage,omitempty"`
	SystemFingerprint string `json:"system_fingerprint,omitempty"`
	XGroq             struct {
		ID string `json:"id,omitempty"`
	} `json:"x_groq,omitempty"`
}
