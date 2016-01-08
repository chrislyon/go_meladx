

BIN=meladx.go
OPTS=-v -r x3@sra.fr < mel1.txt

.PHONY: run
run:
	go run $(BIN) $(OPTS)
