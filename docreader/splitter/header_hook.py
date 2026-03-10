import re
from typing import Callable, Dict, List, Match, Pattern, Union

from pydantic import BaseModel, Field


class HeaderTrackerHook(BaseModel):
    """표 헤더 추적 Hook 설정 클래스, 다양한 표 헤더 인식 시나리오 지원"""

    start_pattern: Pattern[str] = Field(
        description="표 헤더 시작 매칭(정규식 또는 문자열)"
    )
    end_pattern: Pattern[str] = Field(description="표 헤더 종료 매칭(정규식 또는 문자열)")
    extract_header_fn: Callable[[Match[str]], str] = Field(
        default=lambda m: m.group(0),
        description="시작 매칭 결과에서 헤더 내용을 추출하는 함수(기본: 매칭된 전체 내용)",
    )
    priority: int = Field(default=0, description="우선순위(복수 설정 시 높은 우선순위 먼저 매칭)")
    case_sensitive: bool = Field(
        default=True, description="대소문자 구분 여부(문자열 pattern 입력 시에만 적용)"
    )

    def __init__(
        self,
        start_pattern: Union[str, Pattern[str]],
        end_pattern: Union[str, Pattern[str]],
        **kwargs,
    ):
        flags = 0 if kwargs.get("case_sensitive", True) else re.IGNORECASE
        if isinstance(start_pattern, str):
            start_pattern = re.compile(start_pattern, flags | re.DOTALL)
        if isinstance(end_pattern, str):
            end_pattern = re.compile(end_pattern, flags | re.DOTALL)
        super().__init__(
            start_pattern=start_pattern,
            end_pattern=end_pattern,
            **kwargs,
        )


# 표 헤더 Hook 설정 초기화(기본 설정 제공: Markdown 표, 코드 블록 지원)
DEFAULT_CONFIGS = [
    # 코드 블록 설정(```로 시작, ```로 종료)
    # HeaderTrackerHook(
    #     # 코드 블록 시작(언어 지정 지원)
    #     start_pattern=r"^\s*```(\w+).*(?!```$)",
    #     # 코드 블록 종료
    #     end_pattern=r"^\s*```.*$",
    #     extract_header_fn=lambda m: f"```{m.group(1)}" if m.group(1) else "```",
    #     priority=20,  # 코드 블록 우선순위가 표보다 높음
    #     case_sensitive=True,
    # ),
    # Markdown 표 설정(헤더에 밑줄 포함)
    HeaderTrackerHook(
        # 헤더 행 + 구분 행
        start_pattern=r"^\s*(?:\|[^|\n]*)+[\r\n]+\s*(?:\|\s*:?-{3,}:?\s*)+\|?[\r\n]+$",
        # 빈 줄 또는 비표 콘텐츠
        end_pattern=r"^\s*$|^\s*[^|\s].*$",
        priority=15,
        case_sensitive=False,
    ),
]
DEFAULT_CONFIGS.sort(key=lambda x: -x.priority)


# 定义Hook状态数据结构
class HeaderTracker(BaseModel):
    """표 헤더 추적 Hook 상태 클래스"""

    header_hook_configs: List[HeaderTrackerHook] = Field(default=DEFAULT_CONFIGS)
    active_headers: Dict[int, str] = Field(default_factory=dict)
    ended_headers: set[int] = Field(default_factory=set)

    def update(self, split: str) -> Dict[int, str]:
        """현재 split에서 헤더 시작/종료를 감지하고 Hook 상태를 업데이트"""
        new_headers: Dict[int, str] = {}

        # 1. 헤더 종료 마커 확인
        for config in self.header_hook_configs:
            if config.priority in self.active_headers and config.end_pattern.search(
                split
            ):
                self.ended_headers.add(config.priority)
                del self.active_headers[config.priority]

        # 2. 새로운 헤더 시작 마커 확인(활성/종료되지 않은 항목만)
        for config in self.header_hook_configs:
            if (
                config.priority not in self.active_headers
                and config.priority not in self.ended_headers
            ):
                match = config.start_pattern.search(split)
                if match:
                    header = config.extract_header_fn(match)
                    self.active_headers[config.priority] = header
                    new_headers[config.priority] = header

        # 3. 모든 활성 헤더가 종료되었는지 확인(종료 마커 초기화)
        if not self.active_headers:
            self.ended_headers.clear()

        return new_headers

    def get_headers(self) -> str:
        """현재 활성 헤더를 우선순위 순으로 결합한 텍스트 반환"""
        # 우선순위 내림차순으로 헤더 정렬
        sorted_headers = sorted(self.active_headers.items(), key=lambda x: -x[0])
        return (
            "\n".join([header for _, header in sorted_headers])
            if sorted_headers
            else ""
        )
