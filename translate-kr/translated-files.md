# WeKnora 한글화 파일 목록

## 기준 이력

- 원본 기준 브랜치: upstream/main
- 원본 기준 커밋: e57deb2
- 포크 한글화 커밋: b94551c
- 포크 한글화 커밋 메시지: korean comment
- 포크 한글화 커밋 날짜: 2026-03-10 14:04:41 +0900
- 비교 범위: upstream/main..HEAD

## 추출 기준

이 문서는 git 이력만 기준으로 작성했다. 현재 포크는 원본 upstream/main 대비 한 커밋만 앞서 있으며, 그 한 커밋에서 수정된 파일 48개 전부의 diff에 한글 문자열이 포함되어 있었다. 따라서 아래 목록 전체를 이 포크의 한글화 대상 파일로 본다.

추출에 사용한 핵심 확인 포인트:

- git log upstream/main..HEAD 결과가 b94551c 한 건만 반환됨
- git diff --name-status upstream/main..HEAD 결과로 수정 파일 48개 확인
- 각 파일 diff 안에 실제 한글 문자열이 존재함을 추가 확인

## 요약

- 총 한글화 파일 수: 48
- 신규 추가 파일: 1
- 수정 파일: 47

병합 시 우선적으로 주의해야 하는 영역:

- 문서와 README류
- 프롬프트 템플릿과 설정 파일
- 스크립트 주석 및 안내 문구
- 프론트엔드 엔트리 및 배포 설정 파일
- 일부 Go, Python 코드 내 주석 또는 사용자 노출 문자열

## 특이사항

- README_KE.md가 한글 문서로 신규 추가되었다.
- 저장소에는 별도로 README_KO.md도 이미 존재하지만, 이번 한글화 커밋의 diff 범위에는 포함되지 않았다.
- frontend/package-lock.json도 한글화 커밋에 포함되어 있으므로, 후속 병합 시 실제 의존성 변경인지 잠금 파일 재생성 부수효과인지 구분해서 처리하는 편이 안전하다.

## 한글화 파일 목록

### 루트 및 메타 파일

- .env.example
- .github/ISSUE_TEMPLATE/config.yml
- .github/pull_request_template.md
- .github/workflows/docker-image.yml
- .golangci.yml
- README.md
- README_KE.md
- docker-compose.yml

### 클라이언트 및 설정

- client/faq.go
- config/config.yaml
- config/prompt_templates/context_template.yaml
- config/prompt_templates/fallback.yaml
- config/prompt_templates/rewrite_system.yaml
- config/prompt_templates/rewrite_user.yaml
- config/prompt_templates/system_prompt.yaml

### Docker 및 문서 파서 관련

- docker/Dockerfile.app
- docker/Dockerfile.docreader
- docreader/Makefile
- docreader/scripts/download_deps.py
- docreader/scripts/generate_proto.sh
- docreader/splitter/header_hook.py

### 프로젝트 문서

- docs/BUILTIN_MODELS.md
- docs/KnowledgeGraph.md
- docs/QA.md
- docs/ROADMAP.md
- docs/WeKnora.md
- docs/agent-skills.md
- docs/api/README.md

### 프론트엔드 및 배포 진입 파일

- frontend/docker-entrypoint.sh
- frontend/index.html
- frontend/nginx.conf
- frontend/package-lock.json
- frontend/packages/.gitkeep
- frontend/public/config.js
- frontend/src/main.ts
- frontend/vite.config.ts

### 내부 코드 및 실행 스크립트

- internal/errors/errors.go
- internal/models/chat/lkeap.go
- internal/models/chat/qwen.go
- internal/models/chat/remote_api.go
- internal/types/custom_agent.go
- mcp-server/run.py
- scripts/build_images.sh
- scripts/check-env.sh
- scripts/dev.sh
- scripts/get_version.sh
- scripts/quick-dev.sh
- scripts/start_all.sh

## 후속 병합용 메모

- 이 목록은 현재 포크의 한글화 보호 대상 파일 목록으로 사용할 수 있다.
- 나중에 원본 프로젝트에서 업데이트를 가져올 때는, 목록 안의 파일만 우선적으로 3-way 비교 대상으로 두고 한국어 텍스트를 재적용하는 전략이 효율적이다.
- 바로 사용할 프롬프트 초안은 translate-kr/merge-prompt.md에 정리했다.