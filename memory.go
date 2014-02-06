package main

func getRegValue(state *emuState, addr uint16) byte {
    full_addr := uint16(state.bank << RAM_BANK_SHIFT) | addr
    return state.data_ram[full_addr]
}

func setRegValue(state *emuState, addr uint16, value byte) {
    full_addr := uint16(state.bank << RAM_BANK_SHIFT) | addr
    state.data_ram[full_addr] = value
}
