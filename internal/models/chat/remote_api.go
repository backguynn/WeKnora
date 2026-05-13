package chat

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Tencent/WeKnora/internal/logger"
	"github.com/Tencent/WeKnora/internal/models/provider"
	"github.com/Tencent/WeKnora/internal/types"
	secutils "github.com/Tencent/WeKnora/internal/utils"
	"github.com/sashabaranov/go-openai"
)

// LLM 호출 타임아웃 설정입니다. 상위 계층에서 deadline을 설정하지 않았을 때만
// 요청이 영구 대기 상태로 worker를 막지 않도록 보조 타임아웃으로 사용합니다.
// 상위 ctx에 deadline이 이미 있으면(기본값보다 짧거나 길더라도) 그대로 존중하며,
// 추가 기본 타임아웃은 덧붙이지 않습니다. 환경 변수로 재정의할 수 있습니다.
//   - WEKNORA_LLM_CHAT_TIMEOUT_SECONDS    비스트리밍 호출 보조 타임아웃(기본 600s)
//   - WEKNORA_LLM_STREAM_TIMEOUT_SECONDS  스트리밍 호출 보조 타임아웃(기본 1800s)
var (
	defaultChatTimeout   = envDurationSeconds("WEKNORA_LLM_CHAT_TIMEOUT_SECONDS", 300*time.Second)
	defaultStreamTimeout = envDurationSeconds("WEKNORA_LLM_STREAM_TIMEOUT_SECONDS", 600*time.Second)
)

// envDurationSeconds는 "초" 단위 환경 변수를 읽고, 파싱 실패 또는 0 이하일 때 fallback으로 되돌립니다.
func envDurationSeconds(key string, fallback time.Duration) time.Duration {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil || n <= 0 {
		return fallback
	}
	return time.Duration(n) * time.Second
}

// withLLMTimeout은 상위 ctx에 deadline이 없을 때만 보조 타임아웃을 추가합니다.
// 상위 호출자가 이미 deadline을 명시했다면(짧든 길든) 그대로 반환해,
// 최종 타임아웃 정책 결정권을 호출자에게 둡니다.
func withLLMTimeout(ctx context.Context, d time.Duration) (context.Context, context.CancelFunc) {
	if _, ok := ctx.Deadline(); ok {
		return ctx, func() {}
	}
	return context.WithTimeout(ctx, d)
}

// rawHTTPClient is a shared HTTP client for raw HTTP LLM calls with connection-level timeouts.
// Per-request timeout is enforced via context deadline (see defaultChatTimeout / defaultStreamTimeout)
// rather than http.Client.Timeout, so streaming calls are not prematurely terminated.
// Uses SSRFSafeDialContext to prevent DNS rebinding attacks at the connection layer.
var rawHTTPClient = &http.Client{
	Transport: &http.Transport{
		Proxy:               http.ProxyFromEnvironment,
		DialContext:         secutils.SSRFSafeDialContext,
		TLSHandshakeTimeout: 10 * time.Second,
		IdleConnTimeout:     90 * time.Second,
		MaxIdleConnsPerHost: 5,
	},
}

// RemoteAPIChat은 OpenAI 호환 API 기반 채팅을 구현합니다.
// provider별 특수 로직을 포함하지 않는 범용 구현입니다.
type RemoteAPIChat struct {
	modelName string
	client    *openai.Client
	modelID   string
	baseURL   string
	apiKey    string
	provider  provider.ProviderName
	appID     string
	appSecret string
	// customHeaders는 사용자가 모델 설정에서 지정한 사용자 정의 HTTP 헤더입니다.
	// OpenAI Python SDK의 extra_headers와 유사한 개념입니다.
	customHeaders map[string]string

	// requestCustomizer는 하위 구현이 요청을 사용자 정의할 수 있게 합니다.
	// 반환값은 사용자 정의 요청 본문(nil이면 표준 요청 사용)과 원시 HTTP 사용 여부입니다.
	requestCustomizer func(req *openai.ChatCompletionRequest, opts *ChatOptions, isStream bool) (customReq any, useRawHTTP bool)

	// endpointCustomizer는 하위 구현이 요청 endpoint를 사용자 정의할 수 있게 합니다.
	// 빈 문자열을 반환하면 기본 OpenAI 형식 endpoint를 사용합니다.
	endpointCustomizer func(baseURL string, modelID string, isStream bool) (endpoint string)

	// headerCustomizer는 하위 구현이 원시 HTTP 요청 헤더(예: 서명 인증)를 사용자 정의할 수 있게 합니다.
	headerCustomizer func(req *http.Request, body []byte) error
}

