# Makefile extensions for windows.

# -----------------------------------------------------------------------------
# Variables
# -----------------------------------------------------------------------------

SENZING_TOOLS_DATABASE_URL ?= sqlite3://na:na@nowhere/C:\Temp\sqlite\G2C.db

# -----------------------------------------------------------------------------
# OS specific targets
# -----------------------------------------------------------------------------

.PHONY: build-osarch-specific
build-osarch-specific: windows/amd64
	@mv $(TARGET_DIRECTORY)/windows-amd64/$(PROGRAM_NAME) $(TARGET_DIRECTORY)/windows-amd64/$(PROGRAM_NAME).exe


.PHONY: clean-osarch-specific
clean-osarch-specific:
	del /F /S /Q $(TARGET_DIRECTORY)
	del /F /S /Q $(GOPATH)/bin/$(PROGRAM_NAME)
	del /F /S /Q C:\Temp\sqlite


.PHONY: hello-world-osarch-specific
hello-world-osarch-specific:
	@echo "Hello World, from windows."


.PHONY: run-osarch-specific
run-osarch-specific:
	@go run main.go


.PHONY: setup-osarch-specific
setup-osarch-specific:
	@mkdir C:\Temp\sqlite
	@copy testdata\sqlite\G2C.db C:\Temp\sqlite\G2C.db
	@mkdir $(TARGET_DIRECTORY)\
	@mkdir $(TARGET_DIRECTORY)\$(GO_OS)-$(GO_ARCH)	


.PHONY: test-osarch-specific
test-osarch-specific:
	@go test -v -p 1 ./...

# -----------------------------------------------------------------------------
# Makefile targets supported only by this platform.
# -----------------------------------------------------------------------------

.PHONY: only-windows
only-windows:
	@echo "Only windows has this Makefile target."
