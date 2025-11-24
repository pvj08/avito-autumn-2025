OAPI     = oapi-codegen
OPENAPI  = api/openapi.yml
CFG_DIR  = api/oapi_configs

.PHONY: oapi_models oapi_server oapi_client generate \
        compose-up compose-down compose-rebuild compose-logs help

# ==== OpenAPI codegen ====

oapi_models:
	$(OAPI) -config $(CFG_DIR)/models.yaml $(OPENAPI)

oapi_server:
	$(OAPI) -config $(CFG_DIR)/server.yaml $(OPENAPI)

oapi_client:
	$(OAPI) -config $(CFG_DIR)/client.yaml $(OPENAPI)

generate: oapi_models oapi_server oapi_client

# ==== Docker Compose ====

compose-up:
	docker compose up -d

compose-down:
	docker compose down

compose-rebuild: generate
	docker compose down
	docker compose build --no-cache
	docker compose up -d

compose-logs:
	docker compose logs -f

# ==== Help ====

help:
	@echo "Available commands:"
	@echo "  make generate        - generate OpenAPI models/server/client"
	@echo "  make compose-up      - start services via docker compose"
	@echo "  make compose-down    - stop services"
	@echo "  make compose-rebuild - regen OpenAPI, rebuild and restart containers"
	@echo "  make compose-logs    - show docker compose logs"
