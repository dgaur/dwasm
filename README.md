# WASM interpreter/VM

This is an experimental WASM interpreter.  Based on
[WebAssembly spec 1.1, Mar 12 2021 draft](https://webassembly.github.io/spec/core/).

## Build
```
dan@dan-desktop:~/src/dwasm$ make clean
dan@dan-desktop:~/src/dwasm$ make
dan@dan-desktop:~/src/dwasm$ make vet
dan@dan-desktop:~/src/dwasm$ make test
=== RUN   TestULEB128Decoding
=== RUN   TestULEB128Decoding/0x00
=== RUN   TestULEB128Decoding/0x01
...
=== RUN   TestTableSection/decode2
--- PASS: TestTableSection (0.00s)
    --- PASS: TestTableSection/decode1 (0.00s)
    --- PASS: TestTableSection/decodeAB (0.00s)
    --- PASS: TestTableSection/decode2 (0.00s)
PASS
coverage: 38.4% of statements
ok  	wasm	(cached)	coverage: 38.4% of statements
```

## Usage
```
dan@dan-desktop:~/src/dwasm$ ./dwasm 
Usage: ./dwasm [options] input.wasm
```
