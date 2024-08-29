# Makefile extensions for linux-arm64.

# -----------------------------------------------------------------------------
# OS-ARCH specific targets
# -----------------------------------------------------------------------------

.PHONY: build-osarch-specific
build-osarch-specific: linux/arm64
	@mkdir -p $(TARGET_DIRECTORY)/linux
	@cp $(TARGET_DIRECTORY)/linux-arm64/$(PROGRAM_NAME) $(TARGET_DIRECTORY)/linux/$(PROGRAM_NAME)


.PHONY: docker-build-osarch-specific
docker-build-osarch-specific:
	@docker build \
		--platform linux/amd64 \
		--tag $(DOCKER_IMAGE_NAME) \
		--tag $(DOCKER_IMAGE_NAME):$(BUILD_VERSION) \
		.

# -----------------------------------------------------------------------------
# Makefile targets supported only by this platform.
# -----------------------------------------------------------------------------

.PHONY: only-linux-arm64
only-linux-arm64:
	$(info Only linux-arm64 has this Makefile target.)
