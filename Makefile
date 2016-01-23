realm: ./cmd/realm/main.go ./server.go ./zone.go ./registry.go
	go build ./cmd/realm

clean:
	rm -f ./realm

run:
	go run ./cmd/realm/main.go ${ARGS}

.PHONY: clean run
