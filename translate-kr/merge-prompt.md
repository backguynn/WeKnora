# WeKnora 한글화 병합 프롬프트

아래 문서는 Tencent/WeKnora 원본 업데이트를 현재 한국어 포크에 재병합할 때 그대로 재사용할 수 있도록, 실제 병합 과정에서 확인한 규칙과 함정을 반영해 정리한 실행용 프롬프트다. `TARGET_UPSTREAM_REF`만 원하는 upstream ref로 바꿔서 사용하면 된다.

## 실전 메모

- 기준 base는 `e57deb2`, 한글화 커밋은 `b94551c`다.
- 보호 대상 파일 목록의 원본은 [translate-kr/translated-files.md](translate-kr/translated-files.md)다.
- 이번 병합에서 upstream 구조 변경으로 아래 파일은 삭제 유지하는 것이 맞았다.
  - `config/prompt_templates/rewrite_system.yaml`
  - `config/prompt_templates/rewrite_user.yaml`
  - `docs/WeKnora.md`
  - `internal/models/chat/lkeap.go`
  - `internal/models/chat/qwen.go`
- 위 파일들은 과거 한글화 대상이었더라도, upstream가 새 구조로 대체했다면 억지로 복구하지 말고 대체 파일 기준으로 한국어를 재적용해야 한다.
- 병합 후 별도 UI 스윕이 꼭 필요했다. 특히 아래는 하드코딩 중국어가 남기 쉬운 위치였다.
  - `frontend/src/i18n/index.ts`: 기본 locale / fallbackLocale
  - `frontend/src/views/auth/Login.vue`: 언어 선택기 표시명
  - `frontend/src/views/knowledge/components/FAQEntryManager.vue`: 샘플 데이터, CSV/Excel 헤더, import alias
  - `frontend/src/components/document-preview.vue`: 빈 문서 안내 문구
  - `frontend/src/views/settings/ParserEngineSettings.vue`: `t()` 기본 fallback 문구
- 이번 병합에서는 하드코딩 문자열보다 locale 키 누락이 더 큰 리스크였다. `zh-CN.ts` 대비 `ko-KR.ts` 키 diff를 한 번 돌려서 신규 키 누락 여부를 먼저 확인하는 편이 안전했다.
- 특히 아래 계열은 업스트림 기능이 들어오면 한국어 locale에서 빠지거나 영어로 남기 쉬웠다.
   - `frontend/src/i18n/locales/ko-KR.ts`: `datasource.*`, `auth.oidc*`, `agent.editor.web*`, `agentEditor.desc.web*`, `tenant.api.desktop*`
   - `frontend/src/views/agent/AgentEditorModal.vue`: web search provider / web fetch 신규 키 사용 여부
- FAQ CSV/Excel 예제 헤더를 한국어로 바꾸면 import 쪽도 한국어 헤더를 읽도록 같이 수정해야 한다. 표시 문자열만 바꾸면 다시 import가 깨진다.
- 한국어 포크에서는 frontend i18n 기본 언어를 `ko-KR`로 두는 것이 안전하다. 그렇지 않으면 신규 키 누락 시 중국어 fallback이 바로 노출된다.

## 복사용 프롬프트

