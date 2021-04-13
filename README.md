# WASM interpreter/VM

This is an experimental WASM interpreter.  Based on
[WebAssembly spec 1.1, Mar 12 2021 draft](https://webassembly.github.io/spec/core/).  Not useful or functional yet.


## Tools
* `go`, v1.16.  For wasm support, v1.12 or later is required.
* `wat2wasm`, via `wabt`.  Or emulate with `wasmtime`, etc.  Only required for building the samples.


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
=== RUN   TestTypeSection/bad-ftype-delimiter-0xAA
--- PASS: TestTypeSection (0.00s)
    --- PASS: TestTypeSection/decode-no-param-no-result (0.00s)
    --- PASS: TestTypeSection/decode-multiple-param-multiple-result (0.00s)
    --- PASS: TestTypeSection/bad-ftype-delimiter-0xAA (0.00s)
PASS
coverage: 61.8% of statements
ok  	wasm	0.004s	coverage: 61.8% of statements
```

## Usage
```
# Invocation
dan@dan-desktop:~/src/dwasm$ ./dwasm -h
Usage: ./dwasm [options] /path/to/input.wasm
  -d	Dump .wasm sections
  -f function
    	Start/entry function
  -p value
    	Preload int32 value on stack
  -v	Validate .wasm sections
  -x	Start VM + execute

# Summarize basic sections
dan@dan-desktop:~/src/dwasm$ ./dwasm -d samples/factorial.wasm 
2021/04/12 22:30:32 Module:
Custom section:
    custom: 'name', size 20

Type section:
    function: param f64 => result f64 

Function section:
    index: [0]

Export section:
    export: 'fac', type function, index 0x0

Code section:
    function: length 43

# Execute the 'nop' example
dan@dan-desktop:~/src/dwasm$ ./dwasm -x -f fnop samples/fnop.wasm 
2021/04/12 22:31:47 VM exited cleanly

# Execute the 'addTwo' example
dan@dan-desktop:~/src/dwasm$ ./dwasm -x -f addTwo -p 3 -p 4 samples/simple.wasm 2021/04/12 22:31:53 Thread stack: 7
2021/04/12 22:31:53 Thread stack: 4
2021/04/12 22:31:53 Thread stack: 3
2021/04/12 22:31:53 VM exited cleanly

```
