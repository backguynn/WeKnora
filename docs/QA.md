# FAQ

## 1. 로그는 어떻게 확인하나요?
```bash
docker compose logs -f app docreader postgres
```

## 2. 서비스를 어떻게 시작/중지하나요?
```bash
# 서비스 시작
./scripts/start_all.sh

# 서비스 중지
./scripts/start_all.sh --stop

# 데이터베이스 초기화
./scripts/start_all.sh --stop && make clean-db
```

## 3. 서비스 시작 후 문서 업로드가 정상 동작하지 않나요?

대부분 Embedding 모델과 대화 모델이 올바르게 설정되지 않은 경우입니다. 아래 절차로 점검하세요.

1. `.env`의 모델 설정이 완전한지 확인하세요. 로컬 모델을 ollama로 접근하는 경우, 로컬 ollama 서비스가 정상 실행 중인지 확인하고 아래 환경 변수를 올바르게 설정해야 합니다:
```bash
# LLM Model
INIT_LLM_MODEL_NAME=your_llm_model
# Embedding Model
INIT_EMBEDDING_MODEL_NAME=your_embedding_model
# Embedding 모델 벡터 차원
INIT_EMBEDDING_MODEL_DIMENSION=your_embedding_model_dimension
# Embedding 모델 ID(보통 문자열)
INIT_EMBEDDING_MODEL_ID=your_embedding_model_id
```

remote API로 모델을 호출하는 경우 해당 `BASE_URL`과 `API_KEY`를 추가로 제공해야 합니다:
```bash
# LLM 모델 접근 주소
INIT_LLM_MODEL_BASE_URL=your_llm_model_base_url
# LLM 모델 API 키(인증이 필요하면 설정)
INIT_LLM_MODEL_API_KEY=your_llm_model_api_key
# Embedding 모델 접근 주소
INIT_EMBEDDING_MODEL_BASE_URL=your_embedding_model_base_url
# Embedding 모델 API 키(인증이 필요하면 설정)
INIT_EMBEDDING_MODEL_API_KEY=your_embedding_model_api_key
```

재정렬 기능이 필요할 경우 Rerank 모델을 추가로 설정해야 합니다. 설정 예시는 다음과 같습니다:
```bash
# 사용할 Rerank 모델 이름
INIT_RERANK_MODEL_NAME=your_rerank_model_name
# Rerank 모델 접근 주소
INIT_RERANK_MODEL_BASE_URL=your_rerank_model_base_url
# Rerank 모델 API 키(인증이 필요하면 설정)
INIT_RERANK_MODEL_API_KEY=your_rerank_model_api_key
```

2. 메인 서비스 로그에 `ERROR` 로그가 출력되는지 확인하세요.

## 4. 이미지가 없거나 잘못된 이미지 링크가 표시되나요?

멀티모달 기능 사용 시 이미지가 표시되지 않거나 잘못된 링크가 보이는 경우, 아래 절차로 점검하세요.

### 1. 멀티모달 기능이 올바르게 설정되었는지 확인

지식베이스 설정에서 **고급 설정 - 멀티모달 기능**을 켜고, 화면에서 해당 멀티모달 모델을 설정하세요.

### 2. MinIO 서비스가 실행 중인지 확인

멀티모달 기능이 MinIO 스토리지를 사용하는 경우, MinIO 이미지가 정상적으로 실행 중이어야 합니다:

```bash
# MinIO 서비스 시작
docker-compose --profile minio up -d

# 또는 전체 서비스 시작(MinIO, Jaeger, Neo4j, Qdrant 포함)
docker-compose --profile full up -d
```

### 3. MinIO Bucket 권한 확인

MinIO의 해당 bucket에 적절한 읽기/쓰기 권한이 있는지 확인하세요:

1. MinIO 콘솔 접속: `http://localhost:9001` (기본 포트)
2. `.env`에 설정한 `MINIO_ACCESS_KEY_ID`와 `MINIO_SECRET_ACCESS_KEY`로 로그인
3. 해당 bucket으로 들어가 접근 정책을 **공개 읽기** 또는 **공개 읽기/쓰기**로 설정

**중요 안내**:
- Bucket 이름에 특수 문자를 포함하지 마세요(중국어 포함). 소문자, 숫자, 하이픈 사용을 권장합니다.
- 기존 bucket 권한을 변경할 수 없다면 설정에 존재하지 않는 bucket 이름을 입력하세요. 프로젝트가 자동으로 bucket을 생성하고 권한을 설정합니다.