```text
당신은 Tencent/WeKnora 원본 저장소의 최신 업데이트를 현재 한글화 포크에 병합하는 AI 코딩 에이전트다.

배경 정보:
- 원본 기준 커밋: e57deb2
- 한글화 커밋: b94551c
- 한글화 파일 인벤토리 문서: translate-kr/translated-files.md
- 병합 대상 원본 ref: TARGET_UPSTREAM_REF

목표:
- TARGET_UPSTREAM_REF의 최신 코드 로직, 보안 수정, 버그 수정, API 변경 사항을 가져온다.
- 현재 포크에서 이미 적용된 한국어 번역, 한국어 주석, 한국어 문서, 한국어 프롬프트를 최대한 유지한다.
- 단순 텍스트 보존이 아니라, 최신 upstream 구조 위에 한국어 로컬라이징을 다시 얹는 방식으로 병합한다.

## WeKnora 한글화 병합 프롬프트

아래 문서는 Tencent/WeKnora 원본 업데이트를 현재 한국어 포크에 재병합할 때 그대로 재사용할 수 있도록, 실제 병합 과정에서 확인한 규칙과 함정을 반영해 정리한 실행용 프롬프트다. `TARGET_UPSTREAM_REF`만 원하는 upstream ref로 바꿔서 사용하면 된다.

## 실전 메모

- 기준 base는 `e57deb2`, 한글화 커밋은 `b94551c`다.
- 보호 대상 파일 목록의 원본은 [translate-kr/translated-files.md](translate-kr/translated-files.md)다.
- 이번 병합에서 upstream 구조 변경으로 아래 파일은 삭제 유지하는 것이 맞았다.
  - `config/prompt_templates/rewrite_system.yaml`
  - `config/prompt_templates/rewrite_user.yaml`
  - `docs/WeKnora.md`
  - `internal/models/chat/lkeap.go`
  - `internal/models/chat/qwen.go`
- 위 파일들은 과거 한글화 대상이었더라도, upstream가 새 구조로 대체했다면 억지로 복구하지 말고 대체 파일 기준으로 한국어를 재적용해야 한다.
- 병합 후 별도 UI 스윕이 꼭 필요했다. 특히 아래는 하드코딩 중국어가 남기 쉬운 위치였다.
  - `frontend/src/i18n/index.ts`: 기본 locale / fallbackLocale
  - `frontend/src/views/auth/Login.vue`: 언어 선택기 표시명
  - `frontend/src/views/knowledge/components/FAQEntryManager.vue`: 샘플 데이터, CSV/Excel 헤더, import alias
  - `frontend/src/components/document-preview.vue`: 빈 문서 안내 문구
  - `frontend/src/views/settings/ParserEngineSettings.vue`: `t()` 기본 fallback 문구
- FAQ CSV/Excel 예제 헤더를 한국어로 바꾸면 import 쪽도 한국어 헤더를 읽도록 같이 수정해야 한다. 표시 문자열만 바꾸면 다시 import가 깨진다.
- 한국어 포크에서는 frontend i18n 기본 언어를 `ko-KR`로 두는 것이 안전하다. 그렇지 않으면 신규 키 누락 시 중국어 fallback이 바로 노출된다.

## 복사용 프롬프트

```text
당신은 Tencent/WeKnora 원본 저장소의 최신 업데이트를 현재 한글화 포크에 병합하는 AI 코딩 에이전트다.

배경 정보:
- 원본 기준 커밋: e57deb2
- 한글화 커밋: b94551c
- 한글화 파일 인벤토리 문서: translate-kr/translated-files.md
- 병합 대상 원본 ref: TARGET_UPSTREAM_REF

목표:
- TARGET_UPSTREAM_REF의 최신 코드 로직, 보안 수정, 버그 수정, API 변경 사항을 가져온다.
- 현재 포크에서 이미 적용된 한국어 번역, 한국어 주석, 한국어 문서, 한국어 프롬프트를 최대한 유지한다.
- 단순 텍스트 보존이 아니라, 최신 upstream 구조 위에 한국어 로컬라이징을 다시 얹는 방식으로 병합한다.

필수 작업 순서:
1. 새 통합 브랜치에서 작업한다.
2. 먼저 upstream 변경을 병합하고, unmerged 상태를 모두 해소한다.
3. 그 다음 보호 대상 파일에만 한국어를 재적용한다.
4. 마지막으로 프런트 전체에서 신규 중국어 UI 문자열을 다시 스윕한다.
5. 검증을 실행하고, 실패가 환경 의존인지 코드 문제인지 구분해 보고한다.

병합 원칙:
1. 변경 범위를 먼저 비교한다.
   - base: e57deb2
   - upstream delta: e57deb2..TARGET_UPSTREAM_REF
   - localized delta: e57deb2..b94551c
2. translate-kr/translated-files.md에 있는 파일만 기본 보호 대상으로 취급한다.
3. 각 보호 대상 파일은 3-way 비교로 판단한다.
   - base: e57deb2 시점 파일
   - upstream: TARGET_UPSTREAM_REF 시점 파일
   - localized: 현재 포크 시점 파일
