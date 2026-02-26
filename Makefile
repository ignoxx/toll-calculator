SERVICES := $(shell find cmd/* -type d -exec basename {} \;)

.PHONY: $(SERVICES)
$(SERVICES):
	@go run "./cmd/$@"

run-kafka:
	@docker-compose up
