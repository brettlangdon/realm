realm: ./cmd/realm/main.go ./server.go ./zone.go
	go build ./cmd/realm

clean:
	rm -f ./realm

run:
	go run ./cmd/realm/main.go

.PHONY: clean run
