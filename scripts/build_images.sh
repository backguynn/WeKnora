#!/bin/bash
# 이 스크립트는 소스에서 WeKnora의 모든 Docker 이미지를 빌드합니다

# 색상 설정
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # 색상 없음

# 프로젝트 루트 디렉터리(스크립트 디렉터리의 상위)
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJECT_ROOT="$( cd "$SCRIPT_DIR/.." && pwd )"

# 버전 정보
VERSION="1.0.0"
SCRIPT_NAME=$(basename "$0")

# 도움말 표시
show_help() {
    echo -e "${GREEN}WeKnora 이미지 빌드 스크립트 v${VERSION}${NC}"
    echo -e "${GREEN}사용법:${NC} $0 [옵션]"
    echo "옵션:"
    echo "  -h, --help     도움말 표시"
    echo "  -a, --all      모든 이미지 빌드(기본값)"
    echo "  -p, --app      앱 이미지 빌드만"
    echo "  -d, --docreader 문서 리더 이미지 빌드만"
    echo "  -f, --frontend 프론트엔드 이미지 빌드만"
    echo "  -s, --sandbox  샌드박스 이미지 빌드만"
    echo "  -c, --clean    로컬 이미지 전체 정리"
    echo "  -v, --version  버전 정보 표시"
    exit 0
}

# 버전 정보 표시
show_version() {
    echo -e "${GREEN}WeKnora 이미지 빌드 스크립트 v${VERSION}${NC}"
    exit 0
}

# 로그 함수
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

# Docker 설치 여부 확인
check_docker() {
    log_info "Docker 환경을 확인 중..."
    
    if ! command -v docker &> /dev/null; then
        log_error "Docker가 설치되어 있지 않습니다. 먼저 Docker를 설치하세요"
        return 1
    fi
    
    # Docker 서비스 실행 상태 확인
    if ! docker info &> /dev/null; then
        log_error "Docker 서비스가 실행 중이 아닙니다. Docker 서비스를 시작하세요"
        return 1
    fi
    
    log_success "Docker 환경 점검 통과"
    return 0
}

# 플랫폼 확인
check_platform() {
    log_info "시스템 플랫폼 정보를 확인 중..."
    if [ "$(uname -m)" = "x86_64" ]; then
        export PLATFORM="linux/amd64"
        export TARGETARCH="amd64"
    elif [ "$(uname -m)" = "aarch64" ] || [ "$(uname -m)" = "arm64" ]; then
        export PLATFORM="linux/arm64"
        export TARGETARCH="arm64"
    else
        log_warning "알 수 없는 플랫폼 유형입니다: $(uname -m). 기본 플랫폼 linux/amd64를 사용합니다"
        export PLATFORM="linux/amd64"
        export TARGETARCH="amd64"
    fi
    log_info "현재 플랫폼: $PLATFORM"
    log_info "현재 아키텍처: $TARGETARCH"
}

# 버전 정보 가져오기
get_version_info() {
    # VERSION 파일에서 버전 가져오기
    if [ -f "VERSION" ]; then
        VERSION=$(cat VERSION | tr -d '\n\r')
    else
        VERSION="unknown"
    fi
    
    # 커밋 ID 가져오기
    if command -v git >/dev/null 2>&1; then
        COMMIT_ID=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
    else
        COMMIT_ID="unknown"
    fi
    
    # 빌드 시간 가져오기
    BUILD_TIME=$(date -u '+%Y-%m-%d %H:%M:%S UTC')
    
    # Go 버전 가져오기
    if command -v go >/dev/null 2>&1; then
        GO_VERSION=$(go version 2>/dev/null || echo "unknown")
    else
        GO_VERSION="unknown"
    fi
    
    log_info "버전 정보: $VERSION"
    log_info "Commit ID: $COMMIT_ID"
    log_info "빌드 시간: $BUILD_TIME"
    log_info "Go 버전: $GO_VERSION"
}

