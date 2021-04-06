#
# Makefile for the dwasm build
#

# Assume tools are available in the PATH if not explicitly overridden
GO ?= go
WAT2WASM ?= wat2wasm


# Sample .wasm binaries for exercising the VM
SAMPLES_WAT := $(wildcard samples/*.wat)
SAMPLES_WASM := $(SAMPLES_WAT:.wat=.wasm)


# By default, build everything: VM, sample code, etc
DWASM := dwasm
all: $(DWASM) $(SAMPLES_WASM)


# Main CLI/VM binary
$(DWASM): $(wildcard *.go) $(wildcard wasm/*.go) Makefile
	@$(GO) build


# Compile a single .wat source file into the corresponding .wasm
%.wasm : %.wat
	@$(WAT2WASM) $< -o $@

# Compile a single .go source file into the corresponding .wasm
%.wasm : %.go
	@GOOS=js GOARCH=wasm $(GO) build -o $@ $<


.PHONY: clean
clean:
	@rm -f $(DWASM) $(SAMPLES_WASM)
	@$(GO) clean


.PHONY: test
test:
	@$(GO) test -v -cover wasm


.PHONY: vet
vet:
	@$(GO) vet
	
