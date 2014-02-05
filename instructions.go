package main

func decodeDF(instr uint16) (uint16, uint16) {
    d := (instr >> 7) & 1
    f := instr & 0x7f

    return d, f
}

func executeInstruction0(instr uint16, state *emuState) {
    opcode := (instr >> 8) & 0xf
    d := (instr >> 7) & 1
    f := instr & 0x7f

    var newVal byte
    accumVal := state.accum
    regVal := getRegValue(state, f)

    switch opcode {
        case 0x1: newVal = 0
        case 0x2: newVal = regVal - accumVal
        case 0x3: newVal = regVal - 1
        case 0x4: newVal = regVal | accumVal
        case 0x5: newVal = regVal & accumVal
        case 0x6: newVal = regVal ^ accumVal
        case 0x7: newVal = regVal + accumVal
        case 0x8: newVal = regVal
        case 0x9: newVal = ^regVal
        case 0xa: newVal = regVal + 1
    }

    if (d == 1) {
        setRegValue(state, f, newVal)
    } else {
        state.accum = newVal
    }

    state.pc++
}


func executeInstruction1(instr uint16, state *emuState) {
}
func executeInstruction2(instr uint16, state *emuState) {
    call := ((instr >> 11) & 0x1) == 1
    addr := instr & 0x7ff

    if (call) {
    }
    state.pc = addr
}
func executeInstruction3(instr uint16, state *emuState) {
}

func executeInstruction(instr uint16, state *emuState) {
    opcodeClass := (instr >> 12) & 3

    switch opcodeClass {
        case 0: executeInstruction0(instr, state)
        case 1: executeInstruction1(instr, state)
        case 2: executeInstruction2(instr, state)
        case 3: executeInstruction3(instr, state)
    }
}
