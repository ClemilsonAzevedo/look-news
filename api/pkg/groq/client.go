package groq

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const defaultChatCompletionURL = "https://api.groq.com/openai/v1/chat/completions"

func NewClient(options ...ClientOption) *Client {
	client := &Client{
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
		chatCompletionURL: defaultChatCompletionURL,
		apiKey:            os.Getenv("GROQ_API_KEY"),
	}

	for _, option := range options {
		option(client)
	}

	return client
}

func (c *Client) ChatCompletion(
	messages []Message,
	options ...Option,
) (*ChatCompletionResponse, error) {
	return c.ChatCompletionContext(context.Background(), messages, options...)
}

func (c *Client) ChatCompletionContext(
	ctx context.Context,
	messages []Message,
	options ...Option,
) (*ChatCompletionResponse, error) {
	body := requestBody{
		Messages:        filterMessages(messages),
		Model:           "qwen/qwen3.6-27b",
		Temperature:     1,
		MaxTokens:       4096,
		TopP:            1,
		Stream:          false,
		Stop:            nil,
		ReasoningFormat: "parsed",
	}

	for _, option := range options {
		if option != nil {
			option(&body)
		}
	}

	jsonData, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.chatCompletionURL,
		bytes.NewReader(jsonData),
	)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		errorBody, readErr := io.ReadAll(io.LimitReader(resp.Body, 1<<20))

		if readErr != nil {
			return nil, fmt.Errorf(
				"groq returned status %d; could not read error body: %w",
				resp.StatusCode,
				readErr,
			)
		}

		return nil, fmt.Errorf(
			"groq returned status %d: %s",
			resp.StatusCode,
			strings.TrimSpace(string(errorBody)),
		)
	}

	var completion ChatCompletionResponse

	if err := json.NewDecoder(resp.Body).Decode(&completion); err != nil {
		return nil, fmt.Errorf("decode response body: %w", err)
	}

	return &completion, nil
}

func filterMessages(messages []Message) []Message {
	filteredMessages := make([]Message, 0, len(messages))

	for _, message := range messages {
		if strings.TrimSpace(message.Content) != "" {
			filteredMessages = append(filteredMessages, message)
		}
	}

	return filteredMessages
}
