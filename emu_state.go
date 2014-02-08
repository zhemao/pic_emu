package main

type emuState struct {
    code_rom []uint16
    data_ram []byte
    accum byte
    breakpoints []uint16
    pc uint16
    bank uint8
    stack *callStack
    running bool
    verbose bool
}

func (state *emuState) atBreakpoint() bool {
    for _, bp := range state.breakpoints {
        if state.pc == bp {
            return true
        }
    }
    return false
}

func (state *emuState) reset() {
    state.accum = 0
    state.pc = 0
    state.bank = 0
    state.running = false
    state.stack.clear()

    for i, _ := range state.data_ram {
        state.data_ram[i] = 0
    }
}
