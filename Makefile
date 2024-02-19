COMPOSE_FILE := ./developments/docker-compose.yml

postgres:
	docker run --name postgres2 -p 5431:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=6677028a -d postgres

createdb:
	docker exec -it postgres2 createdb --username=root --owner=root shop_vui

start-db:
	docker start postgres2

dropdb:
	docker exec -it postgres2 dropdb shop_vui

migrate-up:
	migrate -path db/migration -database "postgresql://root:6677028a@localhost:5431/shop_vui?sslmode=disable" -verbose up

migrate-up1:
	migrate -path db/migration -database "postgresql://root:6677028a@localhost:5431/shop_vui?sslmode=disable" -verbose up

migrate-down:
	migrate -path db/migration -database "postgresql://root:6677028a@localhost:5431/shop_vui?sslmode=disable" -verbose down

migrate-down1:
	migrate -path db/migration -database "postgresql://root:6677028a@localhost:5431/shop_vui?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate

server:
	go run ./cmd/.

proto:
	protoc --proto_path=proto --go_out=pkg/pb --go_opt=paths=source_relative \
    --go-grpc_out=pkg/pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pkg/pb --grpc-gateway_opt=paths=source_relative \
    proto/*.proto

gen-product-proto:
	protoc --proto_path=proto/product --go_out=pkg/pb/product --go_opt=paths=source_relative \
    --go-grpc_out=pkg/pb/product --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pkg/pb/product --grpc-gateway_opt=paths=source_relative \
    proto/product/*.proto
evans:
	evans --host localhost --port 9091 -r repl


gen-proto:
	docker-compose -f ${COMPOSE_FILE} up generate_pb_go

start-postgres:
	docker compose -f ${COMPOSE_FILE} up postgres -d

adminer:
	docker compose -f ${COMPOSE_FILE} up adminer -d

gen-sql:
	docker-compose -f ${COMPOSE_FILE} up generate_sqlc

migrate:
	docker-compose -f ${COMPOSE_FILE} up migrate

user-migrate:
	docker-compose -f ${COMPOSE_FILE} up user_migrate

start-gateway:
	go run ./main.go gateway

start-user:
	go run ./main.go user

gen-mock-user:
	go run github.com/vektra/mockery/v2 \
		--dir=internal/user/repositories \
		--name=UserRepository \
		--filename=user_repository_interface.go \
		--output=internal/user/mocks \
		--outpkg=repomocks