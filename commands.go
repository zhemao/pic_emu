package main

import (
    "errors"
    "fmt"
    "strconv"
)

type command func([]string, *emuState) error

func setBreakpoint(args []string, state *emuState) error {
    if len(args) == 0 {
        return errors.New("not enough arguments")
    }
    addr, err := strconv.ParseUint(args[0], 10, 16)
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

func toggleVerbose(args []string, state *emuState) error {
    state.verbose = !state.verbose
    return nil
}
