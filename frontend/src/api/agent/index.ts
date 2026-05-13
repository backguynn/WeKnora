import { get, post, put, del } from "../../utils/request";

// 에이전트 설정
// 스마트 추론 모드에서 사용하는 에이전트 유형 프리셋 ID
// 'rag-qa'           : 고전적인 문서/FAQ 청크 RAG
// 'wiki-qa'          : Wiki 그래프 탐색형 질의응답
// 'hybrid-rag-wiki'  : Wiki + 청크 혼합 검색
// 'custom'           : 완전 사용자 정의(프리셋 미적용)
export type AgentType = 'rag-qa' | 'wiki-qa' | 'hybrid-rag-wiki' | 'data-analysis' | 'custom';

export interface CustomAgentConfig {
  // ===== 기본 설정 =====
  agent_mode?: 'quick-answer' | 'smart-reasoning';  // 실행 모드: quick-answer=RAG 모드, smart-reasoning=ReAct Agent 모드
  // 스마트 추론 모드의 유형 프리셋. 시스템 프롬프트 + 도구 + KB 호환성 조합을 한 번에 적용합니다.
  // agent_mode === 'smart-reasoning' 일 때만 사용되며 quick-answer 모드에서는 무시됩니다.
  agent_type?: AgentType;
  system_prompt?: string;           // 통합 시스템 프롬프트({{web_search_status}} 플레이스홀더로 동작 제어)
  system_prompt_id?: string;        // 참조하는 prompt template ID(프리셋이 자동으로 채움)
  context_template?: string;        // 컨텍스트 템플릿(일반 모드)
  context_template_id?: string;

  // ===== 模型设置 =====
  model_id?: string;
  rerank_model_id?: string;         // ReRank 模型 ID
  temperature?: number;
  max_completion_tokens?: number;   // 最大生成token数（普通模式）
  thinking?: boolean;                      // 是否启用思考模式（支持扩展思考的模型）

  // ===== Agent 모드 설정 =====
  max_iterations?: number;          // 최대 반복 횟수
  llm_call_timeout?: number;        // LLM 호출 타임아웃(초)
  allowed_tools?: string[];         // 허용된 도구
  reflection_enabled?: boolean;     // 반성 기능 사용 여부
  // MCP 서비스 선택 모드: all=전체 활성화된 MCP 서비스, selected=지정 서비스, none=사용 안 함
  mcp_selection_mode?: 'all' | 'selected' | 'none';
  mcp_services?: string[];          // 선택한 MCP 서비스 ID 목록

  // ===== Skills设置（仅Agent模式）=====
  // Skills选择模式：all=全部预装, selected=指定, none=不使用
  skills_selection_mode?: 'all' | 'selected' | 'none';
  selected_skills?: string[];       // 选择的Skill名称列表

  // ===== 知识库设置 =====
  // 知识库选择模式：all=全部知识库, selected=指定知识库, none=不使用知识库
  kb_selection_mode?: 'all' | 'selected' | 'none';
  knowledge_bases?: string[];
  // 是否仅在显式 @ 提及时检索知识库（默认: false）
  // true: 只有用户通过 @ 明确提及知识库/文档时才检索
  // false: 根据 kb_selection_mode 自动检索知识库
  retrieve_kb_only_when_mentioned?: boolean;

  // ===== 图片上传/多模态设置 =====
  image_upload_enabled?: boolean;    // 是否启用图片上传（默认: false）
  vlm_model_id?: string;            // VLM模型ID（图片分析用）
  image_storage_provider?: string;   // 图片存储提供商
  audio_upload_enabled?: boolean;    // 是否启用音频上传/ASR转录（默认: false）
  asr_model_id?: string;            // ASR模型ID（音频转录用）

  // ===== 文件类型限制 =====
  // 支持的文件类型（如 ["csv", "xlsx", "xls"]）
  // 为空表示支持所有文件类型
  supported_file_types?: string[];

  // ===== 网络搜索设置 =====
  web_search_enabled?: boolean;
  web_search_provider_id?: string;
  web_search_max_results?: number;
  web_fetch_enabled?: boolean;
  web_fetch_top_n?: number;

  // ===== 多轮对话设置 =====
  multi_turn_enabled?: boolean;     // 是否启用多轮对话
  history_turns?: number;           // 保留历史轮数

  // ===== FAQ策略设置 =====
  faq_priority_enabled?: boolean;
  faq_direct_answer_threshold?: number;
  faq_score_boost?: number;

  // ===== 检索策略设置 =====
  embedding_top_k?: number;         // 向量召回TopK
  keyword_threshold?: number;       // 关键词召回阈值
  vector_threshold?: number;        // 向量召回阈值
  rerank_top_k?: number;            // 重排TopK
  rerank_threshold?: number;        // 重排阈值

