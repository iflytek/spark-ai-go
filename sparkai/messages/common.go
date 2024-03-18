package messages

type Choice struct {
	FinishReason string
	Index        int
	Message      ChatMessage
}

type CompletionUsage struct {
	CompletionTokens float64
	PromptTokens     float64
	TotalTokens      float64
}

type ChatCompletionMessage struct {
	Id      string
	Choices SparkChoices
	Usage   CompletionUsage
}
type SparkHeader struct {
	Code    int
	Status  int
	Sid     string
	Message string
}
type SparkResponse struct {
	Header  SparkHeader
	Payload ChatCompletionMessage
}
type SparkChoices struct {
	Status int           `json:"status"`
	Seq    int           `json:"seq"`
	Text   []SparkChoice `json:"text"`
}
type SparkChoice struct {
	Index        int           `json:"index"`
	Role         string        `json:"role"`
	Content      string        `json:"content"`
	FunctionCall *FunctionCall `json:"function_call"`
}
