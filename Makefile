

BIN=meladx.go
BINAIRE=go_meladx
OPTS=-v -r x3@sra.fr < mel1.txt

RUNTIME=/ado/X3V6/runtime



.PHONY: run
run:
	go run $(BIN) $(OPTS)

binaire:
	go build -o $(BINAIRE) $(BIN)

install: binaire
	cp $(BINAIRE) $(RUNTIME)/bin
	chown x3.dba $(RUNTIME)/bin/$(BINAIRE)
