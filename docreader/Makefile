.PHONY: proto build run docker-build docker-run clean

# protobuf 코드 생성
proto:
	@echo "Generating protobuf code..."
	@sh ./scripts/generate_proto.sh

# Go 클라이언트 빌드
build:
	@echo "Building Go client..."
	@go build -o bin/client ./src/client

# Python 서비스 실행
run:
	@echo "Running Python server..."
	@python src/server/server.py

# 정리
clean:
	@echo "Cleaning up..."
	@rm -rf bin/
	@find . -name "*.pyc" -delete
	@find . -name "__pycache__" -delete 