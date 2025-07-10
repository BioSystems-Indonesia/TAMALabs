.ONESHELL:
.PHONY: build

# Inno Setup Compiler command line tool.
# This assumes Inno Setup is installed in the default location.
# You may need to adjust this path if it's installed elsewhere or if you add it to your system's PATH.
ISCC="ISCC.exe"

# The Inno Setup script file
ISS_FILE=setup.iss

build: build-fe build-be

build-fe:
	cd web && npm run build

build-be: 
	go build -ldflags "-X 'main.version=$(shell git rev-parse --short HEAD)' -H windowsgui" -v -o bin/winapp.exe ./cmd/rest

build-be-win:
	GOOS=windows GOARCH=amd64 make build-be

build-win:
	make build-fe
	make build-be-win

installer: build
	@echo "Creating installer..."
	@ISCC $(ISS_FILE)
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
	atlas migrate diff --env gorm 

# Catch-all target to allow passing arguments
%:
	@:

icon:
	rsrc -arch 386 -ico favicon.ico -manifest elgatama-lims.exe.manifest
	rsrc -arch amd64 -ico favicon.ico -manifest elgatama-lims.exe.manifest
	mv rsrc_windows_amd64.syso cmd/rest
	mv rsrc_windows_386.syso cmd/rest

install:
	go install github.com/air-verse/air@latest
	go install github.com/akavel/rsrc@latest
	go install github.com/google/wire/cmd/wire@latest

wire:
	wire ./...