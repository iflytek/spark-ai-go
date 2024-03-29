package messages

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// ChatMessageType is the type of chat message.
type ChatMessageType string

// ErrUnexpectedChatMessageType is returned when a chat message is of an unexpected type.
var ErrUnexpectedChatMessageType = errors.New("unexpected chat message type")

const (
	// ChatMessageTypeAI is a message sent by an AI.
	ChatMessageTypeAI ChatMessageType = "ai"
	// ChatMessageTypeHuman is a message sent by a human.
	ChatMessageTypeHuman ChatMessageType = "human"
	// ChatMessageTypeSystem is a message sent by the system.
	ChatMessageTypeSystem ChatMessageType = "system"
	// ChatMessageTypeGeneric is a message sent by a generic user.
	ChatMessageTypeGeneric ChatMessageType = "generic"
	// ChatMessageTypeFunction is a message sent by a function.
	ChatMessageTypeFunction ChatMessageType = "function"
)

// ChatMessage represents a message in a chat.
type ChatMessage interface {
	// GetType gets the type of the message.
	GetType() ChatMessageType
	// GetContent gets the content of the message.
	GetContent() string

	UpdateContent(msg string)
}

// Named is an interface for objects that have a name.
type Named interface {
	GetName() string
}

// Statically assert that the types implement the interface.
var (
	_ ChatMessage = AIChatMessage{}
	_ ChatMessage = HumanChatMessage{}
	_ ChatMessage = SystemChatMessage{}
	_ ChatMessage = GenericChatMessage{}
	_ ChatMessage = FunctionChatMessage{}
)

// AIChatMessage is a message sent by an AI.
type AIChatMessage struct {
	// Content is the content of the message.
	Content string `json:"content"`

	// FunctionCall represents the model choosing to call a function.
	FunctionCall *FunctionCall `json:"function_call,omitempty"`
}

func (m AIChatMessage) UpdateContent(msg string) {
	//TODO implement me
	m.Content = msg
}

func (m AIChatMessage) GetType() ChatMessageType { return ChatMessageTypeAI }
func (m AIChatMessage) GetContent() string {
	return m.Content
}
func (m AIChatMessage) GetFunctionCall() *FunctionCall { return m.FunctionCall }

// HumanChatMessage is a message sent by a human.
type HumanChatMessage struct {
	Content string
}

func (m HumanChatMessage) UpdateContent(msg string) {
	//TODO implement me
	m.Content = msg

}

func (m HumanChatMessage) GetType() ChatMessageType { return ChatMessageTypeHuman }
func (m HumanChatMessage) GetContent() string       { return m.Content }

// SystemChatMessage is a chat message representing information that should be instructions to the AI system.
type SystemChatMessage struct {
	Content string
}

func (m SystemChatMessage) UpdateContent(msg string) {
	//TODO implement me
	m.Content = msg

}

func (m SystemChatMessage) GetType() ChatMessageType { return ChatMessageTypeSystem }
func (m SystemChatMessage) GetContent() string       { return m.Content }

// GenericChatMessage is a chat message with an arbitrary speaker.
type GenericChatMessage struct {
	Content string `json:"content"`
	Role    string `json:"role"`
	Name    string `json:"name"`
}

func (m GenericChatMessage) UpdateContent(msg string) {
	//TODO implement me
	m.Content = msg

}

func (m GenericChatMessage) GetType() ChatMessageType {
	return ChatMessageType(strings.ToLower(m.Role))
}
func (m GenericChatMessage) GetContent() string { return m.Content }
func (m GenericChatMessage) GetName() string    { return m.Name }

// FunctionChatMessage is a chat message representing the result of a function call.
type FunctionChatMessage struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

func (m FunctionChatMessage) UpdateContent(msg string) {
	//TODO implement me
	m.Content = msg

}

// FunctionCall is the name and arguments of a function call.
type FunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

func (f *FunctionCall) GetContent() string {
	b, err := json.Marshal(f)
	if err != nil {
		return err.Error()
	}
	return string(b)
}

func (m FunctionChatMessage) GetType() ChatMessageType { return ChatMessageTypeFunction }
func (m FunctionChatMessage) GetContent() string       { return m.Content }
func (m FunctionChatMessage) GetName() string          { return m.Name }

// GetBufferString gets the buffer string of messages.
func GetBufferString(messages []ChatMessage, humanPrefix string, aiPrefix string) (string, error) {
	result := []string{}
	for _, m := range messages {
		role, err := getMessageRole(m, humanPrefix, aiPrefix)
		if err != nil {
			return "", err
		}
		msg := fmt.Sprintf("%s: %s", role, m.GetContent())
		if m, ok := m.(AIChatMessage); ok && m.FunctionCall != nil {
			j, err := json.Marshal(m.FunctionCall)
			if err != nil {
				return "", err
			}
			msg = fmt.Sprintf("%s %s", msg, string(j))
		}
		result = append(result, msg)
	}
	return strings.Join(result, "\n"), nil
}

func getMessageRole(m ChatMessage, humanPrefix, aiPrefix string) (string, error) {
	var role string
	switch m.GetType() {
	case ChatMessageTypeHuman:
		role = humanPrefix
	case ChatMessageTypeAI:
		role = aiPrefix
	case ChatMessageTypeSystem:
		role = "System"
	case ChatMessageTypeGeneric:
		cgm, ok := m.(GenericChatMessage)
		if !ok {
			return "", fmt.Errorf("%w -%+v", ErrUnexpectedChatMessageType, m)
		}
		role = cgm.Role
	case ChatMessageTypeFunction:
		role = "Function"
	default:
		return "", ErrUnexpectedChatMessageType
	}
	return role, nil
}

// FunctionDefinition is a definition of a function that can be called by the model.
type FunctionDefinition struct {
	// Name is the name of the function.
	Name string `json:"name"`
	// Description is a description of the function.
	Description string `json:"description"`
	// Parameters is a list of parameters for the function.
	Parameters any `json:"parameters"`
}

// FunctionCallBehavior is the behavior to use when calling functions.
type FunctionCallBehavior string

const (
	// FunctionCallBehaviorNone will not call any functions.
	FunctionCallBehaviorNone FunctionCallBehavior = "none"
	// FunctionCallBehaviorAuto will call functions automatically.
	FunctionCallBehaviorAuto FunctionCallBehavior = "auto"
)
