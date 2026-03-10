package chat

import (
	"github.com/Tencent/WeKnora/internal/models/provider"
	"github.com/sashabaranov/go-openai"
)

// QwenChat Aliyun Qwen 모델 채팅 구현
// Qwen3 모델은 enable_thinking 파라미터에 대해 특별 처리 필요
type QwenChat struct {
	*RemoteAPIChat
}

// QwenChatCompletionRequest Qwen 모델용 커스텀 요청 구조체
type QwenChatCompletionRequest struct {
	openai.ChatCompletionRequest
	EnableThinking *bool `json:"enable_thinking,omitempty"`
}

// NewQwenChat Qwen 채팅 인스턴스 생성
func NewQwenChat(config *ChatConfig) (*QwenChat, error) {
	config.Provider = string(provider.ProviderAliyun)

	remoteChat, err := NewRemoteAPIChat(config)
	if err != nil {
		return nil, err
	}

	chat := &QwenChat{
		RemoteAPIChat: remoteChat,
	}

	// 요청 커스터마이저 설정
	remoteChat.SetRequestCustomizer(chat.customizeRequest)

	return chat, nil
}

// isQwen3Model Qwen3 모델 여부 확인
func (c *QwenChat) isQwen3Model() bool {
	return provider.IsQwen3Model(c.GetModelName())
}

// customizeRequest Qwen 요청 커스터마이즈
func (c *QwenChat) customizeRequest(req *openai.ChatCompletionRequest, opts *ChatOptions, isStream bool) (any, bool) {
	// Qwen3 모델만 특별 처리
	if !c.isQwen3Model() {
		return nil, false
	}

	// 비스트리밍 요청은 thinking을 명시적으로 비활성화
	if !isStream {
		qwenReq := QwenChatCompletionRequest{
			ChatCompletionRequest: *req,
		}
		enableThinking := false
		qwenReq.EnableThinking = &enableThinking
		return qwenReq, true
	}

	return nil, false
}
