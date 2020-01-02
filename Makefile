BINARY := k8s-qos

all: clean test $(BINARY)

test: deps
	go test ./... -test.v

deps:
	go get -d ./...

$(BINARY): deps
	go build -o $(BINARY) cmd/k8s-qos/main.go

run: build
	./$(BINARY)

clean:
	rm $(BINARY)

