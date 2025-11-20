OAPI       = oapi-codegen
OPENAPI    = api/openapi.yml
CFG_DIR    = api/oapi_configs

.PHONY: oapi_models oapi_server oapi_client generate

oapi_models:
	$(OAPI) -config $(CFG_DIR)/models.yaml $(OPENAPI)

oapi_server:
	$(OAPI) -config $(CFG_DIR)/server.yaml $(OPENAPI)

oapi_client:
	$(OAPI) -config $(CFG_DIR)/client.yaml $(OPENAPI)

generate: oapi_models oapi_server oapi_client