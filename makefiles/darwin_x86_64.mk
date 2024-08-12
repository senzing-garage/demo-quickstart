# Makefile extensions for darwin-x86_64.

# -----------------------------------------------------------------------------
# OS-ARCH specific targets
# -----------------------------------------------------------------------------

.PHONY: build-osarch-specific
build-osarch-specific: darwin/amd64

# -----------------------------------------------------------------------------
# Makefile targets supported only by this platform.
# -----------------------------------------------------------------------------

.PHONY: only-darwin-x86_64
only-darwin-x86_64:
	$(info Only darwin-x86_64 has this Makefile target.)
