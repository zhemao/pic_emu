# PIC Emulator

This program is an emulator and debugger for a PIC16F74 microcontroller.
Currently, the emulator supports all of the PIC16 instructions except
SLEEP and CLRWDT. Peripherals such as the watchdog timer, other timers,
and ADC are not supported. The debugger supports running, stepping forward,
setting breakpoints, printing the contents of memory, writing arbitrary bytes
to memory, manipulating single bits in memory, and triggering interrupts.
Program usage and debugger commands are documented below.

## Usage

    pic_emu code.hex

The emulator uses as its input format Intel Hex files with each instruction
taking up 2 bytes in little endian order. Hex files like this can be
produced from PIC assembly code using the GNU PIC assembler (`gpasm`) from the
[GNU PIC utilities](http://gputils.sourceforge.net/) package.

## Debugger Commands

### A note on integers

Several of the commands take integer arguments. By default, the arguments
given will be parsed as signed decimal integers. However, you can specify
unsigned hex numbers by appending `h` to the end, as you do in the PIC
assembler. So `11h` is decimal integer 17.

### Controlling Execution

`run` - Resets the processors and begins executing instructions starting
from instruction 0. Continues until an error occurs, the program terminates
normally, or a breakpoint is reached.
This command has an abbreviation `r`.

`continue` - Resumes executing instructions from the current processor state
and continues until an error occurs, the program terminates normally, or
a breakpoint is reached. This command has abbreviation `c`.

`int` - Triggers an interrupt. When this command is run, the GIE bit of
the INTCON register is checked. If it is set, the current value of
the program counter is pushed onto the stack and pc is changed to the
interrupt vector (address 4). The GIE bit in the INTCON register is also
cleared.

`step` - Executes a single instruction and then returns control to the user.
This command has the abbreviation `s`.

### Breakpoints

`break addr` - Set a breakpoint at the given instruction address.
This command has abbreviation `b`.

`delete addr` - Removes the breakpoint at the given instruction address.
This command has abbreviation `d`.

### Accessing memory

*Note* - Addresses given to commands in the debugger are always absolute
addresses. That is, they are not affected by banking.

`print addr` - Prints the contents of data memory at the given address.
The address argument can also be one of three special values, "pc" for the
program counter, "w" for the accumulator register, or "tos" for the top of
the instruction stack. This command has abbreviation `p`.

`put addr value` - Writes the given value to memory at the given address.

`set addr bit` - Sets the given bit in the byte at the given address. The bits
are labeled little endian (so bit 0 is the LSB and bit 7 is the MSB).

`clear addr bit` - Clears the given bit in the byte at the given address.

`flip addr bit` - Flips the given bit in the byte at the given address.
