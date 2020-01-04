BINARY := k8s-qos

all: clean test $(BINARY)

test: deps
	go test ./... -test.v -cover

deps:
	go get -t -d ./...

$(BINARY): deps
	go build -o $(BINARY) cmd/k8s-qos/main.go

run: build
	./$(BINARY)

clean:
	rm -f $(BINARY)

