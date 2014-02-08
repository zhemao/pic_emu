package main

import (
    "github.com/kierdavis/go/ihex"
    "github.com/GeertJohan/go.linenoise"
    "bufio"
    "os"
    "fmt"
    "errors"
    "strings"
    "strconv"
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

type command func(string, *emuState) error

func setBreakpoint(arg string, state *emuState) error {
    addr, err := strconv.ParseUint(arg, 16, 16)
    if err != nil {
        return err
    }
    state.breakpoints = append(state.breakpoints, uint16(addr))
    return nil
}

func stepForward(arg string, state *emuState) error {
    if !state.running {
        return errors.New("program isn't running\n")
    }
    return executeInstruction(state.code_rom[state.pc], state)
}

func startRunning(arg string, state *emuState) error {
    state.reset()
    state.running = true
    return continueRunning(arg, state)
}

func continueRunning(arg string, state *emuState) error {
    for !state.atBreakpoint() {
        err := stepForward("", state)
        if err != nil {
            return err
        }
    }
    return nil
}

func printRegister(regName string, state *emuState) error {
    var value int16

    switch regName {
        case "pc" : value = int16(state.pc)
        case "w"  : value = int16(state.accum)
        case "tos": value = int16(state.stack.tos)
        default : {
            regAddr, err := strconv.ParseUint(regName, 16, 8)
            if err != nil {
                return errors.New(
                    fmt.Sprintf("Unrecognized register %s\n", regName))
            }
            value = int16(getRegValue(state, uint16(regAddr)))
        }
    }

    fmt.Println(value)
    return nil
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

    commands := map[string]command {
        "r" : startRunning,
        "c" : continueRunning,
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

        err = operation(arg, state)
        if err != nil {
            if state.running {
                fmt.Printf("Error on instruction 0x%x\n", state.pc)
            }
            state.running = false
            fmt.Println(err)
        }

        line, err = linenoise.Line("pic> ")
    }
}
