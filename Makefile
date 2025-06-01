.ONESHELL:
.PHONY: build

build-fe:
	cd web && npm run build

build-be:
	go build -o bin/app ./cmd/rest

build-be-win:
	GOOS=windows GOARCH=amd64 go build -o bin/winapp.exe ./cmd/rest

build:
	make build-fe
	make build-be

build-win:
	make build-fe
	make build-be-win

dev-fe:
	cd web && npm run dev

dev-be:
	air

migrate-hash:
	atlas migrate hash

migrate-diff:
	$(eval ARGS := $(filter-out $@,$(MAKECMDGOALS)))
	./scripts/migrate-diff.sh $(ARGS)

# Catch-all target to allow passing arguments
%:
	@:

icon:
	rsrc -arch 386 -ico favicon.ico -manifest elgatama-lims.exe.manifest
	rsrc -arch amd64 -ico favicon.ico -manifest elgatama-lims.exe.manifest
	mv rsrc_windows_amd64.syso cmd/rest
	mv rsrc_windows_386.syso cmd/rest
