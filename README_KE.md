<p align="center">
  <picture>
    <img src="./docs/images/logo.png" alt="WeKnora Logo" height="120"/>
  </picture>
</p>

<p align="center">
  <picture>
    <a href="https://trendshift.io/repositories/15289" target="_blank">
      <img src="https://trendshift.io/api/badge/repositories/15289" alt="Tencent%2FWeKnora | Trendshift" style="width: 250px; height: 55px;" width="250" height="55"/>
    </a>
  </picture>
</p>
<p align="center">
    <a href="https://weknora.weixin.qq.com" target="_blank">
        <img alt="공식 홈페이지" src="https://img.shields.io/badge/%EA%B3%B5%EC%8B%9D%20%ED%99%88%ED%8E%98%EC%9D%B4%EC%A7%80-WeKnora-4e6b99">
    </a>
    <a href="https://chatbot.weixin.qq.com" target="_blank">
        <img alt="WeChat 대화 오픈 플랫폼" src="https://img.shields.io/badge/WeChat%20%EB%8C%80%ED%99%94%20%EC%98%A4%ED%94%88%20%ED%94%8C%EB%9E%AB%ED%8F%BC-5ac725">
    </a>
    <a href="https://github.com/Tencent/WeKnora/blob/main/LICENSE">
        <img src="https://img.shields.io/badge/License-MIT-ffffff?labelColor=d4eaf7&color=2e6cc4" alt="License">
    </a>
    <a href="./CHANGELOG.md">
        <img alt="Version" src="https://img.shields.io/badge/version-0.3.0-2e6cc4?labelColor=d4eaf7">
    </a>
</p>

<p align="center">
| <a href="./README.md"><b>English</b></a> | <a href="./README_CN.md"><b>간체 중국어</b></a> | <a href="./README_JA.md"><b>日本語</b></a> | <b>한국어</b> |
</p>

