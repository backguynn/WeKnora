# 내장 모델 관리 가이드

## 개요

내장 모델은 시스템 레벨의 모델 구성으로, 모든 테넌트에 보이지만 민감 정보는 숨김 처리되며 편집하거나 삭제할 수 없습니다. 내장 모델은 시스템 기본 모델 구성을 제공해 모든 테넌트가 동일한 모델 서비스를 사용할 수 있도록 합니다.

## 내장 모델 특징

- **모든 테넌트 공개**: 내장 모델은 모든 테넌트에 공개되며 별도 설정이 필요 없습니다
- **보안 보호**: 내장 모델의 민감 정보(API Key, Base URL)는 숨김 처리되어 상세 내용을 볼 수 없습니다
- **읽기 전용 보호**: 내장 모델은 편집하거나 삭제할 수 없고 기본 모델로만 설정할 수 있습니다
- **통합 관리**: 시스템 관리자가 일괄 유지 관리하여 구성 일관성과 보안을 보장합니다

## 내장 모델 추가 방법

내장 모델은 데이터베이스에 직접 삽입해야 합니다. 아래 단계를 따라 추가합니다.

### 1. 모델 데이터 준비

먼저 내장 모델로 설정할 모델 구성 정보를 준비합니다.
- 모델 이름(name)
- 모델 타입(type): `KnowledgeQA`, `Embedding`, `Rerank` 또는 `VLLM`
- 모델 소스(source): `local` 또는 `remote`
- 모델 파라미터(parameters): base_url, api_key, provider 등
- 테넌트 ID(tenant_id): 충돌 방지를 위해 10000 이하의 테넌트 ID 사용을 권장합니다

**지원 provider**: `generic`(사용자 정의), `openai`, `aliyun`, `zhipu`, `volcengine`, `hunyuan`, `deepseek`, `minimax`, `mimo`, `siliconflow`, `jina`, `openrouter`, `gemini`, `modelscope`, `moonshot`, `qianfan`, `qiniu`, `longcat`, `gpustack`

### 2. SQL 삽입 실행

다음 SQL로 내장 모델을 삽입합니다.

```sql
-- 예시: LLM 내장 모델 삽입
INSERT INTO models (
    id,
    tenant_id,
    name,
    type,
    source,
    description,
    parameters,
    is_default,
    status,
    is_builtin
) VALUES (
    'builtin-llm-001',                    -- 고정 ID 사용, builtin- 접두사 권장
    10000,                                -- 테넌트 ID(첫 번째 테넌트 사용)
    'GPT-4',                              -- 모델 이름
    'KnowledgeQA',                        -- 모델 타입
    'remote',                             -- 모델 소스
    '내장 LLM 모델',                       -- 설명
    '{"base_url": "https://api.openai.com/v1", "api_key": "sk-xxx", "provider": "openai"}'::jsonb,  -- 파라미터(JSON 형식)
    false,                                -- 기본 모델 여부
    'active',                             -- 상태
    true                                  -- 내장 모델 표시
) ON CONFLICT (id) DO NOTHING;

-- 예시: Embedding 내장 모델 삽입
INSERT INTO models (
    id,
    tenant_id,
    name,
    type,
    source,
    description,
    parameters,
    is_default,
    status,
    is_builtin
) VALUES (
    'builtin-embedding-001',
    10000,
    'text-embedding-ada-002',
    'Embedding',
    'remote',
    '내장 Embedding 모델',
    '{"base_url": "https://api.openai.com/v1", "api_key": "sk-xxx", "provider": "openai", "embedding_parameters": {"dimension": 1536, "truncate_prompt_tokens": 0}}'::jsonb,
    false,
    'active',
    true
) ON CONFLICT (id) DO NOTHING;

-- 예시: Rerank 내장 모델 삽입
INSERT INTO models (
    id,
    tenant_id,
    name,
    type,
    source,
    description,
    parameters,
    is_default,
    status,
    is_builtin
) VALUES (
    'builtin-rerank-001',
    10000,
    'bge-reranker-base',
    'Rerank',
    'remote',
    '내장 Rerank 모델',
    '{"base_url": "https://api.jina.ai/v1", "api_key": "jina-xxx", "provider": "jina"}'::jsonb,
    false,
    'active',
    true
) ON CONFLICT (id) DO NOTHING;

-- 예시: VLLM 내장 모델 삽입
INSERT INTO models (
    id,
    tenant_id,
    name,
    type,
    source,
    description,
    parameters,
    is_default,
    status,
    is_builtin
) VALUES (
    'builtin-vllm-001',
    10000,
    'gpt-4-vision',
    'VLLM',
    'remote',
    '내장 VLLM 모델',
    '{"base_url": "https://dashscope.aliyuncs.com/compatible-mode/v1", "api_key": "sk-xxx", "provider": "aliyun"}'::jsonb,
    false,
    'active',
    true
) ON CONFLICT (id) DO NOTHING;
```

### 3. 삽입 결과 확인

다음 SQL로 내장 모델이 정상 삽입됐는지 확인합니다.

```sql
SELECT id, name, type, is_builtin, status 
FROM models 
WHERE is_builtin = true
ORDER BY type, created_at;
```

## 주의 사항

1. **ID 명명 규칙**: `builtin-{type}-{번호}` 형식을 권장합니다. 예: `builtin-llm-001`, `builtin-embedding-001`
2. **테넌트 ID**: 내장 모델은 임의 테넌트에 속할 수 있지만, 첫 번째 테넌트 ID(보통 10000)를 권장합니다
3. **파라미터 형식**: `parameters` 필드는 유효한 JSON 형식이어야 합니다
4. **멱등성**: `ON CONFLICT (id) DO NOTHING`으로 반복 실행 시 오류를 방지합니다
5. **보안**: 내장 모델의 API Key와 Base URL은 프런트엔드에서 자동 숨김 처리되지만, DB 원본 데이터는 유지되므로 접근 권한을 안전하게 관리하세요

## 기존 모델을 내장 모델로 설정

이미 존재하는 모델을 내장 모델로 설정하려면 다음 UPDATE를 사용합니다.

```sql
UPDATE models 
SET is_builtin = true 
WHERE id = '모델ID' AND name = '모델이름';
```

## 내장 모델 제거

내장 모델 표시를 제거(일반 모델로 복원)하려면 다음을 실행합니다.

```sql
UPDATE models 
SET is_builtin = false 
WHERE id = '모델ID';
```

참고: 내장 모델 표시를 제거하면 일반 모델로 돌아가 편집 및 삭제가 가능합니다.