### 4. MINIO_PUBLIC_ENDPOINT 설정

`docker-compose.yml`에서 `MINIO_PUBLIC_ENDPOINT` 기본값은 `http://localhost:9000`입니다.

**중요 안내**: 다른 장치 또는 컨테이너에서 이미지에 접근해야 한다면 `localhost`가 동작하지 않을 수 있으므로, 이를 호스트의 실제 IP 주소로 변경하세요.


## 5. 플랫폼 호환성 안내

**중요 안내**: `OCR_BACKEND=paddle` 모드는 일부 플랫폼에서 정상 동작하지 않을 수 있습니다. PaddleOCR 시작 실패 시 다음 해결 방법을 선택하세요.

### 방법 1: OCR 인식 비활성화

`docker-compose.yml`의 `docreader` 서비스에서 `OCR_BACKEND` 설정을 제거한 뒤 docreader 서비스를 재시작하세요.

**주의**: `no_ocr`로 설정하면 문서 파싱에서 OCR 기능을 사용하지 않습니다. 이미지/스캔 문서의 문자 인식 정확도에 영향을 줄 수 있습니다.

### 방법 2: 외부 OCR 모델 사용(권장)

OCR 기능이 필요하다면 외부 시각 언어 모델(VLM)로 PaddleOCR을 대체할 수 있습니다. `docker-compose.yml`의 `docreader` 서비스에 다음을 설정하세요:

```yaml
environment:
  - OCR_BACKEND=vlm
  - OCR_API_BASE_URL=${OCR_API_BASE_URL:-}
  - OCR_API_KEY=${OCR_API_KEY:-}
  - OCR_MODEL=${OCR_MODEL:-}
```

그 다음 docreader 서비스를 재시작하세요.

**장점**: 외부 OCR 모델은 더 나은 인식 성능을 제공하며 플랫폼 제한이 적습니다.

## 6. 데이터 분석 기능은 어떻게 사용하나요?

데이터 분석 기능을 사용하기 전에 에이전트에 관련 도구가 설정되어 있는지 확인하세요:

1. **지능형 추론**: 도구 설정에서 아래 두 도구를 선택해야 합니다:
   - 데이터 메타 정보 조회
   - 데이터 분석

2. **빠른 Q&A 에이전트**: 도구를 수동으로 선택하지 않아도 간단한 데이터 조회를 수행할 수 있습니다.

### 주의 사항 및 사용 규칙

1. **지원 파일 형식**
   - 현재 **CSV** (`.csv`) 및 **Excel** (`.xlsx`, `.xls`) 파일만 지원합니다.
   - 복잡한 Excel 파일에서 읽기 실패 시 표준 CSV 형식으로 변환해 다시 업로드하는 것을 권장합니다.

2. **쿼리 제한**
   - **읽기 전용 쿼리**만 지원하며 `SELECT`, `SHOW`, `DESCRIBE`, `EXPLAIN`, `PRAGMA` 등을 포함합니다.
   - `INSERT`, `UPDATE`, `DELETE`, `CREATE`, `DROP` 등 데이터 변경 작업은 금지됩니다.

## 7. 방금 저장한 설정이 몇 초 뒤에 다시 사라집니다. 왜 그런가요?

이 문제는 실제로 설정이 시스템에서 지워진 경우보다, 브라우저 프록시·캐시·확장 프로그램 간섭으로 인해 프런트엔드가 비정상 응답을 읽고 이후 오래된 상태로 다시 덮어쓰는 경우가 많습니다.

다음 순서대로 점검해 보세요.

1. 브라우저 프록시, 패킷 캡처 도구, 요청을 자동 수정하는 확장 프로그램을 먼저 끄고 페이지를 다시 엽니다.
2. 브라우저가 `localhost` 또는 현재 접속 도메인을 프록시로 보내고 있지 않은지 확인합니다. PAC를 사용 중이라면 `localhost`, `127.0.0.1`, 실제 배포 도메인을 모두 직접 연결 목록에 추가하세요.
3. 강력 새로고침을 하거나 시크릿 창에서 다시 로그인한 뒤 설정을 한 번 더 저장해 봅니다.
4. 브라우저 개발자 도구의 `Network` 패널에서 설정 저장 요청이 최신 응답을 돌려주는지, 프록시 재작성·캐시 적중·다른 환경으로의 리디렉션이 없는지 확인합니다.
5. 디버그 모드 배포라면 `app` 서비스를 재시작한 뒤 다시 확인해 볼 수 있습니다.

