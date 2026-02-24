BINARY_NAME=sgrep
MAIN_PATH=./cmd/sgrep
BUILD_DIR=./build

GOOS := $(shell go env GOOS)

ifeq ($(GOOS),windows)
    BINARY=$(BINARY_NAME).exe
else
    BINARY=$(BINARY_NAME)
endif

.PHONY: build clean tidy run

build:
	go build -o $(BUILD_DIR)/$(BINARY) $(MAIN_PATH)

run: build
	./$(BUILD_DIR)/$(BINARY) -f=test.txt someone

clean:
	go clean
	rm -f sgrep sgrep.exe

tidy:
	go mod tidy