<p align="center">
  <h4 align="center">

  [개요](#-개요) • [아키텍처](#-아키텍처) • [주요-기능](#-주요-기능) • [시작하기](#-시작하기) • [API-참조](#-api-참조) • [개발자-가이드](#-개발자-가이드)
  
  </h4>
</p>

# 💡 WeKnora - LLM 기반 문서 이해 및 검색 프레임워크

## 📌 개요

[**WeKnora**](https://weknora.weixin.qq.com)는 복잡하고 이질적인 문서를 처리하기 위한 깊은 문서 이해와 의미 기반 검색에 최적화된 LLM 기반 프레임워크입니다.

멀티모달 전처리, 의미 벡터 인덱싱, 지능형 검색, 대규모 언어 모델 추론을 결합한 모듈형 아키텍처를 채택합니다. WeKnora의 핵심은 **RAG (Retrieval-Augmented Generation)** 패러다임을 따르며, 관련 문서 청크와 모델 추론을 결합해 고품질의 문맥 기반 답변을 제공합니다.

**웹사이트:** https://weknora.weixin.qq.com

## ✨ 최신 업데이트

**v0.3.0 하이라이트:**

- 🏢 **공유 스페이스**: 멤버 초대 기반 공유 스페이스, 멤버 간 지식 베이스/에이전트 공유, 테넌트 격리 검색
- 🧩 **에이전트 스킬**: 스마트 추론 에이전트를 위한 사전 탑재 스킬과 보안 격리용 샌드박스 실행 환경
- 🤖 **커스텀 에이전트**: 지식 베이스 선택 모드(전체/지정/비활성)와 함께 커스텀 에이전트 생성/구성/선택 지원
- 📊 **데이터 분석 에이전트**: CSV/Excel 분석을 위한 DataSchema 도구가 포함된 내장 데이터 분석 에이전트
- 🧠 **사고 모드**: LLM/에이전트 사고 모드 지원, 사고 내용 지능형 필터링
- 🔍 **웹 검색 제공자**: DuckDuckGo 외에 Bing과 Google 검색 제공자 추가
- 📋 **강화된 FAQ**: 배치 임포트 드라이런, 유사 질문, 검색 결과의 매칭 질문, 대용량 임포트 오브젝트 스토리지 오프로딩
- 🔑 **API 키 인증**: Swagger 문서 보안을 위한 API 키 인증 메커니즘
- 📎 **입력 중 선택**: 입력창에서 @mention 표시로 지식 베이스와 파일 직접 선택
- ☸️ **Helm 차트**: Neo4j GraphRAG 지원을 포함한 Kubernetes 배포용 완전한 Helm 차트
- 🌍 **i18n**: 한국어 지원 추가
- 🔒 **보안 강화**: SSRF 안전 HTTP 클라이언트, 강화된 SQL 검증, MCP stdio 전송 보안, 샌드박스 기반 실행
- ⚡ **인프라**: Qdrant 벡터 DB 지원, Redis ACL, 로그 레벨 설정, Ollama 임베딩 최적화, `DISABLE_REGISTRATION` 제어

**v0.2.0 하이라이트:**

- 🤖 **에이전트 모드**: 내장 도구, MCP 도구, 웹 검색을 호출할 수 있는 새로운 ReACT 에이전트 모드로, 여러 번의 반복과 반성을 통해 종합 요약 리포트를 제공
- 📚 **다양한 지식 베이스 유형**: FAQ 및 문서 지식 베이스 유형 지원, 폴더/URL 임포트, 태그 관리, 온라인 입력 등 신규 기능 제공
- ⚙️ **대화 전략**: 에이전트 모델, 일반 모드 모델, 검색 임계값, 프롬프트 구성 지원으로 멀티턴 대화 제어 강화
- 🌐 **웹 검색**: 내장 DuckDuckGo 검색 엔진을 포함한 확장 가능한 웹 검색 엔진 지원
- 🔌 **MCP 도구 통합**: MCP를 통해 에이전트 기능 확장, uvx 및 npx 런처 내장, 여러 전송 방식 지원
- 🎨 **새 UI**: 에이전트/일반 모드 전환, 도구 호출 과정 표시, 지식 베이스 관리 UI 개선
- ⚡ **인프라 업그레이드**: MQ 비동기 작업 관리, 자동 DB 마이그레이션, 빠른 개발 모드 지원

## 🔒 보안 안내

**중요:** v0.1.3부터 WeKnora에는 로그인 인증 기능이 포함되어 보안이 강화되었습니다. 운영 환경 배포 시 다음을 권장합니다:

- WeKnora 서비스를 공용 인터넷이 아닌 내부/프라이빗 네트워크 환경에 배포
- 정보 유출 위험 방지를 위해 공용 네트워크에 직접 노출하지 않기
- 배포 환경에 맞는 방화벽 규칙 및 접근 제어 구성
- 보안 패치 및 개선 사항을 위해 최신 버전으로 정기 업데이트

## 🏗️ 아키텍처

![weknora-architecture.png](./docs/images/architecture.png)

WeKnora는 완전한 문서 이해 및 검색 파이프라인을 구축하기 위해 현대적인 모듈형 설계를 채택합니다. 시스템은 문서 파싱, 벡터 처리, 검색 엔진, 대형 모델 추론을 핵심 모듈로 구성하며, 각 컴포넌트는 유연하게 설정 및 확장할 수 있습니다.

## 🎯 주요 기능

- **🤖 에이전트 모드**: 내장 도구로 지식 베이스, MCP 도구, 웹 검색을 활용하는 ReACT 에이전트 모드 지원
- **🔍 정밀 이해**: PDF, Word, 이미지 등 다양한 형식에서 구조화된 콘텐츠를 추출해 통합 의미 뷰로 정리
- **🧠 지능형 추론**: LLM로 문서 문맥과 사용자 의도를 파악하여 정확한 Q&A 및 멀티턴 대화 제공
- **📚 다양한 지식 베이스 유형**: FAQ/문서 지식 베이스 지원, 폴더 임포트, URL 임포트, 태그 관리, 온라인 입력
- **🔧 유연한 확장성**: 파싱, 임베딩, 검색, 생성 단계가 분리되어 커스터마이징 용이
- **⚡ 효율적 검색**: 키워드/벡터/지식 그래프를 결합한 하이브리드 검색과 교차 지식 베이스 검색 지원
- **🌐 웹 검색**: 내장 DuckDuckGo 검색 엔진을 포함한 확장 가능한 웹 검색 엔진 지원
- **🔌 MCP 도구 통합**: MCP를 통한 에이전트 기능 확장, uvx 및 npx 런처 내장, 3가지 전송 방식 지원
- **⚙️ 대화 전략**: 에이전트/일반 모드 모델, 검색 임계값, 프롬프트 설정으로 멀티턴 대화 제어
- **🎯 사용자 친화성**: 직관적인 웹 인터페이스와 표준화된 API 제공
- **🔒 보안 및 제어**: 로컬 배포 및 프라이빗 클라우드 지원으로 데이터 주권 보장

## 📊 적용 시나리오

| 시나리오 | 적용 사례 | 핵심 가치 |
|---------|----------|----------|
| **엔터프라이즈 지식 관리** | 내부 문서 검색, 정책 Q&A, 운영 매뉴얼 검색 | 지식 발견 효율 향상, 교육 비용 절감 |
| **학술 연구 분석** | 논문 검색, 연구 보고서 분석, 학술 자료 정리 | 문헌 조사 가속, 연구 의사결정 지원 |
| **제품 기술 지원** | 제품 매뉴얼 Q&A, 기술 문서 검색, 트러블슈팅 | 고객 서비스 품질 향상, 지원 부담 감소 |
| **법무/컴플라이언스 검토** | 계약 조항 검색, 규정 정책 검색, 사례 분석 | 컴플라이언스 효율 향상, 법적 리스크 감소 |
| **의료 지식 지원** | 의학 문헌 검색, 치료 가이드라인 검색, 사례 분석 | 임상 의사결정 지원, 진단 품질 향상 |

## 🧩 기능 매트릭스

| 모듈 | 지원 | 설명 |
|---------|--------------------------------------------------------------------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| 에이전트 모드 | ✅ ReACT 에이전트 모드 | 내장 도구로 지식 베이스, MCP 도구, 웹 검색을 활용하며 교차 지식 베이스 검색과 여러 번의 반복을 지원 |
| 지식 베이스 유형 | ✅ FAQ / Document | FAQ 및 문서 지식 베이스 유형 지원, 폴더/URL 임포트, 태그 관리, 온라인 입력 |
| 문서 형식 | ✅ PDF / Word / Txt / Markdown / Images (with OCR / Caption) | 이미지 텍스트 추출을 포함한 구조화/비구조화 문서 지원 |
| 모델 관리 | ✅ 중앙 집중 구성, 내장 모델 공유 | 지식 베이스 설정 내 모델 선택을 포함한 중앙 관리, 멀티테넌트 내장 모델 공유 지원 |
| 임베딩 모델 | ✅ 로컬 모델, BGE / GTE APIs 등 | 로컬 배포 및 클라우드 벡터 생성 API와 호환되는 커스텀 임베딩 모델 |
| 벡터 DB 통합 | ✅ PostgreSQL (pgvector), Elasticsearch | 주요 벡터 인덱스 백엔드 지원, 시나리오별 유연한 전환 |
| 검색 전략 | ✅ BM25 / Dense Retrieval / GraphRAG | 희소/밀집 검색과 지식 그래프 강화 검색을 결합한 검색-리랭크-생성 파이프라인 구성 가능 |
| LLM 통합 | ✅ Qwen, DeepSeek 등 지원, 사고/비사고 모드 전환 | Ollama 등 로컬 모델 또는 외부 API 서비스와 유연한 추론 설정 지원 |
| 대화 전략 | ✅ 에이전트 모델, 일반 모드 모델, 검색 임계값, 프롬프트 구성 | 에이전트/일반 모드 모델 및 검색 임계값, 온라인 프롬프트 구성으로 멀티턴 대화 제어 |
| 웹 검색 | ✅ 확장 가능한 검색 엔진, DuckDuckGo / Google | 내장 DuckDuckGo 검색 엔진을 포함한 확장 가능한 웹 검색 엔진 지원 |
| MCP 도구 | ✅ uvx, npx 런처, Stdio/HTTP Streamable/SSE | MCP를 통해 에이전트 기능 확장, uvx/npx 런처 내장, 3가지 전송 방식 지원 |
| QA 기능 | ✅ 문맥 기반 멀티턴 대화, 프롬프트 템플릿 | 복잡한 의미 모델링, 지시 제어, 체인-오브-생각 Q&A, 컨텍스트 윈도우 구성 지원 |
| E2E 테스트 | ✅ 검색+생성 과정 시각화 및 지표 평가 | 리콜 적중률, 답변 커버리지, BLEU/ROUGE 등 지표 평가용 종단간 테스트 도구 |
| 배포 모드 | ✅ 로컬 배포 / Docker 이미지 | 프라이빗/오프라인 배포와 유연한 운영 요구 사항 충족, 빠른 개발 모드 지원 |
| 사용자 인터페이스 | ✅ Web UI + RESTful API | 에이전트/일반 모드 전환, 도구 호출 과정 표시가 가능한 대화형 인터페이스 및 표준 API |
| 작업 관리 | ✅ MQ 비동기 작업, 자동 DB 마이그레이션 | MQ 기반 비동기 작업 상태 유지, 버전 업그레이드 시 자동 스키마 및 데이터 마이그레이션 |

## 🚀 시작하기

### 🛠 사전 준비

다음 도구가 설치되어 있어야 합니다:

* [Docker](https://www.docker.com/)
* [Docker Compose](https://docs.docker.com/compose/)
* [Git](https://git-scm.com/)

### 📦 설치

#### 1. 리포지토리 클론

```bash
# 메인 리포지토리 클론
git clone https://github.com/Tencent/WeKnora.git
cd WeKnora
```

#### 2. 환경 변수 구성

```bash
# 예제 env 파일 복사
cp .env.example .env

# .env 편집 후 필수 값 설정
# 모든 변수는 .env.example 주석에 설명되어 있습니다
```

#### 3. 서비스 시작(Ollama 포함)

.env 파일에서 시작해야 할 이미지를 확인하세요.

```bash
./scripts/start_all.sh
```

또는

```bash
make start-all
```

#### 3.0 Ollama 서비스 시작(선택)

```bash
ollama serve > /dev/null 2>&1 &
```

#### 3.1 기능 조합별 활성화

- 최소 핵심 서비스
```bash
docker compose up -d
```

- 모든 기능 활성화
```bash
docker-compose --profile full up -d
```

- 트레이싱 로그 필요
```bash
docker-compose --profile jaeger up -d
```

- Neo4j 지식 그래프 필요
```bash
docker-compose --profile neo4j up -d
```

- Minio 파일 스토리지 필요
```bash
docker-compose --profile minio up -d
```

- 여러 옵션 조합
```bash
docker-compose --profile neo4j --profile minio up -d
```

#### 4. 서비스 중지

```bash
./scripts/start_all.sh --stop
# 또는
make stop-all
```

### 🌐 서비스 접속

서비스 시작 후 다음 주소로 접근합니다:

* Web UI: `http://localhost`
* Backend API: `http://localhost:8080`
* Jaeger Tracing: `http://localhost:16686`

### 🔌 WeChat 대화 오픈 플랫폼에서 WeKnora 사용

WeKnora는 [WeChat Dialog Open Platform](https://chatbot.weixin.qq.com)의 핵심 기술 프레임워크로, 다음과 같은 사용 방식을 제공합니다:

- **노코드 배포**: 지식을 업로드하면 WeChat 생태계에서 즉시 Q&A 서비스를 배포하여 "묻고 답하는" 경험을 제공
- **효율적인 질문 관리**: 빈번한 질문을 분류/관리할 수 있는 풍부한 데이터 도구로 정확하고 신뢰 가능한 답변 유지
- **WeChat 생태계 통합**: WeChat Dialog Open Platform을 통해 WeChat Official Accounts, Mini Programs 등과 자연스럽게 연동

### 🔗 MCP Server로 WeKnora 접근

#### 1. 리포지토리 클론
```
git clone https://github.com/Tencent/WeKnora
```

#### 2. MCP Server 구성
> 구성은 [MCP Configuration Guide](./mcp-server/MCP_CONFIG.md) 참조를 권장합니다.

MCP 클라이언트를 구성하여 서버에 연결하세요:
```json
{
  "mcpServers": {
    "weknora": {
      "args": [
        "path/to/WeKnora/mcp-server/run_server.py"
      ],
      "command": "python",
      "env":{
        "WEKNORA_API_KEY":"WeKnora 인스턴스에서 개발자 도구를 열고 요청 헤더의 x-api-key(sk로 시작)를 확인해 입력",
        "WEKNORA_BASE_URL":"http(s)://your-weknora-address/api/v1"
      }
    }
  }
}
```

stdio 명령으로 바로 실행:
```
pip install weknora-mcp-server
python -m weknora-mcp-server
```

## 🔧 초기화 구성 가이드

모델 설정을 빠르게 구성하고 시행착오를 줄이기 위해, 기존 설정 파일 초기화 방식에 Web UI 기반 모델 구성 인터페이스를 추가했습니다. 사용 전에 최신 코드로 업데이트되어 있는지 확인하세요. 구체적인 단계는 아래와 같습니다.
처음 사용한다면 1, 2단계는 건너뛰고 3, 4단계부터 진행하세요.

### 1. 서비스 중지

```bash
./scripts/start_all.sh --stop
```

### 2. 기존 데이터 테이블 정리(중요 데이터가 없을 때 권장)

```bash
make clean-db
```

### 3. 서비스 컴파일 및 시작

```bash
./scripts/start_all.sh
```

### 4. Web UI 접속

http://localhost

처음 접속하면 등록/로그인 페이지로 자동 이동합니다. 등록 후 새 지식 베이스를 만들고 구성 페이지에서 관련 설정을 완료하세요.

## 📱 인터페이스 소개

### Web UI 인터페이스

<table>
  <tr>
    <td><b>지식 베이스 관리</b><br/><img src="./docs/images/knowledgebases.png" alt="Knowledge Base Management"></td>
    <td><b>대화 설정</b><br/><img src="./docs/images/settings.png" alt="Conversation Settings"></td>
  </tr>
  <tr>
    <td colspan="2"><b>에이전트 모드 도구 호출 과정</b><br/><img src="./docs/images/agent-qa.png" alt="Agent Mode Tool Call Process"></td>
  </tr>
</table>

**지식 베이스 관리:** FAQ 및 문서 지식 베이스 유형을 생성하고 드래그 앤 드롭, 폴더 임포트, URL 임포트 등 다양한 방법으로 가져올 수 있습니다. 문서 구조를 자동 식별하고 핵심 지식을 추출하여 인덱스를 구축합니다. 태그 관리와 온라인 입력을 지원하며, 처리 진행 상황과 문서 상태를 명확히 표시해 효율적인 지식 베이스 관리를 제공합니다.

**에이전트 모드:** 내장 도구로 지식 베이스를 검색하고, 사용자가 구성한 MCP 도구와 웹 검색 도구를 호출하여 여러 번의 반복과 반성을 통해 종합 요약 리포트를 제공합니다. 교차 지식 베이스 검색을 지원하여 여러 지식 베이스를 동시에 선택해 검색할 수 있습니다.

**대화 전략:** 에이전트/일반 모드 모델, 검색 임계값, 온라인 프롬프트 구성으로 멀티턴 대화 제어와 검색 실행 방식을 정밀하게 제어합니다. 대화 입력창은 에이전트/일반 모드 전환, 웹 검색 활성화/비활성화, 대화 모델 선택을 지원합니다.

### 문서 지식 그래프

WeKnora는 문서를 지식 그래프로 변환하여 문서의 서로 다른 섹션 간 관계를 표시합니다. 지식 그래프 기능이 활성화되면 시스템이 내부 의미 연관 네트워크를 분석/구축하여 문서 이해를 돕고 검색 결과의 관련성과 범위를 확장합니다.

자세한 설정은 [Knowledge Graph Configuration Guide](./docs/KnowledgeGraph.md)를 참고하세요.

### MCP Server

필요한 설정은 [MCP Configuration Guide](./mcp-server/MCP_CONFIG.md)를 참고하세요.

## 📘 API 참조

문제 해결 FAQ: [Troubleshooting FAQ](./docs/QA.md)

상세 API 문서: [API Docs](./docs/api/README.md)

제품 계획 및 예정 기능: [Roadmap](./docs/ROADMAP.md)

## 🧭 개발자 가이드

### ⚡ 빠른 개발 모드(권장)

코드를 자주 수정해야 한다면 **매번 Docker 이미지를 다시 빌드할 필요가 없습니다**. 빠른 개발 모드를 사용하세요:

```bash
# 방법 1: Make 명령 사용(권장)
make dev-start      # 인프라 시작
make dev-app        # 백엔드 시작(새 터미널)
make dev-frontend   # 프론트엔드 시작(새 터미널)

# 방법 2: 원클릭 시작
./scripts/quick-dev.sh

# 방법 3: 스크립트 사용
./scripts/dev.sh start     # 인프라 시작
./scripts/dev.sh app       # 백엔드 시작(새 터미널)
./scripts/dev.sh frontend  # 프론트엔드 시작(새 터미널)
```

**개발 장점:**
- ✅ 프론트엔드 수정 즉시 핫 리로드(재시작 불필요)
- ✅ 백엔드 수정 빠른 재시작(5-10초, Air 핫 리로드 지원)
- ✅ Docker 이미지 재빌드 불필요
- ✅ IDE 브레이크포인트 디버깅 지원

**상세 문서:** [개발 환경 빠른 시작](./docs/开发指南.md)

### 📁 디렉터리 구조

```
WeKnora/
├── client/      # go client
├── cmd/         # Main entry point
├── config/      # Configuration files
├── docker/      # docker images files
├── docreader/   # Document parsing app
├── docs/        # Project documentation
├── frontend/    # Frontend app
├── internal/    # Core business logic
├── mcp-server/  # MCP server
├── migrations/  # DB migration scripts
└── scripts/     # Shell scripts
```

## 🤝 기여하기

커뮤니티 기여를 환영합니다. 제안, 버그, 기능 요청은 [Issue](https://github.com/Tencent/WeKnora/issues)를 제출하거나 Pull Request를 만들어 주세요.

### 🎯 기여 방법

- 🐛 **버그 수정**: 시스템 결함 발견 및 수정
- ✨ **신규 기능**: 새로운 기능 제안 및 구현
- 📚 **문서 개선**: 프로젝트 문서 개선
