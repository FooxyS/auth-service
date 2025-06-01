build:
	CGO_enable=0 GOOS=linux GOARCH=amd64 go build -o build/app cmd/main.go
test:
	go test ./internal/usecase/