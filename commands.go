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

func unsetBreakpoint(args []string, state *emuState) error {
    if len(args) == 0 {
        return errors.New("not enough arguments")
    }
    addr, err := parseInteger(args[0], 14)
    if err != nil {
        return err
    }

    for i, bp := range state.breakpoints {
        if uint16(addr) == bp {
            // to "remove" a breakpoint, replace it with the length of
            // the code ROM. Since that address will never be reached,
            // the breakpoint effectively does not exist
            state.breakpoints[i] = uint16(len(state.code_rom))
            return nil
        }
    }

    return errors.New("could not find breakpoint")
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
    // always take first step, otherwise we'd get stuck on a breakpoint
    err := stepForward(nil, state)

    for !state.atBreakpoint() && err == nil{
        err = stepForward(nil, state)
    }

    return err
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

const INT_VECTOR = 4

func triggerInterrupt(args []string, state *emuState) error {
    intcon := getRegValue(state, REG_INTCON)

    // check that interrupts are enabled
    if intcon & (1 << INTCON_GIE) == 0 {
        return errors.New("interrupts not enabled")
    }

    // clear GIE
    intcon &= ^byte(1 << INTCON_GIE)

    // save PC and change it to interrupt vector
    state.stack.push(state.pc)
    state.pc = INT_VECTOR

    setRegValue(state, REG_INTCON, intcon)

    return nil
}
