.PHONY: all
all:
	@$(MAKE) deps
	@$(MAKE) build

.PHONY: deps
deps:
	go mod download -x

.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o build/bin/telegraftoplayfab ./cmd/
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -a -o build/bin/telegraftoplayfab.exe ./cmd/
	cp cmd/plugin.conf build/bin/

	zip -j build/telegraftoplayfablinux.zip build/bin/telegraftoplayfab build/bin/plugin.conf
	zip -j build/telegraftoplayfabwindows.zip build/bin/telegraftoplayfab.exe build/bin/plugin.conf

.PHONY: test
test:
	go test -race -covermode=atomic -coverprofile=covprofile -coverprofile=coverage.out -parallel 4 -timeout 2h ./...

.PHONY: tidy
tidy:
	go mod verify
	go mod tidy

.PHONY: clean
clean:
	rm -rf build