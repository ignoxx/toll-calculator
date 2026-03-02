SERVICES := $(shell find cmd/* -type d -exec basename {} \;)

.PHONY: $(SERVICES)
$(SERVICES): proto
	@go run "./cmd/$@"

run-kafka:
	@docker-compose up

proto:
	@protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative types/ptypes.proto
