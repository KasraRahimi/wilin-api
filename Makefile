SRC_FILES=*.go
MAIN_FILE=main.go
TRG_FILE=$(MAIN_FILE:%.go=%)

SRC_DIR=.
MAIN_DIR=.
TRG_DIR=./trg

SRC=$(shell find $(SRC_DIR) -name '$(SRC_FILES)')
MAIN=$(MAIN_DIR)/$(MAIN_FILE)
TRG=$(TRG_DIR)/$(TRG_FILE)

all: $(TRG)

run: $(TRG)
	$^

$(TRG): $(MAIN) $(SRC)
	mkdir -p $(TRG_DIR)
	go build -o $@ $<

test:
	go test -v ./...

clean:
	rm -f ./trg/*