4. 우선순위는 다음과 같다.
   - 코드 로직, 보안 수정, 버그 수정, API/스키마 변경: upstream 우선
   - 사용자 노출 문자열, 주석, 문서, 예시 문구, 프롬프트 텍스트: localized 우선
   - 충돌 시 upstream 구조를 유지하면서 localized의 한국어 표현을 재적용한다.
5. 아래 영역은 특별 취급한다.
   - config/prompt_templates/*: YAML 키, 템플릿 변수, 플레이스홀더, 들여쓰기를 절대 깨뜨리지 말고 값 텍스트의 한국어만 유지한다.
   - config/config.yaml, .env.example, docker-compose*.yml, scripts/*: 명령, 키 이름, 환경변수명, 포트, 경로는 upstream 우선. 설명 문구와 주석만 한국어 유지.
   - frontend/package-lock.json: package.json 또는 실제 의존성 트리 변경이 있을 때만 갱신한다. 불필요한 lockfile 재생성은 피한다.
   - frontend/src/i18n/index.ts: 한국어 포크에서는 기본 locale 과 fallbackLocale 을 ko-KR 기준으로 검토한다.
6. 과거 보호 대상 파일이라도 upstream가 구조 통합으로 제거했다면 억지로 되살리지 않는다.
   - 특히 rewrite_system / rewrite_user 분리 템플릿이 rewrite.yaml로 통합되었는지 확인한다.
   - 모델 provider별 개별 파일이 다른 공용 스펙 파일로 흡수되었는지 확인한다.
7. 번역 유지 때문에 아래 요소가 깨지면 안 된다.
   - 환경변수명
   - YAML 키
   - JSON 키
   - 코드 식별자
   - import 경로
   - format string placeholder
   - prompt template 변수
   - CLI flag 및 명령어

병합 후 추가 UI 스윕 규칙:
1. 하드코딩 중국어가 보이는 Vue/TS 파일을 다시 검색한다.
2. 특히 아래를 우선 확인한다.
   - 언어 선택 드롭다운의 zh-CN 라벨
   - FAQ 샘플 데이터와 CSV/Excel 예제 헤더
   - 문서 미리보기의 빈 상태 문구
   - t('...', '중국어 fallback') 형태의 기본 문자열
3. `zh-CN.ts` 대비 `ko-KR.ts` 키 diff를 돌려 신규 locale 키 누락을 먼저 확인한다.
4. 특히 `datasource`, `auth.oidc*`, `agent.editor.*`, `agentEditor.desc.*`, `tenant.api.desktop*` 계열 신규 키를 우선 확인한다.
5. FAQ 예제 헤더를 바꿨다면 import 쪽 alias도 같이 맞춘다.
6. 한국어 locale 파일에 zhCN 표시명이 남아 있으면 `간체 중국어`로 바꾼다.

검증 규칙:
1. 프런트 검증
   - `pushd frontend >/dev/null && npm install && npm run type-check && npm run build-with-types && popd >/dev/null`
2. Go 검증
   - `pushd . >/dev/null && go test ./... && popd >/dev/null`
3. 검증 실패 시 아래처럼 분류해서 보고한다.
   - 코드/타입 오류
   - 외부 서비스 의존 실패
   - 시스템 라이브러리 누락 예: sqlite3.h
   - 기존 테스트 불안정 또는 mock 불일치

원하는 출력 형식:
- 실제로 수정한 파일 목록
- 파일별 병합 판단 근거 1~2문장
- 삭제 유지한 보호 대상 파일과 그 이유
- 추가로 한국어화한 UI 문자열 위치
- 실행한 검증 명령과 결과 요약
- 자동 해결하지 못한 충돌 또는 사람이 확인해야 할 번역 모호점

주의:
- 목록에 없는 파일까지 무조건 한국어화 대상으로 확대 해석하지 마라.
- 최신 upstream 로직을 놓치지 말고, 한국어 유지 때문에 기능 회귀를 만들지 마라.
- 문서와 주석은 한국어를 유지하되, 코드 구조와 동작은 최신 upstream 기준으로 정렬하라.