  // ===== 高级设置（主要用于普通模式）=====
  enable_query_expansion?: boolean; // 是否启用查询扩展
  enable_rewrite?: boolean;         // 是否启用问题改写
  rewrite_prompt_system?: string;   // 改写系统提示词
  rewrite_prompt_user?: string;     // 改写用户提示词模板
  fallback_strategy?: 'fixed' | 'model'; // 兜底策略
  fallback_response?: string;       // 固定兜底回复
  fallback_prompt?: string;         // 兜底提示词（模型生成时）

  // ===== 已废弃字段（保留兼容）=====
  welcome_message?: string;
  suggested_prompts?: string[];
}

// 智能体
export interface CustomAgent {
  id: string;
  name: string;
  description?: string;
  avatar?: string;
  is_builtin: boolean;
  tenant_id?: number;
  created_by?: string;
  config: CustomAgentConfig;
  created_at?: string;
  updated_at?: string;
}

// 创建智能体请求
export interface CreateAgentRequest {
  name: string;
  description?: string;
  avatar?: string;
  config?: CustomAgentConfig;
}

// 更新智能体请求
export interface UpdateAgentRequest {
  name: string;
  description?: string;
  avatar?: string;
  config?: CustomAgentConfig;
}

// 内置智能体 ID（常用的保留常量，便于代码引用）
export const BUILTIN_QUICK_ANSWER_ID = 'builtin-quick-answer';
export const BUILTIN_SMART_REASONING_ID = 'builtin-smart-reasoning';

// AgentMode 常量
export const AGENT_MODE_QUICK_ANSWER = 'quick-answer';
export const AGENT_MODE_SMART_REASONING = 'smart-reasoning';

// Deprecated: Use BUILTIN_QUICK_ANSWER_ID instead
export const BUILTIN_AGENT_NORMAL_ID = BUILTIN_QUICK_ANSWER_ID;
// Deprecated: Use BUILTIN_SMART_REASONING_ID instead
export const BUILTIN_AGENT_AGENT_ID = BUILTIN_SMART_REASONING_ID;

// 获取智能体列表（包括内置智能体）
// disabled_own_agent_ids: 当前租户在对话下拉中停用的「我的」智能体 ID，仅影响本租户
export function listAgents() {
  return get<{ data: CustomAgent[]; disabled_own_agent_ids?: string[] }>('/api/v1/agents');
}

// 获取智能体详情
export function getAgentById(id: string) {
  return get<{ data: CustomAgent }>(`/api/v1/agents/${id}`);
}

// 创建智能体
export function createAgent(data: CreateAgentRequest) {
  return post<{ data: CustomAgent }>('/api/v1/agents', data);
}

// 更新智能体
export function updateAgent(id: string, data: UpdateAgentRequest) {
  return put<{ data: CustomAgent }>(`/api/v1/agents/${id}`, data);
}

// 删除智能体
export function deleteAgent(id: string) {
  return del<{ success: boolean }>(`/api/v1/agents/${id}`);
}

// 复制智能体
export function copyAgent(id: string) {
  return post<{ data: CustomAgent }>(`/api/v1/agents/${id}/copy`);
}

// 判断是否为内置智能体（通过 agent.is_builtin 字段或 ID 前缀判断）
export function isBuiltinAgent(agentId: string): boolean {
  return agentId.startsWith('builtin-');
}

// 占位符定义
export interface PlaceholderDefinition {
  name: string;
  label: string;
  description: string;
}

// 占位符响应
export interface PlaceholdersResponse {
  all: PlaceholderDefinition[];
  system_prompt: PlaceholderDefinition[];
  agent_system_prompt: PlaceholderDefinition[];
  context_template: PlaceholderDefinition[];
  rewrite_system_prompt: PlaceholderDefinition[];
  rewrite_prompt: PlaceholderDefinition[];
  fallback_prompt: PlaceholderDefinition[];
}

// 获取占位符定义
export function getPlaceholders() {
  return get<{ data: PlaceholdersResponse }>('/api/v1/agents/placeholders');
}

// ===== 智能体类型预设 =====

// 后端 kb_filter 结构（见 internal/types/agent_type_preset.go）
export interface AgentTypeKBFilter {
  any_of?: string[];   // KB 至少拥有其一
  all_of?: string[];   // KB 必须全部拥有
  none_of?: string[];  // KB 必须全部不拥有
}

// KB 能力标签（后端 types.KBCapabilities 的 JSON）
export interface KBCapabilities {
  vector: boolean;
  keyword: boolean;
  wiki: boolean;
  graph: boolean;
  faq: boolean;
}

