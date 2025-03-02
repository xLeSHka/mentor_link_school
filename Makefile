

run:
	docker run --name postgres_test -e POSTGRES_PASSWORD=123 -e POSTGRES_USER=root -e POSTGRES_DATABASE=postgres -d -p 5432:5432 postgres
	docker run --name minio_test -d -e MINIO_BUCKET_NAME=testbucket -e MINIO_ROOT_USER=testtest -e MINIO_ROOT_PASSWORD=testtest -p 9000:9000 -p 9001:9001 minio/minio server /data --console-address ":9001"
	docker run -d --name redis_test -p 6379:6379 redis
	go test  ./tests/e2e_test.go -timeout 120s  -v -coverpkg='./internal/...' -coverprofile cover.out.tmp && cat cover.out.tmp > cover.out && rm cover.out.tmp && go tool cover -func cover.out
stop:
	docker stop postgres_test
	docker stop minio_test
	docker stop redis_test
	docker rm postgres_test
	docker rm minio_test
	docker rm redis_test
	docker volume prune
