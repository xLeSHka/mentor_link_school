
run_tests:
	docker-compose -f docker-compose.test.yml up -d
	go test  ./tests/... -timeout 120s  -v -coverpkg=./internal/transport,./internal/service/...,./internal/repository/...,./internal/pkg/...,./internal/connections/...
	docker-compose -f docker-compose.test.yml down -v
