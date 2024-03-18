package spark

import (
	"context"
	"github.com/iflytek/spark-ai-go/sparkai/callbacks"
	"github.com/iflytek/spark-ai-go/sparkai/llms"
	"github.com/iflytek/spark-ai-go/sparkai/llms/spark/internal/sparkclient"
	"github.com/iflytek/spark-ai-go/sparkai/messages"
)

type LLM struct {
	CallbacksHandler callbacks.Handler
	client           *sparkclient.Client
}

// New returns a new Spark LLM.
func New(opts ...Option) (*LLM, error) {
	opt, c, err := NewClient(opts...)
	if err != nil {
		return nil, err
	}
	return &LLM{
		client:           c,
		CallbacksHandler: opt.callbackHandler,
	}, err
}

// Call requests a completion for the given prompt.
func (o *LLM) Call(ctx context.Context, prompt string, options ...llms.CallOption) (string, error) {
	r, err := o.Generate(ctx, []string{prompt}, options...)
	if err != nil {
		return "", err
	}
	if len(r) == 0 {
		return "", ErrEmptyResponse
	}
	return r[0].Text, nil
}

func (o *LLM) Generate(ctx context.Context, prompts []string, options ...llms.CallOption) ([]*llms.Generation, error) {
	if o.CallbacksHandler != nil {
		o.CallbacksHandler.HandleLLMStart(ctx, prompts)
	}

	opts := llms.CallOptions{}
	for _, opt := range options {
		opt(&opts)
	}

	generations := make([]*llms.Generation, 0, len(prompts))
	for _, prompt := range prompts {
		result, err := o.client.CreateCompletion(ctx, &sparkclient.CompletionRequest{
			Prompt:      prompt,
			MaxTokens:   opts.MaxTokens,
			Temperature: opts.Temperature,
			N:           opts.N,
			TopK:        int64(opts.TopK),
			Functions:   opts.Functions,
		})
		if err != nil {
			if o.CallbacksHandler != nil {
				o.CallbacksHandler.HandleLLMError(ctx, err)
			}
			return nil, err
		}
		chatRes := result.(*sparkclient.ChatResponse)

		generations = append(generations, &llms.Generation{
			Text: result.GetContent(),
			Message: &messages.AIChatMessage{
				Content:      chatRes.GetContent(),
				FunctionCall: chatRes.FunctionCall,
			},
		})
	}

	//if o.CallbacksHandler != nil {
	//	o.CallbacksHandler.HandleLLMEnd(ctx, llms.LLMResult{Generations: [][]*llms.Generation{generations}})
	//}

	return generations, nil
}

//// CreateEmbedding creates embeddings for the given input texts.
//func (o *LLM) CreateEmbedding(ctx context.Context, inputTexts []string) ([][]float32, error) {
//	embeddings, err := o.client.CreateEmbedding(ctx, &sparkclient.EmbeddingRequest{
//		Input: inputTexts,
//		Model: o.client.Model,
//	})
//	if err != nil {
//		return nil, err
//	}
//	if len(embeddings) == 0 {
//		return nil, ErrEmptyResponse
//	}
//	if len(inputTexts) != len(embeddings) {
//		return embeddings, ErrUnexpectedResponseLength
//	}
//	return embeddings, nil
//}
