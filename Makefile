.ONESHELL:
.PHONY: build build-fe build-be build-be-win build-tray build-service-helper build-integration-service installer dev-fe dev-be migrate-hash migrate-down migrate-diff icon install wire release

# Inno Setup Compiler command line tool
ISCC = ISCC.exe

# Inno Setup script file
ISS_FILE = setup.iss

# Default build (semua komponen)
build: build-fe build-be build-tray build-integration-service

# Build Frontend (React)
build-fe:
	cd web && npm run build

# Build Backend (Linux/macOS)
build-be:
	@echo "Building backend (local platform)..."
	go build -ldflags "-X 'main.version=$(shell git rev-parse --short HEAD)' -H windowsgui" -v -o bin/TAMALabs.exe ./cmd/rest

# Build Backend (Windows cross-compile)
build-be-win:
	@echo "Building backend for Windows..."
	@go env -w GOOS=windows GOARCH=amd64
	go build -ldflags "-X 'main.version=$(shell git rev-parse --short HEAD)' -H windowsgui" -v -o bin/TAMALabs.exe ./cmd/rest
	@go env -u GOOS GOARCH

# Build Tray Application
build-tray:
	go build -ldflags "-H windowsgui" -v -o bin/TAMALabsTray.exe ./cmd/tray
	@echo "✅ Tray built successfully."

# Build Service Helper
build-service-helper:
	@echo "Building service-helper with manifest..."
	@powershell -Command "if (Test-Path 'cmd/service-helper/rsrc.syso') { Remove-Item 'cmd/service-helper/rsrc.syso' }"
	rsrc -manifest cmd/service-helper/service-helper.exe.manifest -o cmd/service-helper/rsrc.syso
	go build -ldflags "-H windowsgui" -v -o bin/service-helper.exe ./cmd/service-helper
	@echo "✅ Service-helper built successfully."

# Build Integration Service
build-integration-service:
	@echo "Building integration service..."
	cd integration-service && go build -ldflags -v -o ../bin/TAMALabsIntegration.exe .
	@echo "✅ Integration service built successfully."

# Simple Tray Build (no manifest)
build-tray-simple:
	go build -ldflags "-H windowsgui" -v -o bin/TAMALabsTray.exe ./cmd/tray

# Generate version.ini from version.go
gen-version:
	@echo "Generating version.ini from version.go..."
	@go run scripts/gen-version-ini.go

# Create Windows Installer using Docker (Inno Setup)
installer: build gen-version
	@echo "Creating installer..."
	@docker run --rm -v "$(CURDIR):/work" amake/innosetup $(ISS_FILE)
	@echo "✅ Installer created successfully!"

# Development
dev-fe:
	cd web && npm run dev

dev-be:
	air

# Database migrations
migrate-hash:
	atlas migrate hash --env gorm

migrate-down:
	migrate -path ./migrations -database 'sqlite3://tmp/biosystem-lims.db' down 1

migrate-diff:
	@if [ -z "$(desc)" ]; then \
		echo "Error: desc is required. Usage: make migrate-diff desc='add_some_table'"; \
		exit 1; \
	fi
	atlas migrate diff --env gorm $(desc)

# Generate icon resources
icon:
	rsrc -arch 386 -ico favicon.ico -manifest TAMALabs.exe.manifest
	rsrc -arch amd64 -ico favicon.ico -manifest TAMALabs.exe.manifest
	mv rsrc_windows_amd64.syso cmd/rest
	mv rsrc_windows_386.syso cmd/rest

# Install dependencies
install:
	go install github.com/air-verse/air@latest
	go install github.com/akavel/rsrc@latest
	go install github.com/google/wire/cmd/wire@latest

# Dependency Injection generator
wire:
	wire ./...

# Catch-all (ignore unknown targets)
%:
	@:

# release: build Windows release artifacts (cross-compile backend + tray)
# Usage: make release
release:
	@echo "Building Windows release artifacts..."
	@$(MAKE) build-be-win
	@$(MAKE) build-tray
	@echo "✅ Release artifacts are in ./bin (TAMALabs.exe, TAMALabsTray.exe)"
	@echo "Run 'make installer' to build the installer if needed."
