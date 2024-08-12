# Makefile extensions for windows-arm64.

# -----------------------------------------------------------------------------
# OS-ARCH specific targets
# -----------------------------------------------------------------------------

.PHONY: build-osarch-specific
build-osarch-specific: windows/arm64
	@mv $(TARGET_DIRECTORY)/windows-arm64/$(PROGRAM_NAME) $(TARGET_DIRECTORY)/windows-arm64/$(PROGRAM_NAME).exe

# -----------------------------------------------------------------------------
# Makefile targets supported only by this platform.
# -----------------------------------------------------------------------------

.PHONY: only-windows-arm64
only-windows-arm64:
	$(info Only windows-arm64 has this Makefile target.)
