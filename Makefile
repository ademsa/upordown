default: build

clean:
	rm -r ./bin && mkdir ./bin

prepare-dependencies:
	go mod download

build-cli-code:
	go build -o ./bin/upordown-cli ./cmd/cli/main.go

run-cli-code:
	./bin/upordown-cli

build-and-run-cli: build run-cli-code

build-code:
	go build -o ./bin/upordown ./cmd/web/main.go

build: clean prepare-dependencies build-code build-cli-code

run-code:
	./bin/upordown

build-and-run: build run-code