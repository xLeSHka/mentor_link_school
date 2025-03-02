

run_tests:
	docker-compose -f docker-compose.test.yml up
	go test  ./tests/e2e_test.go -timeout 120s  -v -coverpkg='./internal/...' -coverprofile cover.out.tmp && cat cover.out.tmp > cover.out && rm cover.out.tmp && go tool cover -func cover.out
	docker-compose -f docker-compose.yml down -v

