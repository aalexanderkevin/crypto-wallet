grpc:
	@echo "Generate pb.go file from proto"
	@sh proto.sh

migrate:
	@echo ">> Running DB migration"
	@env $$(cat .env | xargs) go run github.com/aalexanderkevin/crypto-wallet/cmd/wallet migrate

test-all: test-unit test-integration-with-infra

test-unit:
	@echo ">> Running Unit Test"
	@env $$(cat .env.testing | xargs) go test -tags=unit -failfast -cover -covermode=atomic ./...

test-integration:
	@echo ">> Running Integration Test"
	@env $$(cat .env.testing | xargs) env POSTGRES_MIGRATION_PATH=$$(pwd)/migrations go test -tags=integration -failfast -cover -covermode=atomic ./...

test-integration-with-infra: test-infra-up test-integration test-infra-down

test-infra-up:
	$(MAKE) test-infra-down
	@echo ">> Starting Test DB"
	docker run -d --rm --name test-postgres -p 5431:5432 --env-file .env.testing postgres:12
	docker ps

test-infra-down:
	@echo ">> Shutting Down Test DB"
	@-docker kill test-postgres