# 앱 이미지 빌드
build_app_image() {
    log_info "앱 이미지 빌드(weknora-app)..."
    
    cd "$PROJECT_ROOT"
    
    # 버전 정보 가져오기
    get_version_info
    
    docker build \
        --platform $PLATFORM \
        --build-arg GOPRIVATE_ARG=${GOPRIVATE:-""} \
        --build-arg GOPROXY_ARG=${GOPROXY:-"https://goproxy.cn,direct"} \
        --build-arg GOSUMDB_ARG=${GOSUMDB:-"off"} \
        --build-arg VERSION_ARG="$VERSION" \
        --build-arg COMMIT_ID_ARG="$COMMIT_ID" \
        --build-arg BUILD_TIME_ARG="$BUILD_TIME" \
        --build-arg GO_VERSION_ARG="$GO_VERSION" \
        -f docker/Dockerfile.app \
        -t wechatopenai/weknora-app:latest \
        .
    
    if [ $? -eq 0 ]; then
        log_success "앱 이미지 빌드 성공"
        return 0
    else
        log_error "앱 이미지 빌드 실패"
        return 1
    fi
}

# 문서 리더 이미지 빌드
build_docreader_image() {
    log_info "문서 리더 이미지 빌드(weknora-docreader)..."
    
    cd "$PROJECT_ROOT"
    
    docker build \
        --platform $PLATFORM \
        --build-arg PLATFORM=$PLATFORM \
        --build-arg TARGETARCH=$TARGETARCH \
        --build-arg APT_MIRROR=${APT_MIRROR:-} \
        -f docker/Dockerfile.docreader \
        -t wechatopenai/weknora-docreader:latest \
        .
    
    if [ $? -eq 0 ]; then
        log_success "문서 리더 이미지 빌드 성공"
        return 0
    else
        log_error "문서 리더 이미지 빌드 실패"
        return 1
    fi
}

# 프론트엔드 이미지 빌드
build_frontend_image() {
    log_info "프론트엔드 이미지 빌드(weknora-ui)..."
    
    cd "$PROJECT_ROOT"
    
    docker build \
        --platform $PLATFORM \
        -f frontend/Dockerfile \
        -t wechatopenai/weknora-ui:latest \
        frontend/
    
    if [ $? -eq 0 ]; then
        log_success "프론트엔드 이미지 빌드 성공"
        return 0
    else
        log_error "프론트엔드 이미지 빌드 실패"
        return 1
    fi
}

# 샌드박스 이미지 빌드
build_sandbox_image() {
    log_info "샌드박스 이미지 빌드(weknora-sandbox)..."

    cd "$PROJECT_ROOT"

    docker build \
        --platform $PLATFORM \
        -f docker/Dockerfile.sandbox \
        -t wechatopenai/weknora-sandbox:latest \
        .

    if [ $? -eq 0 ]; then
        log_success "샌드박스 이미지 빌드 성공"
        return 0
    else
        log_error "샌드박스 이미지 빌드 실패"
        return 1
    fi
}

# 모든 이미지 빌드
build_all_images() {
    log_info "모든 이미지 빌드를 시작합니다..."

    local app_result=0
    local docreader_result=0
    local frontend_result=0
    local sandbox_result=0

    # 앱 이미지 빌드
    build_app_image
    app_result=$?

    # 문서 리더 이미지 빌드
    build_docreader_image
    docreader_result=$?

    # 프론트엔드 이미지 빌드
    build_frontend_image
    frontend_result=$?

    # 샌드박스 이미지 빌드
    build_sandbox_image
    sandbox_result=$?

    # 빌드 결과 표시
    echo ""
    log_info "=== 빌드 결과 ==="
    if [ $app_result -eq 0 ]; then
        log_success "✓ 앱 이미지 빌드 성공"
    else
        log_error "✗ 앱 이미지 빌드 실패"
    fi

    if [ $docreader_result -eq 0 ]; then
        log_success "✓ 문서 리더 이미지 빌드 성공"
    else
        log_error "✗ 문서 리더 이미지 빌드 실패"
    fi

    if [ $frontend_result -eq 0 ]; then
        log_success "✓ 프론트엔드 이미지 빌드 성공"
    else
        log_error "✗ 프론트엔드 이미지 빌드 실패"
    fi

    if [ $sandbox_result -eq 0 ]; then
        log_success "✓ 샌드박스 이미지 빌드 성공"
    else
        log_error "✗ 샌드박스 이미지 빌드 실패"
    fi

    if [ $app_result -eq 0 ] && [ $docreader_result -eq 0 ] && [ $frontend_result -eq 0 ] && [ $sandbox_result -eq 0 ]; then
        log_success "모든 이미지 빌드가 완료되었습니다!"
        return 0
    else
        log_error "일부 이미지 빌드가 실패했습니다"
        return 1
    fi
}

