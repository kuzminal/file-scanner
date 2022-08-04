test: test_scanner test_storage

test_scanner:
	@cd ./internal/scanner && go test . -v

test_storage:
	@cd ./internal/storage && go test . -v

bench:
	@cd ./internal/utils && go test -bench . -benchmem -benchtime=10s

generate:
	@cd ./internal/wrappers && go generate

run:
	@go run ./cmd/app/