

BIN=meladx.go
OPTS=-v -r x3@example.com -s mx.example.com < mel1.txt

.PHONY: run
run:
	go run $(BIN) $(OPTS)
