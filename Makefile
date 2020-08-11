default: build

clean:
	rm -r ./bin && mkdir ./bin

prepare-dependencies:
	go mod download

build-code:
	go build -o ./bin/upordown ./cmd/main.go

build: clean prepare-dependencies build-code

run-code:
	./bin/upordown

build-and-run: build run-code