// NewRemoteAPIChat은 원격 API 채팅 인스턴스를 생성합니다.
func NewRemoteAPIChat(chatConfig *ChatConfig) (*RemoteAPIChat, error) {
	if chatConfig.BaseURL != "" {
		if err := secutils.ValidateURLForSSRF(chatConfig.BaseURL); err != nil {
			return nil, fmt.Errorf("baseURL SSRF check failed: %w", err)
		}
	}

	apiKey := chatConfig.APIKey
	providerName := provider.ProviderName(chatConfig.Provider)
	if providerName == "" {
		providerName = provider.DetectProvider(chatConfig.BaseURL)
	}

	var config openai.ClientConfig
	if providerName == provider.ProviderAzureOpenAI {
		config = openai.DefaultAzureConfig(apiKey, chatConfig.BaseURL)
		config.AzureModelMapperFunc = func(model string) string {
			return model
		}
		if chatConfig.ExtraConfig != nil {
			if v, ok := chatConfig.ExtraConfig["api_version"]; ok {
				config.APIVersion = v
			}
		}
	} else {
		config = openai.DefaultConfig(apiKey)
		if baseURL := chatConfig.BaseURL; baseURL != "" {
			config.BaseURL = baseURL
		}
	}

	// CustomHeaders가 지정되면 SDK가 사용하는 HTTPClient에 RoundTripper를 감싸
	// 모든 요청에 자동으로 헤더를 주입합니다(raw HTTP 경로는 전송 전에 별도 처리).
	if len(chatConfig.CustomHeaders) > 0 {
		if httpClient, ok := config.HTTPClient.(*http.Client); ok {
			config.HTTPClient = secutils.WrapHTTPClientWithHeaders(httpClient, chatConfig.CustomHeaders)
		} else {
			// SDK 기본값으로 HTTPClient가 nil이면, 헤더가 주입된 새 client를 구성합니다.
			config.HTTPClient = secutils.WrapHTTPClientWithHeaders(nil, chatConfig.CustomHeaders)
		}
	}

	modelName := chatConfig.ModelName
	if chatConfig.ExtraConfig != nil {
		if override := strings.TrimSpace(chatConfig.ExtraConfig["remote_model_name"]); override != "" {
			modelName = override
		}
	}
	if providerName == provider.ProviderWeKnoraCloud {
		if chatConfig.AppID == "" {
			return nil, fmt.Errorf("WeKnoraCloud provider: AppID is required")
		}
		if chatConfig.AppSecret == "" {
			return nil, fmt.Errorf("WeKnoraCloud provider: AppSecret is required")
		}
	}

	return &RemoteAPIChat{
		modelName:     modelName,
		client:        openai.NewClientWithConfig(config),
		modelID:       chatConfig.ModelID,
		baseURL:       chatConfig.BaseURL,
		apiKey:        apiKey,
		provider:      providerName,
		appID:         chatConfig.AppID,
		appSecret:     chatConfig.AppSecret,
		customHeaders: chatConfig.CustomHeaders,
	}, nil
}

// SetRequestCustomizer는 요청 사용자 정의 함수를 설정합니다.
func (c *RemoteAPIChat) SetRequestCustomizer(customizer func(req *openai.ChatCompletionRequest, opts *ChatOptions, isStream bool) (any, bool)) {
	c.requestCustomizer = customizer
}

// SetEndpointCustomizer는 endpoint 사용자 정의 함수를 설정합니다.
func (c *RemoteAPIChat) SetEndpointCustomizer(customizer func(baseURL string, modelID string, isStream bool) string) {
	c.endpointCustomizer = customizer
}

// SetHeaderCustomizer는 원시 HTTP 헤더 사용자 정의 함수를 설정합니다.
func (c *RemoteAPIChat) SetHeaderCustomizer(customizer func(req *http.Request, body []byte) error) {
	c.headerCustomizer = customizer
}

// ConvertMessages는 메시지를 OpenAI 형식으로 변환합니다.
// 하위 구현에서도 사용할 수 있도록 공개되어 있습니다.
func (c *RemoteAPIChat) ConvertMessages(messages []Message) []openai.ChatCompletionMessage {
	openaiMessages := make([]openai.ChatCompletionMessage, 0, len(messages))
	for _, msg := range messages {
		openaiMsg := openai.ChatCompletionMessage{
			Role: msg.Role,
		}

		// 이미지 등을 포함하는 다중 콘텐츠 메시지를 우선 처리합니다.
		if len(msg.MultiContent) > 0 {
			openaiMsg.MultiContent = make([]openai.ChatMessagePart, 0, len(msg.MultiContent))
			for _, part := range msg.MultiContent {
				switch part.Type {
				case "text":
					openaiMsg.MultiContent = append(openaiMsg.MultiContent, openai.ChatMessagePart{
						Type: openai.ChatMessagePartTypeText,
						Text: part.Text,
					})
				case "image_url":
					if part.ImageURL != nil {
						openaiMsg.MultiContent = append(openaiMsg.MultiContent, openai.ChatMessagePart{
							Type: openai.ChatMessagePartTypeImageURL,
							ImageURL: &openai.ChatMessageImageURL{
								URL:    part.ImageURL.URL,
								Detail: openai.ImageURLDetail(part.ImageURL.Detail),
							},
						})
					}
				}
			}
		} else if len(msg.Images) > 0 && msg.Role == "user" {
			parts := make([]openai.ChatMessagePart, 0, len(msg.Images)+1)
			for _, imgURL := range msg.Images {
				resolved := resolveImageURLForLLM(imgURL)
				parts = append(parts, openai.ChatMessagePart{
					Type: openai.ChatMessagePartTypeImageURL,
					ImageURL: &openai.ChatMessageImageURL{
						URL:    resolved,
						Detail: openai.ImageURLDetailAuto,
					},
				})
			}
			parts = append(parts, openai.ChatMessagePart{
				Type: openai.ChatMessagePartTypeText,
				Text: msg.Content,
			})
			openaiMsg.MultiContent = parts
		} else if msg.Content != "" {
			openaiMsg.Content = msg.Content
		}

		if len(msg.ToolCalls) > 0 {
			openaiMsg.ToolCalls = make([]openai.ToolCall, 0, len(msg.ToolCalls))
			for _, tc := range msg.ToolCalls {
				toolType := openai.ToolType(tc.Type)
				openaiMsg.ToolCalls = append(openaiMsg.ToolCalls, openai.ToolCall{
					ID:   tc.ID,
					Type: toolType,
					Function: openai.FunctionCall{
						Name:      tc.Function.Name,
						Arguments: tc.Function.Arguments,
					},
				})
			}
		}

		if msg.Role == "tool" {
			openaiMsg.ToolCallID = msg.ToolCallID
			openaiMsg.Name = msg.Name
		}

		openaiMessages = append(openaiMessages, openaiMsg)
	}
	return openaiMessages
}

