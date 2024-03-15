package memory

import (
	"autospark-go/llms"
	"fmt"
)

// ChatMessageHistoryOption is a function for creating new chat message history
// with other than the default values.
type ChatMessageHistoryOption func(m *ChatMessageHistory)

// WithPreviousMessages is an option for NewChatMessageHistory for adding
// previous messages to the history.
func WithPreviousMessages(previousMessages []llms.ChatMessage) ChatMessageHistoryOption {
	return func(m *ChatMessageHistory) {
		m.messages = append(m.messages, previousMessages...)
	}
}

func WithPreviousMessagesByLog(logFile string) ChatMessageHistoryOption {

	return func(m *ChatMessageHistory) {
		storage, err := file_memory.NewChatHistoryFileStorage(logFile)
		if err != nil {
			return
		}
		his, err := storage.Read()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		m.messages = his
	}

}

func applyChatOptions(options ...ChatMessageHistoryOption) *ChatMessageHistory {
	h := &ChatMessageHistory{
		messages: make([]llms.ChatMessage, 0),
	}

	for _, option := range options {
		option(h)
	}

	return h
}
