package openai

import (
	"context"
	"encoding/json"
	"github.com/iflytek/spark-ai-go/sparkai/llms"
	"github.com/iflytek/spark-ai-go/sparkai/messages"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestClient(t *testing.T, opts ...Option) llms.Model {
	t.Helper()
	if openaiKey := os.Getenv("OPENAI_API_KEY"); openaiKey == "" {
		t.Skip("OPENAI_API_KEY not set")
		return nil
	}

	llm, err := New(opts...)
	require.NoError(t, err)
	return llm
}

func TestMultiContentText(t *testing.T) {
	t.Parallel()
	llm := newTestClient(t)

	parts := []messages.ContentPart{
		messages.TextPart("I'm a pomeranian"),
		messages.TextPart("What kind of mammal am I?"),
	}
	content := []messages.MessageContent{
		{
			Role:  messages.ChatMessageTypeHuman,
			Parts: parts,
		},
	}

	rsp, err := llm.GenerateContent(context.Background(), content)
	require.NoError(t, err)

	assert.NotEmpty(t, rsp.Choices)
	c1 := rsp.Choices[0]
	assert.Regexp(t, "dog|canid", strings.ToLower(c1.Content))
}

func TestMultiContentTextChatSequence(t *testing.T) {
	t.Parallel()
	llm := newTestClient(t)

	content := []messages.MessageContent{
		{
			Role:  messages.ChatMessageTypeHuman,
			Parts: []messages.ContentPart{messages.TextPart("Name some countries")},
		},
		{
			Role:  messages.ChatMessageTypeAI,
			Parts: []messages.ContentPart{messages.TextPart("Spain and Lesotho")},
		},
		{
			Role:  messages.ChatMessageTypeHuman,
			Parts: []messages.ContentPart{messages.TextPart("Which if these is larger?")},
		},
	}

	rsp, err := llm.GenerateContent(context.Background(), content)
	require.NoError(t, err)

	assert.NotEmpty(t, rsp.Choices)
	c1 := rsp.Choices[0]
	assert.Regexp(t, "spain.*larger", strings.ToLower(c1.Content))
}

func TestMultiContentImage(t *testing.T) {
	t.Parallel()

	llm := newTestClient(t, WithModel("gpt-4-vision-preview"))

	parts := []messages.ContentPart{
		messages.ImageURLPart("https://github.com/tmc/langchaingo/blob/main/docs/static/img/parrot-icon.png?raw=true"),
		messages.TextPart("describe this image in detail"),
	}
	content := []messages.MessageContent{
		{
			Role:  messages.ChatMessageTypeHuman,
			Parts: parts,
		},
	}

	rsp, err := llm.GenerateContent(context.Background(), content, llms.WithMaxTokens(300))
	require.NoError(t, err)

	assert.NotEmpty(t, rsp.Choices)
	c1 := rsp.Choices[0]
	assert.Regexp(t, "parrot", strings.ToLower(c1.Content))
}

func TestWithStreaming(t *testing.T) {
	t.Parallel()
	llm := newTestClient(t)

	parts := []messages.ContentPart{
		messages.TextPart("I'm a pomeranian"),
		messages.TextPart("Tell me more about my taxonomy"),
	}
	content := []messages.MessageContent{
		{
			Role:  messages.ChatMessageTypeHuman,
			Parts: parts,
		},
	}

	var sb strings.Builder
	rsp, err := llm.GenerateContent(context.Background(), content,
		llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			sb.Write(chunk)
			return nil
		}))

	require.NoError(t, err)

	assert.NotEmpty(t, rsp.Choices)
	c1 := rsp.Choices[0]
	assert.Regexp(t, "dog|canid", strings.ToLower(c1.Content))
	assert.Regexp(t, "dog|canid", strings.ToLower(sb.String()))
}

//nolint:lll
func TestFunctionCall(t *testing.T) {
	t.Parallel()
	llm := newTestClient(t)

	parts := []messages.ContentPart{
		messages.TextPart("What is the weather like in Boston?"),
	}
	content := []messages.MessageContent{
		{
			Role:  messages.ChatMessageTypeHuman,
			Parts: parts,
		},
	}

	functions := []messages.FunctionDefinition{
		{
			Name:        "getCurrentWeather",
			Description: "Get the current weather in a given location",
			Parameters:  json.RawMessage(`{"type": "object", "properties": {"location": {"type": "string", "description": "The city and state, e.g. San Francisco, CA"}, "unit": {"type": "string", "enum": ["celsius", "fahrenheit"]}}, "required": ["location"]}`),
		},
	}

	rsp, err := llm.GenerateContent(context.Background(), content,
		llms.WithFunctions(functions))
	require.NoError(t, err)

	assert.NotEmpty(t, rsp.Choices)
	c1 := rsp.Choices[0]
	assert.Equal(t, "function_call", c1.StopReason)
	assert.NotNil(t, c1.FuncCall)
}

func showResponse(rsp any) string { //nolint:golint,unused
	b, err := json.MarshalIndent(rsp, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(b)
}