// BuildChatCompletionRequest는 표준 채팅 요청 파라미터를 구성합니다.
// 하위 구현에서도 사용할 수 있도록 공개되어 있습니다.
func (c *RemoteAPIChat) BuildChatCompletionRequest(messages []Message, opts *ChatOptions, isStream bool) openai.ChatCompletionRequest {
	req := openai.ChatCompletionRequest{
		Model:    c.modelName,
		Messages: c.ConvertMessages(messages),
		Stream:   isStream,
	}

	if isStream {
		req.StreamOptions = &openai.StreamOptions{IncludeUsage: true}
	}

	if opts != nil {
		req.Temperature = float32(opts.Temperature)
		if opts.TopP > 0 {
			req.TopP = float32(opts.TopP)
		}
		// Token limit handling: prefer max_completion_tokens for OpenAI API v1.0+ compatibility
		// If only max_tokens is set, use it as max_completion_tokens value
		// Don't set MaxTokens to avoid "Unsupported parameter" error with latest OpenAI models
		if opts.MaxCompletionTokens > 0 {
			req.MaxCompletionTokens = opts.MaxCompletionTokens
		} else if opts.MaxTokens > 0 {
			req.MaxCompletionTokens = opts.MaxTokens
		}
		if opts.FrequencyPenalty > 0 {
			req.FrequencyPenalty = float32(opts.FrequencyPenalty)
		}
		if opts.PresencePenalty > 0 {
			req.PresencePenalty = float32(opts.PresencePenalty)
		}

		// Tools 처리
		if len(opts.Tools) > 0 {
			req.Tools = make([]openai.Tool, 0, len(opts.Tools))
			for _, tool := range opts.Tools {
				toolType := openai.ToolType(tool.Type)
				openaiTool := openai.Tool{
					Type: toolType,
					Function: &openai.FunctionDefinition{
						Name:        tool.Function.Name,
						Description: tool.Function.Description,
					},
				}
				if tool.Function.Parameters != nil {
					openaiTool.Function.Parameters = tool.Function.Parameters
				}
				req.Tools = append(req.Tools, openaiTool)
			}
		}

		// ParallelToolCalls 처리
		if opts.ParallelToolCalls != nil {
			val := *opts.ParallelToolCalls
			req.ParallelToolCalls = val
		}

		// ToolChoice 처리(표준 구현)
		if opts.ToolChoice != "" {
			switch opts.ToolChoice {
			case "none", "required", "auto":
				req.ToolChoice = opts.ToolChoice
			default:
				req.ToolChoice = openai.ToolChoice{
					Type: "function",
					Function: openai.ToolFunction{
						Name: opts.ToolChoice,
					},
				}
			}
		}

		if len(opts.Format) > 0 {
			req.ResponseFormat = &openai.ChatCompletionResponseFormat{
				Type: openai.ChatCompletionResponseFormatTypeJSONObject,
			}
			req.Messages[len(req.Messages)-1].Content += fmt.Sprintf("\nUse this JSON schema: %s", opts.Format)
		}
	}

	return req
}

// logRequest는 요청 로그를 기록합니다.
func (c *RemoteAPIChat) logRequest(ctx context.Context, req any, isStream bool) {
	if jsonData, err := json.MarshalIndent(req, "", "  "); err == nil {
		logger.Infof(ctx, "[LLM Request] model=%s, stream=%v, request:\n%s", c.modelName, isStream, secutils.CompactImageDataURLForLog(string(jsonData)))
	}
}

// Chat은 비스트리밍 채팅을 수행합니다.
func (c *RemoteAPIChat) Chat(ctx context.Context, messages []Message, opts *ChatOptions) (*types.ChatResponse, error) {
	// 호출자가 deadline을 설정하지 않았을 때만 보조 타임아웃을 추가해
	// hung 요청이 worker를 영구 점유하지 않도록 합니다.
	// 호출자가 더 짧거나 더 긴 deadline을 명시했다면 그대로 존중합니다.
	timeoutCtx, cancel := withLLMTimeout(ctx, defaultChatTimeout)
	defer cancel()

	req := c.BuildChatCompletionRequest(messages, opts, false)
	var customEndpoint string
	if c.endpointCustomizer != nil {
		customEndpoint = c.endpointCustomizer(c.baseURL, c.modelID, true)
	}
	// 사용자 정의 요청이 필요한지 확인
	if c.requestCustomizer != nil {
		customReq, useRawHTTP := c.requestCustomizer(&req, opts, false)
		if useRawHTTP && customReq != nil {
			return c.chatWithRawHTTP(timeoutCtx, customEndpoint, customReq)
		}
	}

	// 사용자 정의 endpoint 사용
	if customEndpoint != "" {
		return c.chatWithRawHTTP(timeoutCtx, customEndpoint, &req)
	}

	c.logRequest(timeoutCtx, req, false)
	resp, err := c.client.CreateChatCompletion(timeoutCtx, req)
	if err != nil {
		if isMultimodalNotSupportedError(err) {
			logger.Warnf(timeoutCtx, "[LLM Request] Model %s does not support multimodal, retrying without images", c.modelName)
			cleaned := stripImagesFromMessages(messages)
			req = c.BuildChatCompletionRequest(cleaned, opts, false)
			resp, err = c.client.CreateChatCompletion(timeoutCtx, req)
		}
		if err != nil {
			return nil, fmt.Errorf("create chat completion: %w", err)
		}
	}

	result, err := c.parseCompletionResponse(&resp)
	if err != nil {
		return nil, err
	}
	logger.Infof(timeoutCtx, "[LLM Usage] model=%s, prompt_tokens=%d, completion_tokens=%d, total_tokens=%d",
		c.modelName, result.Usage.PromptTokens, result.Usage.CompletionTokens, result.Usage.TotalTokens)
	return result, nil
}

