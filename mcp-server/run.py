#!/usr/bin/env python3
"""
WeKnora MCP Server 간편 실행 스크립트

간단한 실행 스크립트로 기본 기능만 제공합니다.
더 많은 옵션은 main.py를 사용하세요.
"""

import os
import sys
from pathlib import Path


def main():
    """간단한 실행 함수"""
    # 현재 디렉터리를 Python 경로에 추가
    current_dir = Path(__file__).parent.absolute()
    if str(current_dir) not in sys.path:
        sys.path.insert(0, str(current_dir))

    # 환경 변수 확인
    base_url = os.getenv("WEKNORA_BASE_URL", "http://localhost:8080/api/v1")
    api_key = os.getenv("WEKNORA_API_KEY", "")

    print("WeKnora MCP Server")
    print(f"Base URL: {base_url}")
    print(f"API Key: {'설정됨' if api_key else '미설정'}")
    print("-" * 40)

    try:
        # 임포트 후 실행
        from main import sync_main

        sync_main()
    except ImportError:
        print("오류: 필요한 모듈을 임포트할 수 없습니다")
        print("다음 명령을 실행했는지 확인하세요: pip install -r requirements.txt")
        sys.exit(1)
    except KeyboardInterrupt:
        print("\n서버가 중지되었습니다")
    except Exception as e:
        print(f"오류: {e}")
        sys.exit(1)


if __name__ == "__main__":
    main()
