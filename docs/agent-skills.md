# Agent Skills 문서

## 개요

Agent Skills는 Agent가 "사용 설명서"를 읽어 새 능력을 학습하는 확장 메커니즘입니다. 기존의 하드코딩 도구와 달리 Skills는 System Prompt에 주입되어 Agent의 능력을 확장하며 **Progressive Disclosure(점진적 공개)** 설계 원칙을 따릅니다.
현재는 **지능형 추론** 능력이 있는 Agent만 지원합니다. 프런트엔드에서 Agent 편집 페이지에 관련 설정이 있습니다.

### 핵심 특성

- **비침투적 확장**: 기존 Agent ReAct 흐름에 영향을 주지 않음
- **온디맨드 로딩**: 3단계 점진 로딩으로 Token 사용 최적화
- **샌드박스 실행**: 스크립트를 격리 환경에서 안전하게 실행
- **유연한 구성**: 다중 디렉터리, 화이트리스트 필터 지원

## 설계 철학

### Progressive Disclosure(점진적 공개)

Skills는 3단계 로딩 메커니즘을 사용해 필요할 때만 LLM에 상세 정보를 제공합니다:

```
┌─────────────────────────────────────────────────────────────────┐
│ Level 1: 메타데이터 (Metadata)                                  │
│ • 항상 System Prompt에 로드                                     │
│ • 약 100 tokens/skill                                           │
│ • 포함: 스킬 이름 + 간단한 설명                                  │
└─────────────────────────────────────────────────────────────────┘
                              ↓ 사용자 요청이 매칭될 때
┌─────────────────────────────────────────────────────────────────┐
│ Level 2: 지침 (Instructions)                                    │
│ • read_skill 도구로 온디맨드 로드                                │
│ • SKILL.md의 지침 내용                                           │
│ • 포함: 상세 지침, 코드 예시, 사용 방법                           │
└─────────────────────────────────────────────────────────────────┘
                              ↓ 추가 정보가 필요할 때
┌─────────────────────────────────────────────────────────────────┐
│ Level 3: 추가 리소스 (Resources)                                │
│ • read_skill 도구로 특정 파일 로드                               │
│ • 보조 문서, 구성 템플릿, 스크립트 파일                           │
│ • execute_skill_script로 스크립트 실행                           │
└─────────────────────────────────────────────────────────────────┘
```

## Skill 디렉터리 구조

각 Skill은 `SKILL.md` 본문 파일과 선택적 부가 리소스를 포함하는 디렉터리입니다:

```
my-skill/
├── SKILL.md           # 필수: 본문 파일(YAML frontmatter 포함)
├── REFERENCE.md       # 선택: 보조 문서
├── templates/         # 선택: 템플릿 파일
│   └── config.yaml
└── scripts/           # 선택: 실행 스크립트
    ├── analyze.py
    └── generate.sh
```

## SKILL.md 형식

### YAML Frontmatter

각 `SKILL.md`는 YAML frontmatter로 시작하며 메타데이터를 정의합니다:

```markdown
---
name: pdf-processing
description: Extract text and tables from PDF files, fill forms, merge documents. Use when working with PDF files or when the user mentions PDFs, forms, or document extraction.
---

# PDF Processing

This skill provides utilities for working with PDF documents.

## Quick Start

Use pdfplumber to extract text from PDFs:

```python
import pdfplumber

with pdfplumber.open("document.pdf") as pdf:
    text = pdf.pages[0].extract_text()
    print(text)