// chatWithRawHTTP는 원시 HTTP 요청으로 채팅을 수행합니다(사용자 정의 요청용).
func (c *RemoteAPIChat) chatWithRawHTTP(ctx context.Context, endpoint string, customReq any) (*types.ChatResponse, error) {
	jsonData, err := json.Marshal(customReq)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	if endpoint == "" {
		endpoint = c.baseURL + "/chat/completions"
	}
	if err := secutils.ValidateURLForSSRF(endpoint); err != nil {
		return nil, fmt.Errorf("endpoint SSRF check failed: %w", err)
	}
	logger.Infof(ctx, "[LLM Request] Remote HTTP, endpoint=%s, model=%s, raw HTTP request:\n%s",
		endpoint, c.modelName, secutils.CompactImageDataURLForLog(string(jsonData)))

	httpReq, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	if c.headerCustomizer != nil {
		if err := c.headerCustomizer(httpReq, jsonData); err != nil {
			return nil, fmt.Errorf("customize headers: %w", err)
		}
	} else if c.provider == provider.ProviderAzureOpenAI {
		httpReq.Header.Set("api-key", c.apiKey)
	} else {
		httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	// 사용자 정의 header 주입(예약 헤더는 내부에서 자동 건너뜀)
	secutils.ApplyCustomHeaders(httpReq, c.customHeaders)

	logger.Infof(ctx, "[LLM Request] Remote HTTP, endpoint=%s, model=%s",
		endpoint, c.modelName)

	resp, err := rawHTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var chatResp openai.ChatCompletionResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	result, err := c.parseCompletionResponse(&chatResp)
	if err != nil {
		return nil, err
	}
	logger.Infof(ctx, "[LLM Usage] model=%s, prompt_tokens=%d, completion_tokens=%d, total_tokens=%d",
		c.modelName, result.Usage.PromptTokens, result.Usage.CompletionTokens, result.Usage.TotalTokens)
	return result, nil
}

// parseCompletionResponse는 비스트리밍 응답을 파싱합니다.
func (c *RemoteAPIChat) parseCompletionResponse(resp *openai.ChatCompletionResponse) (*types.ChatResponse, error) {
	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response from API")
	}

	choice := resp.Choices[0]

	// 사고형 모델 출력에서 <think></think>로 감싼 사고 과정을 제거합니다.
	// Thinking=false여도 사고 내용을 반환하는 경우와, Thinking=false를 지원하지 않는 일부 사고형 모델(예: Miniax-M2.1)을 위한 보조 처리입니다.
	content := removeThinkingContent(choice.Message.Content)

	response := &types.ChatResponse{
		Content:      content,
		FinishReason: string(choice.FinishReason),
		Usage: types.TokenUsage{
			PromptTokens:     resp.Usage.PromptTokens,
			CompletionTokens: resp.Usage.CompletionTokens,
			TotalTokens:      resp.Usage.TotalTokens,
		},
	}

	if len(choice.Message.ToolCalls) > 0 {
		response.ToolCalls = make([]types.LLMToolCall, 0, len(choice.Message.ToolCalls))
		for _, tc := range choice.Message.ToolCalls {
			response.ToolCalls = append(response.ToolCalls, types.LLMToolCall{
				ID:   tc.ID,
				Type: string(tc.Type),
				Function: types.FunctionCall{
					Name:      tc.Function.Name,
					Arguments: tc.Function.Arguments,
				},
			})
		}
	}

	return response, nil
}

// removeThinkingContent는 사고형 모델 출력의 <think></think> 구간을 제거합니다.
// 내용이 <think>로 시작할 때만 처리합니다.
func removeThinkingContent(content string) string {
	const thinkStartTag = "<think>"
	const thinkEndTag = "</think>"

	trimmed := strings.TrimSpace(content)
	if !strings.HasPrefix(trimmed, thinkStartTag) {
		return content
	}

	// 중첩된 경우까지 고려해 마지막 </think> 태그를 찾습니다.
	if lastEndIdx := strings.LastIndex(trimmed, thinkEndTag); lastEndIdx != -1 {
		if result := strings.TrimSpace(trimmed[lastEndIdx+len(thinkEndTag):]); result != "" {
			return result
		}
		return ""
	}

	return "" // </think>를 찾지 못한 경우이며, 사고 내용이 너무 길어 잘렸을 수 있으므로 빈 문자열을 반환합니다.
}

