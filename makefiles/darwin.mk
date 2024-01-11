# Makefile extensions for darwin.

# -----------------------------------------------------------------------------
# Variables
# -----------------------------------------------------------------------------

SENZING_DIR ?= /opt/senzing/g2
SENZING_TOOLS_SENZING_DIRECTORY ?= $(SENZING_DIR)

LD_LIBRARY_PATH := $(SENZING_TOOLS_SENZING_DIRECTORY)/lib:$(SENZING_TOOLS_SENZING_DIRECTORY)/lib/macos
DYLD_LIBRARY_PATH := $(LD_LIBRARY_PATH)
SENZING_TOOLS_DATABASE_URL ?= sqlite3://na:na@/tmp/sqlite/G2C.db

# -----------------------------------------------------------------------------
# OS specific targets
# -----------------------------------------------------------------------------

.PHONY: build-osarch-specific
build-osarch-specific: darwin/amd64


.PHONY: clean-osarch-specific
clean-osarch-specific:
	@docker rm --force $(DOCKER_CONTAINER_NAME) 2> /dev/null || true
	@docker rmi --force $(DOCKER_IMAGE_NAME) $(DOCKER_BUILD_IMAGE_NAME) 2> /dev/null || true
	@rm -rf $(TARGET_DIRECTORY) || true
	@rm -f $(GOPATH)/bin/$(PROGRAM_NAME) || true
	@rm -rf /tmp/sqlite || true


.PHONY: hello-world-osarch-specific
hello-world-osarch-specific:
	@echo "Hello World, from darwin."


.PHONY: run-osarch-specific
run-osarch-specific:
	@go run -exec macos_exec_dyld.sh main.go


.PHONY: setup-osarch-specific
setup-osarch-specific:
	@mkdir /tmp/sqlite
	@cp testdata/sqlite/G2C.db /tmp/sqlite/G2C.db
	@mkdir -p $(TARGET_DIRECTORY)/$(GO_OS)-$(GO_ARCH) || true


.PHONY: test-osarch-specific
test-osarch-specific:
	@go test -exec macos_exec_dyld.sh -v -p 1 ./...

# -----------------------------------------------------------------------------
# Makefile targets supported only by this platform.
# -----------------------------------------------------------------------------

.PHONY: only-darwin
only-darwin:
	@echo "Only darwin has this Makefile target."