```

## 메타데이터 검증 규칙

| 필드 | 요구 사항 |
|------|----------|
| `name` | 1-50자, 한글/영문/숫자만 허용, 예약어 불가 |
| `description` | 1-500자, 스킬 용도와 트리거 조건 설명 |

**예약어**: `system`, `default`, `internal`, `core`, `base`, `root`, `admin`

## 구성

### AgentConfig 구성 항목

```go
type AgentConfig struct {
    // ... 기타 구성 ...

    // Skills 관련 구성
    SkillsEnabled  bool     `json:"skills_enabled"`   // Skills 활성화 여부
    SkillDirs      []string `json:"skill_dirs"`       // Skill 디렉터리 목록
    AllowedSkills  []string `json:"allowed_skills"`   // 화이트리스트(비어 있으면 모두 허용)
}
```

### 구성 예시

```json
{
  "skills_enabled": true,
  "skill_dirs": [
    "/path/to/project/skills",
    "/home/user/.agent-skills"
  ],
  "allowed_skills": ["pdf-processing", "code-review"]
}
```

### Sandbox 구성(환경 변수)

Sandbox 관련 구성은 환경 변수로 설정합니다:

| 환경 변수 | 설명 | 기본값 |
|---------|------|--------|
| `WEKNORA_SANDBOX_MODE` | sandbox 모드: `docker`, `local`, `disabled` | `disabled` |
| `WEKNORA_SANDBOX_TIMEOUT` | 스크립트 실행 타임아웃(초) | `60` |
| `WEKNORA_SANDBOX_DOCKER_IMAGE` | 사용자 정의 Docker 이미지 | `wechatopenai/weknora-sandbox:latest` |

### Sandbox 모드

| 모드 | 설명 |
|------|------|
| `docker` | Docker 컨테이너 격리(권장) |
| `local` | 로컬 프로세스 실행(기본 보안 제한) |
| `disabled` | 스크립트 실행 비활성화 |

## Agent 도구

Skills 기능은 두 가지 도구로 Agent와 상호작용합니다:

### read_skill

스킬 내용 또는 특정 파일을 읽습니다.

**파라미터**:
```json
{
  "skill_name": "pdf-processing",      // 필수: 스킬 이름
  "file_path": "FORMS.md"              // 선택: 상대 경로
}
```

**사용 시나리오**:
1. Level 2 로드: `skill_name`만 전달
2. Level 3 리소스 로드: `skill_name`과 `file_path`를 함께 전달

**예시 호출**:
```json
// 스킬 본문 로드
{"skill_name": "pdf-processing"}

// 보조 문서 로드
{"skill_name": "pdf-processing", "file_path": "FORMS.md"}

// 스크립트 내용 확인
{"skill_name": "pdf-processing", "file_path": "scripts/analyze.py"}
```

### execute_skill_script

샌드박스에서 스킬 스크립트를 실행합니다.

**파라미터**:
```json
{
  "skill_name": "pdf-processing",           // 필수: 스킬 이름
  "script_path": "scripts/analyze.py",      // 필수: 스크립트 상대 경로
  "args": ["input.pdf", "--format", "json"] // 선택: 명령줄 파라미터
}
```

**지원 스크립트 타입**:
- Python (`.py`)
- Shell (`.sh`)
- JavaScript/Node.js (`.js`)
- Ruby (`.rb`)
- Go (`.go`)

## 프리로드 스킬(Preloaded Skills)

시스템에는 지식베이스 QA와 문서 처리 능력을 강화하기 위한 5개의 프리로드 스킬이 내장되어 있습니다:

### 1. citation-generator - 인용 생성기

**용도**: 표준 인용 형식을 자동 생성

**트리거 시나리오**:
- 참고문헌 생성이 필요할 때
- 지식베이스 내용의 출처를 표기할 때
- 인용 정보를 요청받았을 때

**핵심 기능**:
| 기능 | 설명 |
|------|------|
| 출처 표기 | 답변에 사용된 각 지식의 출처 표기 |
| 인용 형식화 | APA, MLA, Chicago, 간소화 형식 지원 |
| 참고문헌 목록 | 답변 끝에 전체 참고문헌 목록 생성 |

**간소화 인용 형식 예시**:
```
회사 정책[직원핸드북2024.pdf, 15쪽]에 따르면 연차 신청은 사전에...
```

---

### 2. data-processor - 데이터 처리기

**용도**: 데이터 처리 및 분석

**트리거 시나리오**:
- "이 데이터를 분석해줘", "통계 내줘", "합계/평균을 계산해줘"
- "JSON/CSV 형식으로 변환해줘"
- "핵심 정보를 추출해줘", "표로 정리해줘"
- "보고서 생성", "데이터 요약"

**핵심 기능**:
| 기능 | 설명 |
|------|------|
| 데이터 분석 | 검색된 문서 데이터를 통계적으로 분석 |
| 형식 변환 | JSON/CSV/Markdown 간 변환 |
| 데이터 추출 | 비정형 텍스트에서 구조화 정보 추출 |
| 보고서 생성 | 데이터 분석 보고서와 요약 생성 |

**사용 가능한 스크립트**:
- `scripts/analyze.py` - 데이터 분석 스크립트
- `scripts/format_converter.py` - 형식 변환 스크립트
- `scripts/extract_info.py` - 정보 추출 스크립트

**스크립트 사용 예시**:
```bash
# 데이터 분석
echo '{"items": [1, 2, 3, 4, 5]}' | python scripts/analyze.py

# 형식 변환(JSON -> CSV)
echo '[{"name": "A", "value": 1}]' | python scripts/format_converter.py --to csv

