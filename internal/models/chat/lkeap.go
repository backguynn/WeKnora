package chat

import (
	"strings"

	"github.com/Tencent/WeKnora/internal/models/provider"
	"github.com/sashabaranov/go-openai"
)

// LKEAPChat Tencent Cloud Knowledge Engine Atomic Ability (LKEAP) 채팅 구현
// DeepSeek-R1, DeepSeek-V3 계열 모델 지원, chain-of-thought 기능 제공
// 참고: https://cloud.tencent.com/document/product/1772/115963
//
// 표준 OpenAI API와의 차이:
// 1. thinking 파라미터 형식이 다름: LKEAP는 {"type": "enabled"/"disabled"} 사용
// 2. DeepSeek V3.x 계열만 thinking 파라미터를 명시적으로 설정 필요, R1 계열은 기본 활성
type LKEAPChat struct {
	*RemoteAPIChat
}

// LKEAPThinkingConfig chain-of-thought 설정(LKEAP 전용 형식)
type LKEAPThinkingConfig struct {
	Type string `json:"type"` // "enabled" 또는 "disabled"
}

// LKEAPChatCompletionRequest LKEAP 커스텀 요청 구조체
type LKEAPChatCompletionRequest struct {
	openai.ChatCompletionRequest
	Thinking *LKEAPThinkingConfig `json:"thinking,omitempty"` // chain-of-thought 토글(V3.x 계열만)
}

// NewLKEAPChat LKEAP 채팅 인스턴스 생성
func NewLKEAPChat(config *ChatConfig) (*LKEAPChat, error) {
	// provider 설정 확인
	config.Provider = string(provider.ProviderLKEAP)

	remoteChat, err := NewRemoteAPIChat(config)
	if err != nil {
		return nil, err
	}

	chat := &LKEAPChat{
		RemoteAPIChat: remoteChat,
	}

	// 요청 커스터마이저 설정, LKEAP 전용 thinking 파라미터 추가
	remoteChat.SetRequestCustomizer(chat.customizeRequest)

	return chat, nil
}

// isDeepSeekV3Model DeepSeek V3.x 계열 모델 여부 확인
func (c *LKEAPChat) isDeepSeekV3Model() bool {
	return strings.Contains(strings.ToLower(c.GetModelName()), "deepseek-v3")
}

// customizeRequest LKEAP 요청 커스터마이즈
func (c *LKEAPChat) customizeRequest(req *openai.ChatCompletionRequest, opts *ChatOptions, isStream bool) (any, bool) {
	// DeepSeek V3.x 계열에만 thinking 파라미터 특별 처리
	// R1 계열 모델은 chain-of-thought 기본 활성, 추가 파라미터 불필요
	if !c.isDeepSeekV3Model() || opts == nil || opts.Thinking == nil {
		return nil, false // 표준 요청 사용
	}

	// LKEAP 전용 요청 구성
	lkeapReq := LKEAPChatCompletionRequest{
		ChatCompletionRequest: *req,
	}

	thinkingType := "disabled"
	if *opts.Thinking {
		thinkingType = "enabled"
	}
	lkeapReq.Thinking = &LKEAPThinkingConfig{Type: thinkingType}

	return lkeapReq, true // 원본 HTTP 요청 사용
}
