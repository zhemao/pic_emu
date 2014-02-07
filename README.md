# PIC Emulator

This program is an emulator and debugger for a PIC16F74 microcontroller.
Currently, the emulator supports all of the PIC16 instructions except
SLEEP and CLRWDT. Peripherals such as the watchdog timer, other timers,
and ADC are not supported. The debugger supports running, stepping forward,
setting breakpoints, and printing the contents of memory. I plan to eventually
support setting interrupts and modifying the contents of memory.