// ChatStream은 스트리밍 채팅을 수행합니다.
func (c *RemoteAPIChat) ChatStream(ctx context.Context, messages []Message, opts *ChatOptions) (<-chan types.StreamResponse, error) {
	// 호출자가 deadline을 지정하지 않았을 때만 보조 타임아웃을 추가합니다.
	// 사고/추론형 모델은 첫 token이 나오기까지 수십 초에서 수분이 걸릴 수 있어 스트리밍 기본 타임아웃을 더 길게 둡니다.
	timeoutCtx, cancel := withLLMTimeout(ctx, defaultStreamTimeout)

	req := c.BuildChatCompletionRequest(messages, opts, true)

	var customEndpoint string
	if c.endpointCustomizer != nil {
		customEndpoint = c.endpointCustomizer(c.baseURL, c.modelID, true)
	}

	// 사용자 정의 요청이 필요한지 확인
	if c.requestCustomizer != nil {
		customReq, useRawHTTP := c.requestCustomizer(&req, opts, true)
		if useRawHTTP && customReq != nil {
			ch, err := c.chatStreamWithRawHTTP(timeoutCtx, customEndpoint, customReq)
			return wrapStreamCancel(ch, err, cancel)
		}
	}
	// 사용자 정의 endpoint 사용
	if customEndpoint != "" {
		ch, err := c.chatStreamWithRawHTTP(timeoutCtx, customEndpoint, &req)
		return wrapStreamCancel(ch, err, cancel)
	}
	c.logRequest(timeoutCtx, req, true)

	streamChan := make(chan types.StreamResponse)

	stream, err := c.client.CreateChatCompletionStream(timeoutCtx, req)
	if err != nil {
		if isMultimodalNotSupportedError(err) {
			logger.Warnf(timeoutCtx, "[LLM Stream] Model %s does not support multimodal, retrying without images", c.modelName)
			cleaned := stripImagesFromMessages(messages)
			req = c.BuildChatCompletionRequest(cleaned, opts, true)
			stream, err = c.client.CreateChatCompletionStream(timeoutCtx, req)
		}
		if err != nil {
			cancel()
			close(streamChan)
			return nil, fmt.Errorf("create chat completion stream: %w", err)
		}
	}

	go func() {
		defer cancel()
		c.processStream(timeoutCtx, stream, streamChan)
	}()

	return streamChan, nil
}

// wrapStreamCancel은 하위 channel이 닫힌 뒤 cancel을 실행해 timeout context 누수를 막습니다.
// 하위 호출이 바로 error를 반환하면 즉시 cancel하고 error를 그대로 전달합니다.
func wrapStreamCancel(in <-chan types.StreamResponse, err error, cancel context.CancelFunc) (<-chan types.StreamResponse, error) {
	if err != nil {
		cancel()
		return nil, err
	}
	out := make(chan types.StreamResponse)
	go func() {
		defer cancel()
		defer close(out)
		for v := range in {
			out <- v
		}
	}()
	return out, nil
}

// chatStreamWithRawHTTP는 원시 HTTP 요청으로 스트리밍 채팅을 수행합니다.
func (c *RemoteAPIChat) chatStreamWithRawHTTP(ctx context.Context, endpoint string, customReq any) (<-chan types.StreamResponse, error) {
	jsonData, err := json.Marshal(customReq)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	if endpoint == "" {
		endpoint = c.baseURL + "/chat/completions"
	}
	if err := secutils.ValidateURLForSSRF(endpoint); err != nil {
		return nil, fmt.Errorf("endpoint SSRF check failed: %w", err)
	}

	if prettyJSON, pErr := json.MarshalIndent(customReq, "", "  "); pErr == nil {
		logger.Infof(ctx, "[LLM Stream Request] endpoint=%s, model=%s, stream=true, request:\n%s",
			endpoint, c.modelName, secutils.CompactImageDataURLForLog(string(prettyJSON)))
	} else {
		logger.Infof(ctx, "[LLM Stream] endpoint=%s, model=%s", endpoint, c.modelName)
	}
	httpReq, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	if c.headerCustomizer != nil {
		if err := c.headerCustomizer(httpReq, jsonData); err != nil {
			return nil, fmt.Errorf("customize headers: %w", err)
		}
	} else if c.provider == provider.ProviderAzureOpenAI {
		httpReq.Header.Set("api-key", c.apiKey)
	} else {
		httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	}
	httpReq.Header.Set("Accept", "text/event-stream")

	// 사용자 정의 header를 주입합니다. 예약 헤더는 내부에서 자동으로 건너뜁니다.
	secutils.ApplyCustomHeaders(httpReq, c.customHeaders)

	resp, err := rawHTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	streamChan := make(chan types.StreamResponse)

	go c.processRawHTTPStream(ctx, resp, streamChan)

	return streamChan, nil
}

// processStream은 OpenAI SDK 스트리밍 응답을 처리합니다.
func (c *RemoteAPIChat) processStream(ctx context.Context, stream *openai.ChatCompletionStream, streamChan chan types.StreamResponse) {
	defer close(streamChan)
	defer stream.Close()

	state := newStreamState()

	for {
		response, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				if state.usage != nil {
					logger.Infof(ctx, "[LLM Usage] model=%s, prompt_tokens=%d, completion_tokens=%d, total_tokens=%d",
						c.modelName, state.usage.PromptTokens, state.usage.CompletionTokens, state.usage.TotalTokens)
				}
				toolCalls := state.buildOrderedToolCalls()
				streamChan <- types.StreamResponse{
					ResponseType: types.ResponseTypeAnswer,
					Content:      "",
					Done:         true,
					ToolCalls:    toolCalls,
					Usage:        state.usage,
					FinishReason: state.lastFinishReason,
				}
			} else {
				streamChan <- types.StreamResponse{
					ResponseType: types.ResponseTypeError,
					Content:      err.Error(),
					Done:         true,
				}
			}
			return
		}

		if response.Usage != nil {
			state.usage = &types.TokenUsage{
				PromptTokens:     response.Usage.PromptTokens,
				CompletionTokens: response.Usage.CompletionTokens,
				TotalTokens:      response.Usage.TotalTokens,
			}
		}

		if len(response.Choices) > 0 {
			c.processStreamDelta(ctx, &response.Choices[0], state, streamChan, response.Choices[0].Delta.ReasoningContent)
		}
	}
}

