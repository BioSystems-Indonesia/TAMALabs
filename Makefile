.ONESHELL:
.PHONY: build

# Inno Setup Compiler command line tool.
# This assumes Inno Setup is installed in the default location.
# You may need to adjust this path if it's installed elsewhere or if you add it to your system's PATH.
ISCC="ISCC.exe"

# The Inno Setup script file
ISS_FILE=setup.iss

build: build-fe build-be build-tray build-service-helper

build-fe:
	cd web && npm run build

build-be-win-ps:
	@echo "Building for Windows using PowerShell..."
	@powershell -Command "$$env:GOOS='windows'; $$env:GOARCH='amd64'; go build -ldflags \"-X 'main.version=$$(git rev-parse --short HEAD)' -H windowsgui\" -v -o bin/TAMALabs.exe ./cmd/rest"

build-win:
	$(MAKE) build-fe
	$(MAKE) build-be-win

# Ensure rsrc installed and generate syso inside cmd/tray
build-tray:
	@echo "Building TAMALabsTray with administrator manifest..."
	@powershell -Command "if (Test-Path 'cmd/tray/rsrc.syso') { Remove-Item 'cmd/tray/rsrc.syso' }"
	rsrc -manifest TAMALabsTray.exe.manifest -o cmd/tray/rsrc.syso
	@powershell -Command "$$env:GO111MODULE='on'; go build -ldflags \"-H windowsgui\" -v -o bin/TAMALabsTray.exe ./cmd/tray"
	@echo "TAMALabsTray built (rsrc in cmd/tray/rsrc.syso)"

# Use build-tray inside build-be and build-be-win
build-be:
	go build -ldflags "-X 'main.version=$(shell git rev-parse --short HEAD)' -H windowsgui" -v -o bin/TAMALabs.exe ./cmd/rest
	$(MAKE) build-tray

build-be-win:
	@echo "Building backend service..."
	@go env -w GOOS=windows GOARCH=amd64
	@go build -ldflags "-X 'main.version=$(shell git rev-parse --short HEAD)' -H windowsgui" -v -o bin/TAMALabs.exe ./cmd/rest
	@echo "Building tray application with manifest..."
	$(MAKE) build-tray
	@go env -u GOOS GOARCH

build-service-helper:
	@echo "Building service helper with administrator manifest..."
	@powershell -Command "if (Test-Path 'cmd/service-helper/rsrc.syso') { Remove-Item 'cmd/service-helper/rsrc.syso' }"
	rsrc -manifest cmd/service-helper/service-helper.exe.manifest -o cmd/service-helper/rsrc.syso
	@powershell -Command "$$env:GO111MODULE='on'; go build -ldflags \"-H windowsgui\" -v -o bin/service-helper.exe ./cmd/service-helper"
	@echo "Service helper built (rsrc in cmd/service-helper/rsrc.syso)"

build-tray-simple:
	go build -ldflags "-H windowsgui" -v -o bin/TAMALabsTray.exe ./cmd/tray

installer: build	
	@echo "Creating installer..."
	@docker run --rm -v "$(CURDIR):/work" amake/innosetup $(ISS_FILE)
	@echo "Installer created successfully!"

dev-fe:
	cd web && npm run dev

dev-be:
	air

migrate-hash:
	atlas migrate hash --env gorm

migrate-down:
	migrate -path ./migrations -database 'sqlite3://tmp/biosystem-lims.db' down 1

migrate-diff:
	@if not defined desc ( \
		echo Error: desc is required. Usage: make migrate-diff desc="add_some_table" & \
		exit /b 1 \
	)
	atlas migrate diff --env gorm $(desc)

# Catch-all target to allow passing arguments
%:
	@:

icon:
	rsrc -arch 386 -ico favicon.ico -manifest TAMALabs.exe.manifest
	rsrc -arch amd64 -ico favicon.ico -manifest TAMALabs.exe.manifest
	mv rsrc_windows_amd64.syso cmd/rest
	mv rsrc_windows_386.syso cmd/rest

install:
	go install github.com/air-verse/air@latest
	go install github.com/akavel/rsrc@latest
	go install github.com/google/wire/cmd/wire@latest

wire:
	wire ./...