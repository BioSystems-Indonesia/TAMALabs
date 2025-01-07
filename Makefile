.ONESHELL:
.PHONY: build

build-fe:
	cd web && npm run build

build-be:
	go build -o bin/app cmd/rest/main.go

build-be-win:
	GOOS=windows GOARCH=amd64 go build -o bin/winapp.exe cmd/rest/main.go

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

