#
# Makefile for the dwasm build
#

# Assume the 'go' tools are available in the PATH if not explicitly overridden
GO ?= go

DWASM := dwasm

# Sample .wasm binaries for exercising the VM
HELLO_WORLD := samples/hello_world.wasm
SAMPLES := $(HELLO_WORLD)



all: $(DWASM) $(SAMPLES)

$(DWASM): $(wildcard *.go) Makefile
	@$(GO) build

# Compile a single .go source file into the corresponding .wasm
%.wasm : %.go
	@GOOS=js GOARCH=wasm $(GO) build -o $@ $<

.PHONY: clean
clean:
	@rm -f $(DWASM) $(SAMPLES)
	@$(GO) clean


.PHONY: test
test:
	@$(GO) test -v -cover


.PHONY: vet
vet:
	@$(GO) vet
	
