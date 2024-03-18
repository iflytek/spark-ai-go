package llms

import (
	"context"
	"errors"
	"github.com/iflytek/spark-ai-go/sparkai/messages"
)

// LLM is a Large Language Model.
type LLM interface {
	Call(ctx context.Context, prompt string, options ...CallOption) (string, error)
	Generate(ctx context.Context, prompts []string, options ...CallOption) ([]*Generation, error)
}

// ChatLLM is a LLM that can be used for chatting.
type ChatLLM interface {
	Call(ctx context.Context, messages []messages.ChatMessage, options ...CallOption) (*messages.AIChatMessage, error)
	Generate(ctx context.Context, messages [][]messages.ChatMessage, options ...CallOption) ([]*Generation, error)
}

// Model is an interface multi-modal models implement.
// Note: this is an experimental API.
type Model interface {
	// GenerateContent asks the model to generate content from parts.
	GenerateContent(ctx context.Context, messages []messages.MessageContent, options ...CallOption) (*messages.ContentResponse, error)
}

// GenerateFromSinglePrompt is a convenience function for calling an LLM with
// a single string prompt, expecting a single string response. It's useful for
// simple, string-only interactions and provides a slightly more ergonomic API
// than the more general [llms.Model.GenerateContent].
func GenerateFromSinglePrompt(ctx context.Context, llm Model, prompt string, options ...CallOption) (string, error) {
	msg := messages.MessageContent{
		Role:  messages.ChatMessageTypeHuman,
		Parts: []messages.ContentPart{messages.TextContent{prompt}},
	}

	resp, err := llm.GenerateContent(ctx, []messages.MessageContent{msg}, options...)
	if err != nil {
		return "", err
	}

	choices := resp.Choices
	if len(choices) < 1 {
		return "", errors.New("empty response from model")
	}
	c1 := choices[0]
	return c1.Content, nil
}

// Generation is a single generation from a LLM.
type Generation struct {
	// Text is the generated text.
	Text string `json:"text"`
	// Message stores the potentially generated message.
	Message *messages.AIChatMessage `json:"message"`
	// GenerationInfo is the generation info. This can contain vendor-specific information.
	GenerationInfo map[string]any `json:"generation_info"`
	// StopReason is the reason the generation stopped.
	StopReason string `json:"stop_reason"`
}

// LLMResult is the class that contains all relevant information for an LLM Result.
type LLMResult struct {
	Generations [][]*Generation
	LLMOutput   map[string]any
}