# 로컬 이미지 정리
clean_images() {
    log_info "로컬 WeKnora 이미지를 정리합니다..."
    
    # 관련 컨테이너 중지
    log_info "관련 컨테이너를 중지합니다..."
    docker stop $(docker ps -q --filter "ancestor=wechatopenai/weknora-app:latest" 2>/dev/null) 2>/dev/null || true
    docker stop $(docker ps -q --filter "ancestor=wechatopenai/weknora-docreader:latest" 2>/dev/null) 2>/dev/null || true
    docker stop $(docker ps -q --filter "ancestor=wechatopenai/weknora-ui:latest" 2>/dev/null) 2>/dev/null || true
    
    # 관련 컨테이너 삭제
    log_info "관련 컨테이너를 삭제합니다..."
    docker rm $(docker ps -aq --filter "ancestor=wechatopenai/weknora-app:latest" 2>/dev/null) 2>/dev/null || true
    docker rm $(docker ps -aq --filter "ancestor=wechatopenai/weknora-docreader:latest" 2>/dev/null) 2>/dev/null || true
    docker rm $(docker ps -aq --filter "ancestor=wechatopenai/weknora-ui:latest" 2>/dev/null) 2>/dev/null || true
    
    # 이미지 삭제
    log_info "로컬 이미지를 삭제합니다..."
    docker rmi wechatopenai/weknora-app:latest 2>/dev/null || true
    docker rmi wechatopenai/weknora-docreader:latest 2>/dev/null || true
    docker rmi wechatopenai/weknora-ui:latest 2>/dev/null || true
    docker rmi wechatopenai/weknora-sandbox:latest 2>/dev/null || true
    
    docker image prune -f
    
    log_success "이미지 정리가 완료되었습니다"
    return 0
}

# 커맨드라인 인자 파싱
BUILD_ALL=false
BUILD_APP=false
BUILD_DOCREADER=false
BUILD_FRONTEND=false
BUILD_SANDBOX=false
CLEAN_IMAGES=false

# 인자가 없으면 모든 이미지를 기본으로 빌드
if [ $# -eq 0 ]; then
    BUILD_ALL=true
fi

while [ "$1" != "" ]; do
    case $1 in
        -h | --help )       show_help
                            ;;
        -a | --all )        BUILD_ALL=true
                            ;;
        -p | --app )        BUILD_APP=true
                            ;;
        -d | --docreader )  BUILD_DOCREADER=true
                            ;;
        -f | --frontend )   BUILD_FRONTEND=true
                            ;;
        -s | --sandbox )    BUILD_SANDBOX=true
                            ;;
        -c | --clean )      CLEAN_IMAGES=true
                            ;;
        -v | --version )    show_version
                            ;;
        * )                 log_error "알 수 없는 옵션: $1"
                            show_help
                            ;;
    esac
    shift
done

# Docker 환경 확인
check_docker
if [ $? -ne 0 ]; then
    exit 1
fi

# 플랫폼 확인
check_platform

# 정리 작업 실행
if [ "$CLEAN_IMAGES" = true ]; then
    clean_images
    exit $?
fi

# 빌드 작업 실행
if [ "$BUILD_ALL" = true ]; then
    build_all_images
    exit $?
fi

if [ "$BUILD_APP" = true ]; then
    build_app_image
    exit $?
fi

if [ "$BUILD_DOCREADER" = true ]; then
    build_docreader_image
    exit $?
fi

if [ "$BUILD_FRONTEND" = true ]; then
    build_frontend_image
    exit $?
fi

if [ "$BUILD_SANDBOX" = true ]; then
    build_sandbox_image
    exit $?
fi

exit 0