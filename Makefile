run:
	go run ./cmd/web
deps:
	go mod tidy
test:
	go test -v ./...
vet:
	go vet ./...