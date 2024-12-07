.ONESHELL:
.PHONY: build

build-fe:
	cd web && npm run build
	statik -src=./web/dist

build-be:
	go build -o bin/app cmd/rest/main.go

build:
	make build-fe
	make build-be

dev-fe:
	cd web && npm run dev

dev-be:
	air --build.cmd "go build -o bin/app cmd/rest/main.go" --build.bin "./bin/app" --build.exclude_dir "node_modules,bin,web"

