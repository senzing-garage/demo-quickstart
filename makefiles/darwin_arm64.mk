# Makefile extensions for darwin-arm64.

# -----------------------------------------------------------------------------
# OS-ARCH specific targets
# -----------------------------------------------------------------------------

.PHONY: build-osarch-specific
build-osarch-specific: darwin/arm64

# -----------------------------------------------------------------------------
# Makefile targets supported only by this platform.
# -----------------------------------------------------------------------------

.PHONY: only-darwin-arm64
only-darwin-arm64:
	@echo "Only darwin-arm64 has this Makefile target."
