lint-install:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.8

lint:
	golangci-lint run

gotenberg-run:
	docker run \
		-p 3000:3000 \
		--name gotenberg \
		--add-host="host.docker.internal:host-gateway" \
		gotenberg/gotenberg:8

gotenberg:
	docker start gotenberg

# MinIO targets
minio-run:
	docker run -d \
		-p 9000:9000 \
		-p 9001:9001 \
		--name minio \
		-e "MINIO_ROOT_USER=minioadmin" \
		-e "MINIO_ROOT_PASSWORD=minioadmin" \
		minio/minio server /data --console-address ":9001"

minio:
	docker start minio

minio-stop:
	docker stop minio

# API Server targets
api-deps:
	go mod tidy
	go mod download

api-run: api-deps
	cd examples && go run minio_api_server.go

api-test:
	./examples/test_api.sh

# Combined targets
dev: minio api-run

clean:
	docker stop minio gotenberg || true
	docker rm minio gotenberg || true

help:
	@echo "Available targets:"
	@echo "  lint-install  - Install golangci-lint"
	@echo "  lint          - Run linter"
	@echo "  gotenberg-run - Run Gotenberg in Docker"
	@echo "  gotenberg     - Start Gotenberg container"
	@echo "  minio-run     - Run MinIO in Docker"
	@echo "  minio         - Start MinIO container"
	@echo "  minio-stop    - Stop MinIO container"
	@echo "  api-deps      - Install API dependencies"
	@echo "  api-run       - Run API server"
	@echo "  api-test      - Run API tests"
	@echo "  dev           - Start MinIO and API server"
	@echo "  clean         - Stop and remove containers"
