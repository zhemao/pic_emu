package main

import "errors"

func decodeDF(instr uint16) (uint16, uint16) {
    d := (instr >> 7) & 1
    f := instr & 0x7f

    return d, f
}

func executeInstruction0(instr uint16, state *emuState) error {
    opcode := (instr >> 8) & 0xf
    d := (instr >> 7) & 1
    f := instr & 0x7f

    var newVal int16
    accumVal := int16(state.accum)
    regVal := int16(getRegValue(state, f))

    status := getRegValue(state, REG_STATUS)

    c := byte((status >> STATUS_C) & 0x1)

    if opcode == 0 {
        // MOVWF
        if d == 1 {
            setRegValue(state, f, byte(accumVal))
        } else if (f & 0xfe) == 8 {
            // RETURN or RETFIE
            newInstr, err := state.stack.pop()
            if err != nil {
                panic(err)
            }
            state.pc = newInstr
            if f == 9 {
                intcon := getRegValue(state, REG_INTCON)
                intcon |= 1 << INTCON_GIE
                setRegValue(state, REG_INTCON, intcon)
            }
            return nil
        }
        if f == 0x63 {
            state.pc++
            return errors.New("entering standby mode\n")
        }
        if (f & 0x1f) != 0 {
            // not NOP
            return errors.New("invalid instruction\n")
        }

        state.pc++
        return nil
    }

    switch opcode {
        case 0x1: newVal = 0 // CLRF or CLRW
        case 0x2: newVal = regVal - accumVal // SUBWF
                  c = byte(((newVal ^ regVal ^ -accumVal) >> 8) | 0x1)
        case 0x3: newVal = regVal - 1 // DECF
        case 0x4: newVal = regVal | accumVal // IORWF
        case 0x5: newVal = regVal & accumVal // ANDWF
        case 0x6: newVal = regVal ^ accumVal // XORWF
        case 0x7: newVal = regVal + accumVal // ADDWF
                  c = byte(((newVal ^ regVal ^ accumVal) >> 8) | 0x1)
        case 0x8: newVal = regVal // MOVF
        case 0x9: newVal = ^regVal // COMF
        case 0xa: newVal = regVal + 1 // INCF
        case 0xb: newVal = regVal - 1 // DECFSZ
        case 0xc: newVal = ((regVal << 1) | int16(c)) & 0xff
                  c = byte(regVal | 0x1) // RLF
        case 0xd: newVal = ((regVal >> 1) | int16(c << 7)) & 0xff
                  c = byte((regVal >> 7) | 0x1) // RRF
        case 0xe: newVal = ((regVal >> 4) & 0x7) | ((regVal & 0x7) << 4) // SWAPF
        case 0xf: newVal = regVal + 1 // INCFSZ
    }

    var z byte
    if (newVal == 0) {
        z = 1
    } else {
        z = 0
    }

    dc := byte(((newVal ^ regVal ^ accumVal) >> 4) & 0x1)

    if (d == 1) {
        setRegValue(state, f, byte(newVal))
    } else {
        state.accum = byte(newVal)
    }

    skipInstr := (opcode == 0xb || opcode == 0xf)
    addsubInstr := (opcode == 0x7 || opcode == 0x2)
    rotInstr := (opcode == 0xc || opcode == 0xd)
    zeroAffectInstr := (opcode != 0x0 && !skipInstr && opcode != 0xe)

    if (addsubInstr) {
        status |= (dc << STATUS_DC)
    }

    if (addsubInstr || rotInstr) {
        status |= (c << STATUS_C)
    }

    if (zeroAffectInstr) {
        status |= (z << STATUS_Z)
    }

    setRegValue(state, REG_STATUS, status)

    if skipInstr && z == 1 {
        state.pc += 2
    } else {
        state.pc++
    }

    return nil
}


func executeInstruction1(instr uint16, state *emuState) error {
    return errors.New("not supported yet")
}
func executeInstruction2(instr uint16, state *emuState) error {
    call := ((instr >> 11) & 0x1) == 1
    addr := instr & 0x7ff

    if (call) {
        state.stack.push(state.pc + 1)
    }
    state.pc = addr

    return nil
}
func executeInstruction3(instr uint16, state *emuState) error {
    opcode := (instr >> 8) & 0xf
    k := int16(instr & 0xff)

    accumVal := int16(state.accum)
    var newVal int16

    if opcode & 0xc == 0 {
        // MOVLW
        newVal = k
    } else if opcode & 0xc == 0x4 {
        // RETLW
        state.accum = byte(k)
        newpc, err := state.stack.pop()
        if err != nil {
            return err
        }
        state.pc = newpc
        return nil
    } else if opcode & 0xe == 0xc {
        newVal = k - accumVal
    } else if opcode & 0xe == 0xe {
        newVal = k + accumVal
    } else {
        switch opcode {
            case 0x8: newVal = k | accumVal
            case 0x9: newVal = k & accumVal
            case 0xa: newVal = k ^ accumVal
        }
    }

    state.accum = byte(newVal)
    state.pc++

    return nil
}

func executeInstruction(instr uint16, state *emuState) error {
    opcodeClass := (instr >> 12) & 3

    switch opcodeClass {
        case 0: return executeInstruction0(instr, state)
        case 1: return executeInstruction1(instr, state)
        case 2: return executeInstruction2(instr, state)
        case 3: return executeInstruction3(instr, state)
    }

    return errors.New("invalid instruction")
}