# 정보 추출
echo "2024년 매출은 100만 위안" | python scripts/extract_info.py
```

---

### 3. doc-coauthoring - 문서 협업(Claude 공식 Skill에서 파생)

**용도**: 사용자가 구조화된 문서를 작성하도록 안내

**트리거 시나리오**:
- 문서 작성: "write a doc", "draft a proposal", "create a spec"
- 문서 유형: PRD, 설계 문서, 의사결정 문서, RFC

**워크플로**:

```
Stage 1: 컨텍스트 수집 (Context Gathering)
        ↓
Stage 2: 정교화 및 구조화 (Refinement & Structure)
        ↓
Stage 3: 독자 테스트 (Reader Testing)
```

**3단계 설명**:
| 단계 | 목표 | 핵심 활동 |
|------|------|----------|
| Stage 1 | 사용자와 Claude 사이의 정보 격차 축소 | 메타 정보 질문, 컨텍스트 수집, 문제 명확화 |
| Stage 2 | 섹션별 문서 구성 | 브레인스토밍, 선별 정리, 반복 수정 |
| Stage 3 | 문서의 독자 효과 테스트 | 독자 질문 예측, 서브 에이전트 테스트, 사각지대 보완 |

---

### 4. document-analyzer - 문서 분석기

**용도**: 문서 구조와 내용을 심층 분석

**트리거 시나리오**:
- 문서 구조 분석
- 핵심 정보 추출
- 문서 유형 식별
- 내용 품질 평가

**핵심 기능**:
| 기능 | 설명 |
|------|------|
| 구조 분석 | 문서의 장/절 계층 및 조직 구조 식별 |
| 핵심 정보 추출 | 핵심 주장, 주요 데이터, 중요한 결론 추출 |
| 문서 유형 식별 | 문서 유형(보고서, 매뉴얼, 논문, 계약서 등) 판별 |
| 내용 품질 평가 | 문서의 완전성, 일관성, 가독성 평가 |

**분석 흐름**:
1. **문서 개요** - 문서 기본 정보 획득
2. **구조 분석** - 제목 계층과 섹션 구조 식별
3. **내용 추출** - 핵심 주제, 주요 주장, 근거 데이터 추출
4. **품질 평가** - 완전성, 일관성, 명확성 평가

---

### 스킬 디렉터리 구조

프리로드 스킬은 `skills/preloaded/` 디렉터리에 있습니다:

```
skills/preloaded/
├── citation-generator/
│   └── SKILL.md
├── data-processor/
│   ├── SKILL.md
│   └── scripts/
│       ├── analyze.py
│       ├── format_converter.py
│       └── extract_info.py
├── doc-coauthoring/
│   └── SKILL.md
├── document-analyzer/
│   └── SKILL.md
└── summary-generator/
    └── SKILL.md
