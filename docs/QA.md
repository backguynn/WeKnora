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


## P.S.
위 방법으로 해결되지 않으면 issue에 문제 내용을 남기고, 문제 분석을 위해 필요한 로그 정보를 제공해 주세요.
