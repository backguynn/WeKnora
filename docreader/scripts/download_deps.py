#!/usr/bin/env python
# -*- coding: utf-8 -*-

import sys
import os
import logging
from paddleocr import PaddleOCR

# 현재 디렉터리를 Python 경로에 추가
current_dir = os.path.dirname(os.path.abspath(__file__))
if current_dir not in sys.path:
    sys.path.append(current_dir)

# ImageParser 임포트
from parser.image_parser import ImageParser

# 로그 설정
logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s - %(name)s - %(levelname)s - %(message)s",
    datefmt="%Y-%m-%d %H:%M:%S",
)
logger = logging.getLogger(__name__)


def init_ocr_model():
    """Initialize PaddleOCR model to pre-download and cache models"""
    try:
        logger.info("Initializing PaddleOCR model for pre-download...")
        
        # 코드와 동일한 설정 사용
        ocr_config = {
            "use_gpu": False,
            "text_det_limit_type": "max",
            "text_det_limit_side_len": 960,
            "use_doc_orientation_classify": True,  # 문서 방향 분류 활성화
            "use_doc_unwarping": False,
            "use_textline_orientation": True,  # 텍스트 라인 방향 감지 활성화
            "text_recognition_model_name": "PP-OCRv4_server_rec",
            "text_detection_model_name": "PP-OCRv4_server_det",
            "text_det_thresh": 0.3,
            "text_det_box_thresh": 0.6,
            "text_det_unclip_ratio": 1.5,
            "text_rec_score_thresh": 0.0,
            "ocr_version": "PP-OCRv4",
            "lang": "ch",
            "show_log": False,
            "use_dilation": True,
            "det_db_score_mode": "slow",
        }
        
        # PaddleOCR 초기화(모델 다운로드 및 캐시 트리거)
        ocr = PaddleOCR(**ocr_config)
        logger.info("PaddleOCR model initialization completed successfully")
        
        # OCR 기능 테스트로 모델 정상 동작 확인
        import numpy as np
        from PIL import Image
        
        # 간단한 테스트 이미지 생성
        test_image = np.ones((100, 300, 3), dtype=np.uint8) * 255
        test_pil = Image.fromarray(test_image)
        
        # OCR 테스트 1회 실행
        result = ocr.ocr(np.array(test_pil), cls=False)
        logger.info("PaddleOCR test completed successfully")
        
    except Exception as e:
        logger.error(f"Failed to initialize PaddleOCR model: {str(e)}")
        raise
