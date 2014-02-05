package main

import (
    "github.com/kierdavis/go/ihex"
    "github.com/GeertJohan/go.linenoise"
    "bufio"
    "os"
    "fmt"
    "strings"
    "strconv"
    "encoding/binary"
)

const ROM_SIZE = 4096
const NUM_RAM_BANKS = 4
const RAM_BANK_SIZE = 128

func BytesToWords(bytes []byte) []uint16 {
    words := make([]uint16, len(bytes) / 2)

    for i := range words {
        w := binary.LittleEndian.Uint16(bytes[2 * i : 2 * i + 2])
        words[i] = w
    }

    return words
}

type emuState struct {
    code_rom []uint16
    data_ram []byte
    accum byte
    breakpoints []uint16
    pc uint16
    bank uint8
}

type command func(string, *emuState)

func atBreakpoint(state *emuState) bool {
    for _, bp := range state.breakpoints {
        if state.pc == bp {
            return true
        }
    }
    return false
}

func setBreakpoint(arg string, state *emuState) {
    addr, err := strconv.ParseUint(arg, 16, 16)
    if err != nil {
        panic(err)
    }
    state.breakpoints = append(state.breakpoints, uint16(addr))
}

func stepForward(arg string, state *emuState) {
    executeInstruction(state.code_rom[state.pc], state)
}

func runCode(arg string, state *emuState) {
    for !atBreakpoint(state) {
        stepForward("", state)
    }
}

func printRegister(regName string, state *emuState) {
    var value int16

    switch regName {
        case "pc" : value = int16(state.pc)
        case "w"  : value = int16(state.accum)
        default : {
            regAddr, err := strconv.ParseUint(regName, 16, 8)
            if err != nil {
                fmt.Printf("Unrecognized register %s\n", regName)
                return
            }
            value = int16(getRegValue(state, uint16(regAddr)))
        }
    }

    fmt.Println(value)
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

    commands := map[string]command {
        "r" : runCode,
        "b" : setBreakpoint,
        "n" : stepForward,
        "p" : printRegister,
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
        arg := ""
        if len(parts) > 1 {
            arg = parts[1]
        }

        operation(arg, state)

        line, err = linenoise.Line("pic> ")
    }
}
