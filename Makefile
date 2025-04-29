PROTO_DIR = ./api/
OUT_DIR = ./pkg/api/

PROTO_FILES = $(shell find $(PROTO_DIR) -name "*.proto")

gen:
	@echo "Generating gRPC files..."
	@mkdir -p $(OUT_DIR)
	protoc -I$(PROTO_DIR) -I ./third_party/googleapis \
		--go_out=$(OUT_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(OUT_DIR) --go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=$(OUT_DIR) --grpc-gateway_opt=paths=source_relative \
		--openapiv2_out=$(OUT_DIR) --openapiv2_opt=logtostderr=true \
		$(PROTO_FILES)

build:
	docker-compose --env-file ./backend/env/users_postgres.env -f ./backend/docker/docker-compose.yml up -d --build postgresql_users
	docker-compose --env-file ./backend/env/user.env -f ./backend/docker/docker-compose.yml up -d --build user_service