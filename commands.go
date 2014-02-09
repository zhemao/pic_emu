package main

import (
    "errors"
    "fmt"
    "strconv"
    "strings"
)

func parseInteger(numStr string, nbits int) (int16, error) {
    var num int16
    var err error
    if strings.HasSuffix(numStr, "h") {
        var hexnum uint64
        hexStr := numStr[0 : len(numStr) - 1]
        hexnum, err = strconv.ParseUint(hexStr, 16, nbits)
        num = int16(hexnum)
    } else {
        var decnum int64
        decnum, err = strconv.ParseInt(numStr, 10, nbits)
        num = int16(decnum)
    }

    if err != nil {
        return 0, err
    }

    return num, nil
}

type command func([]string, *emuState) error

func setBreakpoint(args []string, state *emuState) error {
    if len(args) == 0 {
        return errors.New("not enough arguments")
    }
    addr, err := parseInteger(args[0], 14)
    if err != nil {
        return err
    }
    state.breakpoints = append(state.breakpoints, uint16(addr))
    return nil
}

func stepForward(args []string, state *emuState) error {
    if !state.running {
        return errors.New("program isn't running\n")
    }
    if state.verbose {
        fmt.Printf("instruction %d\n", state.pc)
    }
    return executeInstruction(state.code_rom[state.pc], state)
}

func startRunning(args []string, state *emuState) error {
    state.reset()
    state.running = true
    return continueRunning(args, state)
}

func continueRunning(args []string, state *emuState) error {
    for !state.atBreakpoint() {
        err := stepForward(nil, state)
        if err != nil {
            return err
        }
    }
    return nil
}

func printRegister(args []string, state *emuState) error {
    if len(args) == 0 {
        return errors.New("not enough arguments")
    }
    var value int16
    regName := args[0]

    switch regName {
        case "pc" : value = int16(state.pc)
        case "w"  : value = int16(state.accum)
        case "tos": value = int16(state.stack.tos)
        default : {
            regAddr, err := parseInteger(regName, 9)
            if err != nil {
                return errors.New(
                    fmt.Sprintf("Unrecognized register %s\n", regName))
            }
            value = int16(state.data_ram[regAddr])
        }
    }

    fmt.Println(value)
    return nil
}

func toggleVerbose(args []string, state *emuState) error {
    state.verbose = !state.verbose
    return nil
}

func putByte(args []string, state *emuState) error {
    if len(args) < 2 {
        return errors.New("not enough arguments")
    }

    addr, err := parseInteger(args[0], 9)
    if err != nil {
        return err
    }

    value, err := parseInteger(args[1], 8)
    if err != nil {
        return err
    }

    state.data_ram[addr] = byte(value)

    return nil
}

func parseAddrAndBit(args []string) (uint16, uint8, error) {
    addr, err := parseInteger(args[0], 9)
    if err != nil {
        return 0, 0, err
    }

    bit, err := parseInteger(args[1], 4)
    if err != nil {
        return 0, 0, err
    }

    return uint16(addr), uint8(bit), nil
}

func setBit(args []string, state *emuState) error {
    if len(args) < 2 {
        return errors.New("not enough arguments")
    }

    addr, bit, err := parseAddrAndBit(args)
    if err != nil {
        return err
    }

    state.data_ram[addr] |= byte(1 << bit)
    return nil
}

func clearBit(args []string, state *emuState) error {
    if len(args) < 2 {
        return errors.New("not enough arguments")
    }

    addr, bit, err := parseAddrAndBit(args)
    if err != nil {
        return err
    }

    state.data_ram[addr] &= ^byte(1 << bit)
    return nil
}

func flipBit(args []string, state *emuState) error {
    if len(args) < 2 {
        return errors.New("not enough arguments")
    }

    addr, bit, err := parseAddrAndBit(args)
    if err != nil {
        return err
    }

    state.data_ram[addr] ^= byte(1 << bit)

    return nil
}
