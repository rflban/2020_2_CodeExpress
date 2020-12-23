.PHONY:

build:
	go build -o main_service app/main/main.go
	go build -o track_service app/track_microservice/main.go
	go build -o session_service app/session_microservice/main.go  

tests:
	go test -coverpkg=./... -coverprofile cover.out.tmp ./...
	cat cover.out.tmp | grep -v "mock_*" | grep -v ".pb.go" | grep -v ".pb" | grep -v "_easyjson.go"> cover.out
	go tool cover -func cover.out

linter:
	golangci-lint run ./...
