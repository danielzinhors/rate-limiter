gerar_mocks:
	mockgen -source=./ratelimiter/adapters/storage.go -destination ./ratelimiter/mocks/storage.go -package mocks
	mockgen -source=./ratelimiter/response_writer/response.go -destination ./ratelimiter/mocks/response.go -package mocks

testar:
	go test ./... -v

