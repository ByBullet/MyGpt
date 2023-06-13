package messages

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/sashabaranov/go-openai"
)

const AccessKey = ""
const AccessInterval = time.Second * 6

var client *openai.Client

func init() {
	client = openai.NewClient(AccessKey)
}

type Communication struct {
	messages       []openai.ChatCompletionMessage //访问过的信息
	lastAccessTime time.Time                      //上一次访问的时间
}

func (m *Communication) Chat(text string) (string, error) {
	m.messages = append(m.messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: text,
	})

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    openai.GPT3Dot5Turbo,
			Messages: m.messages,
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		m.messages = make([]openai.ChatCompletionMessage, 0, 10)
		return "", err
	}

	content := resp.Choices[0].Message.Content
	m.messages = append(m.messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleAssistant,
		Content: content,
	})
	return content, nil
}

/*
NewMessage 创建Message结构体
*/
func NewMessage() (result *Communication) {
	result = new(Communication)
	result.messages = make([]openai.ChatCompletionMessage, 0, 10)
	result.lastAccessTime = time.Now().Add(-AccessInterval)
	return result
}

type MessageStorage struct {
	sync.Map
}

func (ms *MessageStorage) putMessage(key string, m *Communication) {
	ms.Store(key, m)
}

func (ms *MessageStorage) GetMessage(k, msgItem string) (*Communication, error) {
	m, b := ms.Load(k)
	if !b {
		ms.putMessage(k, NewMessage())
		m, _ = ms.Load(k)
	}

	msg := m.(*Communication)
	lastAccess := msg.lastAccessTime.Add(AccessInterval)
	now := time.Now()
	if lastAccess.After(now) {
		fmt.Println(lastAccess.Format("2006-01-02 15:04:05"), "After", now.Format("2006-01-02 15:04:05"))
		return nil, fmt.Errorf("访问过于频繁")
	}
	msg.lastAccessTime = time.Now()
	msg.messages = append(msg.messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: msgItem,
	})
	return msg, nil
}
