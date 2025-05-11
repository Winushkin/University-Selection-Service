PROTO_DIR = ./backend/api/
OUT_DIR = ./backend/pkg/api/

PROTO_FILES = $(shell find $(PROTO_DIR) -name "*.proto")

gen:
	@echo "Generating gRPC files..."
	@mkdir -p $(OUT_DIR)
	protoc -I$(PROTO_DIR) -I ./backend/third_party/googleapis \
		--go_out=$(OUT_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(OUT_DIR) --go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=$(OUT_DIR) --grpc-gateway_opt=paths=source_relative \
		--openapiv2_out=$(OUT_DIR) --openapiv2_opt=logtostderr=true \
		--openapiv2_opt logtostderr=true 										\
            --openapiv2_opt generate_unbound_methods=true 					\
            --openapiv2_opt allow_merge=true 										\
            --openapiv2_opt merge_file_name=allservices 						\
		$(PROTO_FILES)

build:
	docker compose --env-file ./backend/env/users_postgres.env -f ./backend/docker/docker-compose.yml up -d --build postgresql_users
	docker compose --env-file ./backend/env/universities.env -f ./backend/docker/docker-compose.yml up -d --build postgresql_universities
	docker compose --env-file ./backend/env/user.env -f ./backend/docker/docker-compose.yml up -d --build user_service
	docker compose --env-file ./backend/env/analytic.env -f ./backend/docker/docker-compose.yml up -d --build analytic_service
	docker compose -f ./backend/docker/docker-compose.yml up -d --build gateway_service
	docker compose -f ./backend/docker/docker-compose.yml up -d --build nginx_service

down:
	migrate -database ${POSTGRESQL_URL} -path backend/db/migrations/universities down


up:
	migrate -database ${POSTGRESQL_URL} -path backend/db/migrations/universities up