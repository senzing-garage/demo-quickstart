# Makefile extensions for windows.

# -----------------------------------------------------------------------------
# OS-ARCH specific targets
# -----------------------------------------------------------------------------

.PHONY: build-osarch-specific
build-osarch-specific: windows/amd64
	@mv $(TARGET_DIRECTORY)/windows-amd64/$(PROGRAM_NAME) $(TARGET_DIRECTORY)/windows-amd64/$(PROGRAM_NAME).exe

# -----------------------------------------------------------------------------
# Makefile targets supported only by this platform.
# -----------------------------------------------------------------------------

.PHONY: only-windows-x86_64
only-windows-x86_64:
	@echo "Only windows-x86_64 has this Makefile target."
