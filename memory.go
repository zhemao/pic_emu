package main

func getRegValue(state *emuState, addr uint16) byte {
    full_addr := uint16(4 * state.bank) + addr
    return state.data_ram[full_addr]
}

func setRegValue(state *emuState, addr uint16, value byte) {
    full_addr := uint16(4 * state.bank) + addr
    state.data_ram[full_addr] = value
}