```bash
docker compose restart app
```

재시작 직후 잠시 정상으로 돌아오지만 다시 같은 현상이 반복된다면, 백엔드 설정 유실로 단정하기보다 브라우저 프록시·캐시·다중 환경 혼선 문제를 먼저 의심하는 편이 맞습니다.

## 8. SSRF 검사 화이트리스트(`SSRF_WHITELIST`)

선택형 설정입니다. `.env`에서 `SSRF_WHITELIST`를 지정하면 URL 검사 등에서 특정 대상을 화이트리스트에 추가해 일반 SSRF 제한을 우회할 수 있습니다. 값은 쉼표로 구분한 규칙 목록이며, 각 항목은 다음 형태를 지원합니다.

- **정확한 도메인**: 예 `api.internal`
- **와일드카드 도메인**: 예 `*.example.com`
- **IPv4**: 예 `203.0.113.5`
- **IPv6**: 예 `2001:db8::1`(대괄호 제외)
- **CIDR**: 예 `10.0.0.0/8`, `2001:db8::/32`

화이트리스트에 포함된 대상은 URL 검사 등에서 일반 SSRF 규칙을 우회합니다. **운영 환경에서는 반드시 신중히 설정**하고, 실제로 필요하며 신뢰할 수 있는 대상만 추가하세요.

예시(`.env.example`과 동일하며 필요에 따라 주석을 해제하고 수정하세요):

```bash
# SSRF_WHITELIST=internal.service,*.corp.example,172.16.0.0/12,2001:db8::1,fd00::/8
```


## 9. Langfuse 관측성 추적은 어떻게 켜고 확인하나요?

WeKnora는 Langfuse를 통해 Agent의 ReAct 루프, 대형 모델 Token 사용량, 도구 호출, 비동기 작업 파이프라인을 전체 추적할 수 있습니다.

**활성화 절차**:
1. 사용 가능한 Langfuse 인스턴스를 준비합니다. 클라우드형과 자체 호스팅형 모두 지원합니다.
2. `.env` 파일에 다음 환경 변수를 설정합니다.
```bash
LANGFUSE_PUBLIC_KEY=pk-lf-...
LANGFUSE_SECRET_KEY=sk-lf-...
LANGFUSE_HOST=https://cloud.langfuse.com # 또는 자체 호스팅 주소
```
3. 서비스를 재시작하면 시스템이 지원되는 모든 모델 호출과 Agent 실행 흐름을 자동 추적합니다. Langfuse의 Traces 패널에서 각 대화와 백그라운드 작업의 세부 실행 폭포도 및 Token 통계를 확인할 수 있습니다.

## 10. Wiki 모드는 무엇이고 어떻게 사용하나요?

Wiki 모드는 Agent가 원본 문서를 바탕으로 구조화되고 서로 연결된 Markdown Wiki 지식베이스를 자동 생성·유지할 수 있게 해 주며, 복잡한 지식을 체계적으로 축적하고 그래프 형태로 연결하는 데 적합합니다.

**사용 방법**:
1. 해당 **지식베이스 설정**으로 들어가 **인덱싱 전략(Indexing Strategy)** 을 엽니다.
2. **Wiki** 인덱스를 켭니다. 필요하면 **지식 그래프**도 함께 활성화할 수 있습니다.
3. 이 지식베이스에 문서를 업로드하면, 시스템이 비동기 작업을 자동 실행하여 대형 모델로 문서의 엔티티와 핵심 개념을 추출하고 구조화된 Wiki 페이지 및 페이지 간 지식 그래프 링크를 생성합니다.
4. 해당 지식베이스의 `Wiki` 탭에서 전용 Wiki 브라우저로 페이지를 조회·관리할 수 있으며, 시각화된 지식 그래프를 통해 서로 다른 내용의 관계를 확인할 수 있습니다.

## P.S.
위 방법으로 해결되지 않으면 issue에 문제 내용을 남기고, 문제 분석을 위해 필요한 로그 정보를 제공해 주세요.
