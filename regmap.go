package main

// bank 0
const (
    REG_TMR0 = 0x01
    REG_PCL = 0x02
    REG_STATUS = 0x03
    REG_FSR = 0x04
    REG_PORTA = 0x05
    REG_PORTB = 0x06
    REG_PORTC = 0x07
    REG_PORTD = 0x08
    REG_PORTE = 0x09
    REG_PCLATH = 0x0a
    REG_INTCON = 0x0b
    REG_PIR1 = 0x0c
    REG_PIR2 = 0x0d
    REG_TMR1L = 0x0e
    REG_TMR1H = 0x0f
    REG_T1CON = 0x10
    REG_TMR2 = 0x11
    REG_T2CON = 0x12
    REG_SSPBUF = 0x13
    REG_SSPCON = 0x14
    REG_CCPR1L = 0x15
    REG_CCPR1H = 0x16
    REG_CCP1CON = 0x17
    REG_RCSTA = 0x18
    REG_TXREG = 0x19
    REG_RCREG = 0x1a
    REG_CCPR2L = 0x1b
    REG_CCPR2H = 0x1c
    REG_CCP2CON = 0x1d
    REG_ADRES = 0x1e
    REG_ADCON0 = 0x1f
)

// bank 1
const (
    REG_TRISA = 0x05
    REG_TRISB = 0x06
    REG_TRISC = 0x07
    REG_TRISD = 0x08
    REG_TRISE = 0x09
    REG_PIE1 = 0x0c
    REG_PIE2 = 0x0d
    REG_PCON = 0x0e
    REG_PR2 = 0x12
    REG_SSPADD = 0x13
    REG_SSPSTAT = 0x14
    REG_TXSTA = 0x18
    REG_SPBRG = 0x19
    REG_ADCON1 = 0x1f
)

// status bits
const (
    STATUS_C = 0
    STATUS_DC = 1
    STATUS_Z = 2
    STATUS_RP = 5
)

const (
    INTCON_GIE = 7
)