// processRawHTTPStream은 원시 HTTP 스트리밍 응답을 처리합니다.
func (c *RemoteAPIChat) processRawHTTPStream(ctx context.Context, resp *http.Response, streamChan chan types.StreamResponse) {
	defer close(streamChan)
	defer resp.Body.Close()

	state := newStreamState()
	reader := NewSSEReader(resp.Body)

	for {
		event, err := reader.ReadEvent()
		if err != nil {
			if err == io.EOF {
				if state.usage != nil {
					logger.Infof(ctx, "[LLM Usage] model=%s, prompt_tokens=%d, completion_tokens=%d, total_tokens=%d",
						c.modelName, state.usage.PromptTokens, state.usage.CompletionTokens, state.usage.TotalTokens)
				}
				toolCalls := state.buildOrderedToolCalls()
				streamChan <- types.StreamResponse{
					ResponseType: types.ResponseTypeAnswer,
					Content:      "",
					Done:         true,
					ToolCalls:    toolCalls,
					Usage:        state.usage,
				}
			} else {
				logger.Errorf(ctx, "Stream read error: %v", err)
				streamChan <- types.StreamResponse{
					ResponseType: types.ResponseTypeError,
					Content:      err.Error(),
					Done:         true,
				}
			}
			return
		}

		if event == nil {
			continue
		}

		if event.Done {
			if state.usage != nil {
				logger.Infof(ctx, "[LLM Usage] model=%s, prompt_tokens=%d, completion_tokens=%d, total_tokens=%d",
					c.modelName, state.usage.PromptTokens, state.usage.CompletionTokens, state.usage.TotalTokens)
			}
			toolCalls := state.buildOrderedToolCalls()
			streamChan <- types.StreamResponse{
				ResponseType: types.ResponseTypeAnswer,
				Content:      "",
				Done:         true,
				ToolCalls:    toolCalls,
				Usage:        state.usage,
			}
			return
		}

		if event.Data == nil {
			continue
		}

		// 지역 구조체로 한 번에 파싱하면서 표준 필드와 vLLM reasoning 필드를 함께 수집해 성능 손실을 줄입니다.
		var streamResp struct {
			openai.ChatCompletionStreamResponse
			Choices []struct {
				Index int `json:"index"`
				Delta struct {
					openai.ChatCompletionStreamChoiceDelta
					Reasoning string `json:"reasoning,omitempty"`
				} `json:"delta"`
				FinishReason openai.FinishReason `json:"finish_reason"`
			} `json:"choices"`
		}

		if err := json.Unmarshal(event.Data, &streamResp); err != nil {
			logger.Errorf(ctx, "Failed to parse stream response: %v", err)
			continue
		}

		if streamResp.Usage != nil {
			state.usage = &types.TokenUsage{
				PromptTokens:     streamResp.Usage.PromptTokens,
				CompletionTokens: streamResp.Usage.CompletionTokens,
				TotalTokens:      streamResp.Usage.TotalTokens,
			}
		}

		if len(streamResp.Choices) > 0 {
			choice := streamResp.Choices[0]
			// 표준 경로와 vLLM 경로를 모두 지원하는 공통 추출 로직
			reasoning := choice.Delta.Reasoning
			if reasoning == "" {
				reasoning = choice.Delta.ReasoningContent
			}

			// 표준 SDK 호환 choice 객체를 만들어 하위 로직에 넘겨 기존 처리 흐름을 유지합니다.
			sdkChoice := openai.ChatCompletionStreamChoice{
				Index:        choice.Index,
				Delta:        choice.Delta.ChatCompletionStreamChoiceDelta,
				FinishReason: choice.FinishReason,
			}
			c.processStreamDelta(ctx, &sdkChoice, state, streamChan, reasoning)
		}
	}
}

// streamState는 스트리밍 처리 상태를 보관합니다.
type streamState struct {
	toolCallMap      map[int]*types.LLMToolCall
	lastFunctionName map[int]string
	nameNotified     map[int]bool
	hasThinking      bool
	fieldExtractors  map[int]*jsonFieldExtractor // per tool-call-index extractors for streaming field extraction
	usage            *types.TokenUsage           // captured from the final stream chunk when include_usage is enabled
	lastFinishReason string                      // last observed finish_reason for EOF handler fallback

	// Diagnostic flags (fire-once) used to log earliest signals of tool_call
	// presence/absence at the OpenAI-protocol level. These are independent of
	// the higher-level ResponseTypeToolCall marker (which only fires once
	// function name has stabilized) and let us distinguish between
	//   (A) no tool_calls field ever observed (true natural-stop), and
	//   (B) tool_calls field observed but marker not yet emitted.
	firstToolCallSeen    bool // true once any delta carried tool_calls
	noToolCallStopLogged bool // true once we logged "stop without tool_calls"
	firstContentSeen     bool // true once delta.Content first appeared
	firstReasoningSeen   bool // true once reasoning_content first appeared
	streamStartedAt      time.Time
}

func newStreamState() *streamState {
	return &streamState{
		toolCallMap:      make(map[int]*types.LLMToolCall),
		lastFunctionName: make(map[int]string),
		nameNotified:     make(map[int]bool),
		hasThinking:      false,
		fieldExtractors:  make(map[int]*jsonFieldExtractor),
		streamStartedAt:  time.Now(),
	}
}

// elapsedMs returns the milliseconds elapsed since the stream state was
// initialized. Used to attach time-since-stream-start to fire-once diagnostic
// logs so a single grep can reveal the temporal layout of a single stream
// (TTFC / TTFT / first-tool-call / natural-stop confirmation, etc).
func (s *streamState) elapsedMs() int64 {
	if s.streamStartedAt.IsZero() {
		return 0
	}
	return time.Since(s.streamStartedAt).Milliseconds()
}