// 预设的"自动填充"配置载荷：仅包含被预设覆盖的字段；其他字段不动
export interface AgentTypePresetConfig {
  system_prompt_id?: string;
  temperature?: number;
  max_iterations?: number;
  allowed_tools?: string[];
  retain_retrieval_history?: boolean;
  faq_priority_enabled?: boolean;
  web_search_enabled?: boolean;
  supported_file_types?: string[];
  kb_selection_mode?: 'all' | 'selected' | 'none';
}

export interface AgentTypePresetI18n {
  label: string;
  description: string;
}

export interface AgentTypePreset {
  id: AgentType;
  i18n: Record<string, AgentTypePresetI18n>;
  config?: AgentTypePresetConfig;     // 为空表示"自定义"类型（无预设）
  kb_filter?: AgentTypeKBFilter;      // 为空表示所有 KB 可选
}

// 拉取类型预设列表（编辑器用）
export function getAgentTypePresets() {
  return get<{ data: AgentTypePreset[] }>('/api/v1/agents/type-presets');
}

// ===== IM渠道 =====

export interface IMChannel {
  id: string;
  tenant_id?: number;
  agent_id: string;
  platform: 'wecom' | 'feishu' | 'slack' | 'telegram' | 'dingtalk' | 'mattermost' | 'wechat';
  name: string;
  enabled: boolean;
  mode: 'webhook' | 'websocket' | 'longpoll';
  output_mode: 'stream' | 'full';
  session_mode?: 'user' | 'thread';
  knowledge_base_id?: string;
  credentials: Record<string, any>;
  created_at?: string;
  updated_at?: string;
}

export function listIMChannels(agentId: string) {
  return get<{ data: IMChannel[] }>(`/api/v1/agents/${agentId}/im-channels`);
}

// Tenant-wide overview row. Credentials are intentionally omitted — use
// listIMChannels(agentId) when you need to edit a specific channel.
export interface IMChannelOverview {
  id: string;
  tenant_id: number;
  agent_id: string;
  agent_name: string; // empty string for built-in agents
  platform: IMChannel['platform'];
  name: string;
  enabled: boolean;
  mode: IMChannel['mode'];
  output_mode: IMChannel['output_mode'];
  session_mode?: IMChannel['session_mode'];
  bot_identity: string;
  created_at: string;
  updated_at: string;
}

export function listAllIMChannels() {
  return get<{ data: IMChannelOverview[] }>('/api/v1/im-channels');
}

export function createIMChannel(agentId: string, data: Partial<IMChannel>) {
  return post<{ data: IMChannel }>(`/api/v1/agents/${agentId}/im-channels`, data);
}

export function updateIMChannel(id: string, data: Partial<IMChannel>) {
  return put<{ data: IMChannel }>(`/api/v1/im-channels/${id}`, data);
}

export function deleteIMChannel(id: string) {
  return del<{ success: boolean }>(`/api/v1/im-channels/${id}`);
}

export function toggleIMChannel(id: string) {
  return post<{ data: IMChannel }>(`/api/v1/im-channels/${id}/toggle`);
}

// ===== 推荐问题 =====

// 推荐问题
export interface SuggestedQuestion {
  question: string;
  source: 'faq' | 'document' | 'agent_config' | 'wiki';
  knowledge_base_id?: string;
}

// 获取智能体推荐问题
// 根据智能体关联的知识库范围返回推荐问题，用于前端对话面板快捷提问
export function getSuggestedQuestions(
  agentId: string,
  params?: { knowledge_base_ids?: string[]; knowledge_ids?: string[]; limit?: number }
) {
  const query = new URLSearchParams();
  if (params?.knowledge_base_ids?.length) query.set('knowledge_base_ids', params.knowledge_base_ids.join(','));
  if (params?.knowledge_ids?.length) query.set('knowledge_ids', params.knowledge_ids.join(','));
  if (params?.limit) query.set('limit', String(params.limit));
  const qs = query.toString();
  return get<{ data: { questions: SuggestedQuestion[] } }>(`/api/v1/agents/${agentId}/suggested-questions${qs ? '?' + qs : ''}`);
}
// ===== WeChat QR Code Login =====

export interface WeChatQRCodeResult {
  qrcode_url: string;
  qrcode: string;
}

export interface WeChatQRCodeStatus {
  status: 'wait' | 'scaned' | 'confirmed' | 'expired';
  credentials?: {
    bot_token: string;
    ilink_bot_id: string;
    ilink_user_id: string;
  };
  baseurl?: string;
}

export function getWeChatQRCode() {
  return post<{ data: WeChatQRCodeResult }>('/api/v1/wechat/qrcode');
}

export function pollWeChatQRCodeStatus(qrcode: string) {
  return post<{ data: WeChatQRCodeStatus }>('/api/v1/wechat/qrcode/status', { qrcode });
}