```

## 사용자 정의 Skill 생성

현재는 사용자가 직접 사용자 정의 Skill을 만드는 기능을 지원하지 않습니다.

## 샌드박스 보안 메커니즘

### 스크립트 보안 검증(Script Validator)

스크립트 실행 전, 시스템은 다층 보안 검증을 수행해 잠재적 악성 동작을 차단합니다:

#### 검증 유형

| 유형 | 설명 | 예시 |
|------|------|------|
| **위험 명령 탐지** | 시스템을 손상시킬 수 있는 명령 탐지 | `rm -rf /`, `mkfs`, `shutdown`, 포크 봄 |
| **위험 패턴 매칭** | 정규식으로 고위험 동작 패턴 매칭 | `curl \| bash`, `base64 -d`, `eval()` |
| **네트워크 접근 탐지** | 네트워크 요청 시도 탐지 | `curl`, `wget`, `socket.connect`, `requests.get` |
| **리버스 Shell 탐지** | 원격 제어 백도어 탐지 | `/dev/tcp/`, `bash -i`, `nc -e` |
| **파라미터 인젝션 탐지** | 명령줄 파라미터 인젝션 탐지 | `&&`, `\|`, `$()`, 백틱 |
| **Stdin 인젝션 탐지** | 표준 입력 내 삽입 명령 탐지 | 명령 치환 구문 |

#### 차단되는 위험 명령

**시스템 파괴 계열**:
- `rm -rf /`, `rm -rf /*` - 루트 디렉터리 재귀 삭제
- `mkfs`, `dd if=/dev/zero` - 파일 시스템/디스크 작업
- 포크 봄: `:(){ :|:& };:`

**시스템 제어 계열**:
- `shutdown`, `reboot`, `halt`, `poweroff`
- `killall`, `pkill`
- `systemctl`, `service`

**권한 상승 계열**:
- `chmod 777 /`, `chown root`
- `setuid`, `setgid`, `passwd`
- `/etc/passwd`, `/etc/shadow`, `/etc/sudoers` 접근

**자격 증명 탈취 계열**:
- `.ssh/`, `id_rsa`, `id_ed25519` 접근
- 민감 구성 파일 읽기

**컨테이너 탈출 계열**:
- `docker`, `kubectl`, `nsenter`
- `unshare`, `capsh`

#### 차단되는 위험 패턴

**코드 인젝션**:
```
# 아래 패턴은 차단됩니다
curl ... | bash           # 다운로드 후 실행
wget ... | sh             # 다운로드 후 실행
eval()                    # 동적 코드 실행
exec()                    # 명령 실행
os.system()               # 시스템 명령 실행
subprocess.Popen(shell=True)  # Shell 명령 실행
```

**인코딩 우회 시도**:
```
# 아래 패턴은 차단됩니다
base64 -d                 # Base64 디코딩 실행
echo ... | base64 -d      # 파이프 디코딩
xxd -r                    # Hex 디코딩
```

**Python 특유 위험**:
```python
# 아래 패턴은 차단됩니다
__import__()              # 동적 임포트
pickle.load()             # 역직렬화(임의 코드 실행 가능)
yaml.load()               # 안전하지 않은 YAML 로드
yaml.unsafe_load()        # 명시적으로 안전하지 않은 로드
```

#### Shell 연산자 차단

파라미터에 다음 연산자가 포함되면 차단됩니다:

| 연산자 | 설명 |
|--------|------|
| `&&`, `\|\|` | 명령 연결 |
| `;` | 명령 구분 |
| `\|` | 파이프 |
| `$()`, `` ` `` | 명령 치환 |
| `>`, `>>`, `<` | 리다이렉션 |
| `2>`, `&>` | 오류/복합 리다이렉션 |
| `\n`, `\r` | 줄바꿈 인젝션 |

#### 검증 결과

검증 실패 시 상세 오류 정보를 반환합니다:

```go
type ValidationError struct {
    Type    string // 오류 유형: dangerous_command, dangerous_pattern, arg_injection 등
    Pattern string // 매칭된 패턴
    Context string // 컨텍스트 정보
    Message string // 사람이 읽을 수 있는 설명
}
```

**예시 오류**:
```
security validation failed [dangerous_command]: Script contains dangerous command: rm -rf / (pattern: rm -rf /, context: ...cleanup && rm -rf / && echo done...)
```

#### 사용 예시

```go
// 검증기 생성
validator := sandbox.NewScriptValidator()

// 스크립트 내용 검증
result := validator.ValidateScript(scriptContent)
if !result.Valid {
    for _, err := range result.Errors {
        log.Printf("Security error: %s", err.Error())
    }
    return errors.New("script validation failed")
}

// 명령줄 파라미터 검증
argsResult := validator.ValidateArgs(args)

// 표준 입력 검증
stdinResult := validator.ValidateStdin(stdin)

// 또는 한번에 전체 검증
fullResult := validator.ValidateAll(scriptContent, args, stdin)
```

---

### Docker 샌드박스

Docker 모드는 가장 강력한 격리를 제공합니다:

- **non-root 사용자**: 컨테이너 내부에서 일반 사용자로 실행
- **Capability 제한**: 모든 Linux capabilities 제거
- **읽기 전용 파일 시스템**: 루트 파일 시스템을 읽기 전용으로 설정
- **리소스 제한**: 메모리 256MB, CPU 제한
- **네트워크 격리**: 기본적으로 네트워크 접근 없음
- **임시 마운트**: Skill 디렉터리를 읽기 전용으로 마운트
- **스크립트 사전 검증**: 실행 전 보안 검증 수행

#### 샌드박스 이미지

시스템은 전용 샌드박스 이미지 `wechatopenai/weknora-sandbox`를 사용하며, Python 3.11, Node.js 20, 일반적인 CLI 도구와 Python 라이브러리를 미리 포함합니다. 실행 시 의존성을 임시 설치할 필요가 없습니다.

**이미지 사전 풀링**(첫 배포 시 권장, 첫 실행 시 다운로드 대기를 방지):

```bash
# 방법 1: 직접 풀링
docker pull wechatopenai/weknora-sandbox:latest

# 방법 2: 로컬 빌드
sh scripts/build_images.sh -s
```

> 사전 풀링을 하지 않으면 애플리케이션 시작 시 자동으로 비동기 풀링(`EnsureImage`)이 실행되지만, 첫 실행에는 다운로드 대기가 필요할 수 있습니다.

**이미지 내장 환경**:
- Python 3.11 + pip(requests, pyyaml, pandas, beautifulsoup4)
- Node.js 20 + npm
- CLI 도구: jq, curl, bash, grep, sed, awk 등

```bash
# Docker 실행 예시
docker run --rm \
  --user 1000:1000 \
  --cap-drop ALL \
  --read-only \
  --memory=256m \
  --network=none \
  -v /path/to/skill:/skill:ro \
  -w /skill \
  wechatopenai/weknora-sandbox:latest \
  python scripts/analyze.py input.pdf
```

### Local 샌드박스

Local 모드는 기본적인 보호를 제공합니다:

- **명령 화이트리스트**: 지정된 인터프리터만 허용
- **작업 디렉터리 제한**: Skill 디렉터리로 제한
- **환경 변수 필터링**: 안전한 변수만 전달
- **타임아웃 제어**: 기본 30초 타임아웃
- **경로 탐색 방지**: Skill 디렉터리 외부 접근 차단
- **스크립트 사전 검증**: 실행 전 보안 검증 수행

**허용되는 명령**:
- `python`, `python3`
- `node`, `nodejs`
- `bash`, `sh`
- `ruby`
- `go run`

## API 참조

### SkillManager

```go
type Manager interface {
    // 초기화, 모든 Skills 탐색
    Initialize(ctx context.Context) error
    
    // 모든 Skill 메타데이터(Level 1) 가져오기
    GetAllMetadata() []*SkillMetadata
    
    // Skill 지침(Level 2) 로드
    LoadSkill(ctx context.Context, skillName string) (*Skill, error)
    
    // Skill 파일 내용(Level 3) 읽기
    ReadSkillFile(ctx context.Context, skillName, filePath string) (string, error)
    
    // Skill 내 모든 파일 목록
    ListSkillFiles(ctx context.Context, skillName string) ([]string, error)
    
    // Skill 스크립트 실행
    ExecuteScript(ctx context.Context, skillName, scriptPath string, args []string) (*sandbox.ExecuteResult, error)
    
    // 활성화 여부 확인
    IsEnabled() bool
}
```

### Skill 구조

```go
type Skill struct {
    Name         string // 스킬 이름
    Description  string // 스킬 설명
    BasePath     string // 디렉터리 절대 경로
    FilePath     string // SKILL.md 절대 경로
    Instructions string // SKILL.md 본문 지침 내용
    Loaded       bool   // Level 2 로드 여부
}

type SkillMetadata struct {
    Name        string // 스킬 이름
    Description string // 스킬 설명
    BasePath    string // 디렉터리 경로
}
```

### ExecuteResult 구조

```go
type ExecuteResult struct {
    ExitCode int           // 종료 코드
    Stdout   string        // 표준 출력
    Stderr   string        // 표준 오류
    Duration time.Duration // 실행 시간
    Error    error         // 실행 오류
}
```

## 예시: 전체 워크플로

다음은 Agent가 사용자 요청을 처리하는 전체 흐름입니다:

```
사용자: "report.pdf에서 표 데이터 추출해줘"

Agent 생각:
  → System Prompt의 Skills 목록 확인
  → "pdf-processing" 스킬 매칭 확인

Agent 행동 1: read_skill 호출
  → {"skill_name": "pdf-processing"}
  → SKILL.md 지침 내용 확인
  → pdfplumber 사용 방법 학습

Agent 행동 2: execute_skill_script 호출
  → {"skill_name": "pdf-processing", 
     "script_path": "scripts/extract_text.py",
     "args": ["report.pdf"]}
  → 샌드박스에서 스크립트 실행, 표 데이터 반환

Agent 응답:
  → 사용자에게 추출된 표 데이터 제공
  → 데이터 활용 방법 제안
```

## 문제 해결

### Skill이 발견되지 않는 경우

1. `skill_dirs` 설정이 올바른지 확인
2. 디렉터리에 `SKILL.md` 파일이 있는지 확인
3. YAML frontmatter 형식 검증

```bash
# demo 실행으로 검증
go run ./cmd/skills-demo/main.go
```

### 스크립트 실행 실패

1. `sandbox_mode` 설정 확인
2. Docker 모드: Docker 서비스가 실행 중인지 확인
3. Local 모드: 인터프리터 설치 여부 확인
4. 스크립트 권한 및 문법 확인

### 메타데이터 검증 오류

자주 발생하는 오류:
- `skill name too long`: 이름이 50자를 초과함
- `skill name contains invalid characters`: 허용되지 않는 문자가 포함됨
- `skill name is reserved`: 예약어를 사용함
- `skill description too long`: 설명이 500자를 초과함