# Makefile extensions for darwin-arm64.

# -----------------------------------------------------------------------------
# OS-ARCH specific targets
# -----------------------------------------------------------------------------

.PHONY: build-osarch-specific
build-osarch-specific: darwin/arm64


.PHONY: docker-build-osarch-specific
docker-build-osarch-specific:
	@docker build \
		--platform linux/arm64 \
		--tag $(DOCKER_IMAGE_NAME) \
		--tag $(DOCKER_IMAGE_NAME):$(BUILD_VERSION) \
		.

# -----------------------------------------------------------------------------
# Makefile targets supported only by this platform.
# -----------------------------------------------------------------------------

.PHONY: only-darwin-arm64
only-darwin-arm64:
	$(info Only darwin-arm64 has this Makefile target.)
