package main

import (
    "github.com/kierdavis/go/ihex"
    "github.com/GeertJohan/go.linenoise"
    "bufio"
    "os"
    "fmt"
    "strings"
    "encoding/binary"
)

const ROM_SIZE = 4096
const NUM_RAM_BANKS = 4
const RAM_BANK_SIZE = 128
const RAM_BANK_SHIFT = 7
const STACK_SIZE = 8

func BytesToWords(bytes []byte) []uint16 {
    words := make([]uint16, len(bytes) / 2)

    for i := range words {
        w := binary.LittleEndian.Uint16(bytes[2 * i : 2 * i + 2])
        words[i] = w
    }

    return words
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: pic_emu code.hex")
        os.Exit(1)
    }

    f, err := os.Open(os.Args[1])
    if err != nil {
        panic(err)
    }
    defer f.Close();

    reader := bufio.NewReader(f)
    ix, err := ihex.ReadIHex(reader)
    if err != nil {
        panic(err)
    }

    state := new(emuState)
    state.code_rom = BytesToWords(ix.ExtractDataToEnd(0))
    if len(state.code_rom) > ROM_SIZE {
        fmt.Println("ERROR: code size too large")
        os.Exit(1)
    }
    state.data_ram = make([]byte, NUM_RAM_BANKS * RAM_BANK_SIZE)
    state.breakpoints = make([]uint16, 0)
    state.pc = 0
    state.bank = 0
    state.running = false
    state.stack = newStack(STACK_SIZE)
    state.verbose = false

    commands := map[string]command {
        "r" : startRunning,
        "run" : startRunning,
        "c" : continueRunning,
        "continue" : continueRunning,
        "b" : setBreakpoint,
        "break" : setBreakpoint,
        "d" : unsetBreakpoint,
        "delete" : unsetBreakpoint,
        "s" : stepForward,
        "step" : stepForward,
        "p" : printRegister,
        "print" : printRegister,
        "v" : toggleVerbose,
        "verbose" : toggleVerbose,
        "put" : putByte,
        "set" : setBit,
        "clear" : clearBit,
        "flip" : flipBit,
        "int" : triggerInterrupt,
    }

    line, err := linenoise.Line("pic> ")

    for err == nil  {
        if line == "" {
            continue
        }
        if line == "q" {
            break
        }

        parts := strings.Split(line, " ")
        operation := commands[parts[0]]
        if operation == nil {
            fmt.Printf("No such operation %s\n", parts[0])
            line, err = linenoise.Line("pic> ")
            continue
        }

        var args []string

        if len(parts) > 1 {
            args = parts[1:]
        } else {
            args = make([]string, 0)
        }

        err = operation(args, state)
        if err != nil {
            if state.running {
                fmt.Printf("Error on instruction 0x%x\n", state.pc)
            }
            fmt.Println(err)
        }

        line, err = linenoise.Line("pic> ")
    }
}
