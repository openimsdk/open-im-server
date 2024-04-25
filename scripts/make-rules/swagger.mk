# ==============================================================================
# Makefile helper functions for swagger
#

## swagger.run: Generate swagger document.
.PHONY: swagger.run
swagger.run: tools.verify.swagger
	@echo "===========> Generating swagger API docs"
	@$(TOOLS_DIR)/swagger generate spec --scan-models -w $(ROOT_DIR)/cmd/genswaggertypedocs -o $(ROOT_DIR)/api/swagger/swagger.yaml

## swagger.serve: Serve swagger spec and docs.
.PHONY: swagger.serve
swagger.serve: tools.verify.swagger
	@$(TOOLS_DIR)/swagger serve -F=redoc --no-open --port 36666 $(ROOT_DIR)/api/swagger/swagger.yaml

## swagger.help: Display help information about the release package
.PHONY: swagger.help
swagger.help: scripts/make-rules/swagger.mk
	$(call smallhelp)