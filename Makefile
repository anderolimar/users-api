check-swag:
	which swag || go install github.com/swaggo/swag/cmd/swag

check-mock:
	which mockgen || go install github.com/golang/mock/mockgen@v1.6.0

swag: check-swag
	go mod tidy && swag init

mock: 
	mockgen -source ./users/repository.go -destination ./users/mock_repository.go -package users
	mockgen -source ./users/service.go -destination ./users/mock_service.go -package users
envup: 
	docker-compose up -d

envdown: 
	docker-compose down

test: 
	go test ./...

attack:
	vegeta attack -duration=10s -rate=100 -targets=target.conf | vegeta report	

