.ONESHELL:
.PHONY: build

build: build-fe build-be

build-fe: npm-install web/dist

build-be: bin/app

npm-install: web/node_modules

web/dist: web/src web/package.json web/tsconfig.json web/tsconfig.app.json web/vite.config.ts web/index.html
	cd web && npm run build

web/node_modules: web/package.json web/package-lock.json
	cd web && npm install

bin/app:
	go build -ldflags "-X main.Version=$(git rev-parse --short HEAD)" -v -o bin/app ./cmd/rest

build-be-win:
	GOOS=windows GOARCH=amd64 go build -o bin/winapp.exe ./cmd/rest

build-win:
	make build-fe
	make build-be-win

dev-fe:
	cd web && npm run dev

dev-be:
	air

migrate-hash:
	atlas migrate hash

migrate-down:
	migrate -path ./migrations -database 'sqlite3://tmp/biosystem-lims.db' down 1

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

install:
	go install github.com/air-verse/air@latest
	go install github.com/akavel/rsrc@latest