func (s *streamState) buildOrderedToolCalls() []types.LLMToolCall {
	if len(s.toolCallMap) == 0 {
		return nil
	}
	result := make([]types.LLMToolCall, 0, len(s.toolCallMap))
	for i := 0; i < len(s.toolCallMap); i++ {
		if tc, ok := s.toolCallMap[i]; ok && tc != nil {
			result = append(result, *tc)
		}
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

// processStreamDelta는 스트리밍 응답의 개별 delta를 처리합니다.
func (c *RemoteAPIChat) processStreamDelta(ctx context.Context, choice *openai.ChatCompletionStreamChoice, state *streamState, streamChan chan types.StreamResponse, reasoningContent string) {
	delta := choice.Delta
	isDone := string(choice.FinishReason) != ""

	// Track finish_reason for EOF handler fallback
	if isDone {
		state.lastFinishReason = string(choice.FinishReason)
	}

	// tool calls 처리
	if len(delta.ToolCalls) > 0 {
		c.processToolCallsDelta(ctx, delta.ToolCalls, state, streamChan)
	}

	// Earliest reliable "no tool_calls" signal at the OpenAI-protocol level:
	// finish_reason=stop arrived AND we never observed a tool_calls field on
	// any prior delta. Logged once per stream so callers can grep for the
	// natural-stop entry point without waiting for the higher-level summary.
	if isDone &&
		string(choice.FinishReason) == "stop" &&
		!state.firstToolCallSeen &&
		!state.noToolCallStopLogged {
		logger.Infof(ctx, "[LLM Stream] Natural-stop at OpenAI layer "+
			"(finish=stop, tool_calls field never observed, thinking_seen=%t, "+
			"first_content_seen=%t, elapsed_ms=%d)",
			state.hasThinking, state.firstContentSeen, state.elapsedMs())
		state.noToolCallStopLogged = true
	}

	// 사고 내용 전송(ReasoningContent, DeepSeek 등 지원)
	if reasoningContent != "" {
		// Earliest reasoning_content signal at the OpenAI-protocol level. Fired
		// once per stream so we can distinguish "model emitted thinking before
		// answer" vs "model never produced thinking" when triaging logs.
		if !state.firstReasoningSeen {
			state.firstReasoningSeen = true
			logger.Infof(ctx, "[LLM Stream] First reasoning_content at OpenAI layer "+
				"(len=%d, preview=%q, elapsed_ms=%d)",
				len(reasoningContent), truncateForDebug(reasoningContent, 80), state.elapsedMs())
		}
		state.hasThinking = true
		streamChan <- types.StreamResponse{
			ResponseType: types.ResponseTypeThinking,
			Content:      reasoningContent,
			Done:         false,
		}
	}

	// 답변 내용 전송
	if delta.Content != "" {
		// Earliest delta.Content signal at the OpenAI-protocol level. Fired once
		// per stream so we can measure TTFC (time-to-first-content) and tell
		// "answer started before any tool_call" from "tool_call came first".
		if !state.firstContentSeen {
			state.firstContentSeen = true
			logger.Infof(ctx, "[LLM Stream] First delta.Content at OpenAI layer "+
				"(len=%d, preview=%q, tool_call_seen=%t, thinking_seen=%t, elapsed_ms=%d)",
				len(delta.Content), truncateForDebug(delta.Content, 80),
				state.firstToolCallSeen, state.firstReasoningSeen, state.elapsedMs())
		}
		// If we had thinking content and this is the first answer chunk,
		// send a thinking done event first
		if state.hasThinking {
			streamChan <- types.StreamResponse{
				ResponseType: types.ResponseTypeThinking,
				Content:      "",
				Done:         true,
			}
			state.hasThinking = false // Only send once
		}
		streamChan <- types.StreamResponse{
			ResponseType: types.ResponseTypeAnswer,
			Content:      delta.Content,
			Done:         isDone,
			ToolCalls:    state.buildOrderedToolCalls(),
			FinishReason: string(choice.FinishReason),
		}
	}

	if isDone && len(state.toolCallMap) > 0 {
		streamChan <- types.StreamResponse{
			ResponseType: types.ResponseTypeAnswer,
			Content:      "",
			Done:         true,
			ToolCalls:    state.buildOrderedToolCalls(),
			FinishReason: string(choice.FinishReason),
		}
	}

	// Ensure thinking done is sent when stream finishes without any answer content
	// (e.g., model only produced reasoning then hit finish_reason with empty content).
	if isDone && state.hasThinking {
		streamChan <- types.StreamResponse{
			ResponseType: types.ResponseTypeThinking,
			Content:      "",
			Done:         true,
		}
		state.hasThinking = false
	}

	// Catch-all: isDone but none of the above branches sent a response with
	// FinishReason (empty content, no tool calls, no thinking). This prevents
	// the finish_reason from being lost in the streaming pipeline.
	if isDone && delta.Content == "" && len(state.toolCallMap) == 0 && !state.hasThinking {
		streamChan <- types.StreamResponse{
			ResponseType: types.ResponseTypeAnswer,
			Done:         true,
			FinishReason: string(choice.FinishReason),
		}
	}
}

// processToolCallsDelta는 tool calls의 증분 업데이트를 처리합니다.
func (c *RemoteAPIChat) processToolCallsDelta(ctx context.Context, toolCalls []openai.ToolCall, state *streamState, streamChan chan types.StreamResponse) {
	// Earliest signal at the OpenAI-protocol level that this stream will
	// produce at least one tool call. Fires *before* the function name has
	// stabilized, i.e. earlier than the higher-level ResponseTypeToolCall
	// marker downstream consumers see. Useful for distinguishing
	// "tool_calls field arrived but marker not yet emitted" from
	// "tool_calls field truly absent" when triaging stream behavior.
	if !state.firstToolCallSeen && len(toolCalls) > 0 {
		state.firstToolCallSeen = true
		var firstID, firstName string
		for _, tc := range toolCalls {
			if tc.ID != "" {
				firstID = tc.ID
			}
			if tc.Function.Name != "" {
				firstName = tc.Function.Name
			}
			if firstID != "" || firstName != "" {
				break
			}
		}
		logger.Infof(ctx, "[LLM Stream] First tool_calls delta at OpenAI layer "+
			"(count=%d, first_id=%q, first_name=%q, "+
			"first_content_seen=%t, thinking_seen=%t, elapsed_ms=%d)",
			len(toolCalls), firstID, firstName,
			state.firstContentSeen, state.firstReasoningSeen, state.elapsedMs())
	}

	for _, tc := range toolCalls {
		var toolCallIndex int
		if tc.Index != nil {
			toolCallIndex = *tc.Index
		}
		toolCallEntry, exists := state.toolCallMap[toolCallIndex]
		if !exists || toolCallEntry == nil {
			toolCallEntry = &types.LLMToolCall{
				Type: string(tc.Type),
				Function: types.FunctionCall{
					Name:      "",
					Arguments: "",
				},
			}
			state.toolCallMap[toolCallIndex] = toolCallEntry
		}

		if tc.ID != "" {
			toolCallEntry.ID = tc.ID
		}
		if tc.Type != "" {
			toolCallEntry.Type = string(tc.Type)
		}
		if tc.Function.Name != "" {
			// 방어적 검증: 일부 제공자(vLLM Ascend 등)는 각 스트림 chunk마다 전체 도구명을 반복 전송합니다.
			// 현재 저장된 이름과 새 이름이 같으면 중복으로 보고 덧붙이지 않습니다.
			if toolCallEntry.Function.Name != tc.Function.Name {
				toolCallEntry.Function.Name += tc.Function.Name
			}
		}

		argsUpdated := false
		if tc.Function.Arguments != "" {
			toolCallEntry.Function.Arguments += tc.Function.Arguments
			argsUpdated = true
		}

		currName := toolCallEntry.Function.Name
		if currName != "" &&
			currName == state.lastFunctionName[toolCallIndex] &&
			argsUpdated &&
			!state.nameNotified[toolCallIndex] &&
			toolCallEntry.ID != "" {
			streamChan <- types.StreamResponse{
				ResponseType: types.ResponseTypeToolCall,
				Content:      "",
				Done:         false,
				Data: map[string]interface{}{
					"tool_name":    currName,
					"tool_call_id": toolCallEntry.ID,
				},
			}
			state.nameNotified[toolCallIndex] = true
		}

		state.lastFunctionName[toolCallIndex] = currName

		// Stream final_answer tool arguments as answer-type chunks
		if toolCallEntry.Function.Name == "final_answer" && argsUpdated {
			extractor, exists := state.fieldExtractors[toolCallIndex]
			if !exists {
				extractor = newJSONFieldExtractor("answer")
				state.fieldExtractors[toolCallIndex] = extractor
				// Detect non-incremental arrival: if the first args chunk is large,
				// the model likely returned all arguments at once (non-streaming tool call)
				if len(tc.Function.Arguments) > 200 {
					logger.Warnf(ctx, "[LLM Stream] final_answer args arrived in large chunk (%d bytes), "+
						"model may not support incremental tool call streaming", len(tc.Function.Arguments))
				}
			}
			answerChunk := extractor.Feed(tc.Function.Arguments)
			if answerChunk != "" {
				streamChan <- types.StreamResponse{
					ResponseType: types.ResponseTypeAnswer,
					Content:      answerChunk,
					Done:         false,
					Data: map[string]interface{}{
						"source": "final_answer_tool",
					},
				}
			}
		}

		// Stream thinking tool's thought field as thinking-type chunks
		if toolCallEntry.Function.Name == "thinking" && argsUpdated {
			extractor, exists := state.fieldExtractors[toolCallIndex]
			if !exists {
				extractor = newJSONFieldExtractor("thought")
				state.fieldExtractors[toolCallIndex] = extractor
			}
			thoughtChunk := extractor.Feed(tc.Function.Arguments)
			if thoughtChunk != "" {
				streamChan <- types.StreamResponse{
					ResponseType: types.ResponseTypeThinking,
					Content:      thoughtChunk,
					Done:         false,
					Data: map[string]interface{}{
						"source":       "thinking_tool",
						"tool_call_id": toolCallEntry.ID,
					},
				}
			}
		}
	}
}

// GetModelName은 모델 이름을 반환합니다.
func (c *RemoteAPIChat) GetModelName() string {
	return c.modelName
}

// GetModelID는 모델 ID를 반환합니다.
func (c *RemoteAPIChat) GetModelID() string {
	return c.modelID
}

// GetProvider는 provider 이름을 반환합니다.
func (c *RemoteAPIChat) GetProvider() provider.ProviderName {
	return c.provider
}

// GetBaseURL은 baseURL을 반환합니다.
func (c *RemoteAPIChat) GetBaseURL() string {
	return c.baseURL
}

// GetAPIKey는 apiKey를 반환합니다.
func (c *RemoteAPIChat) GetAPIKey() string {
	return c.apiKey
}
