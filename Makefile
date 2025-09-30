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