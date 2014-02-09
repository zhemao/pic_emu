list p=16f74
title "add test"

#include <p16f74.inc>

  org 00h
  bsf INTCON, GIE
  goto main

  org 04h
  goto isrService

  org 05h

main
  clrw ; w = 0, Z is set
  ;; test zero flag getting unset after arithmetic
  movlw 7 ; w = 7
  addlw 5 ; w = 12, Z is unset
  btfsc STATUS, Z
  goto errout
  ;; test DC gets set properly
  movwf 21h
  movlw 5 ; w = 5
  addwf 21h, w ; w = 17, DC is set
  btfss STATUS, DC
  goto errout
  ;; test arithmetic is correct
  sublw 11h
  btfss STATUS, Z
  goto errout
  ;; test right shift, left shift OK
  movlw 3
  movwf 21h
  bcf STATUS, C
  bcf STATUS, Z
  rrf 21h, f
  rlf 21h, f
  sublw 3
  btfss STATUS, Z
  goto errout
  ;; test C gets set properly
  addlw 80h
  addlw 80h ;; 128 + 128 = 256
  btfss STATUS, C
  goto errout

  ;; test INC/DEC functions
  clrf 21h ; 0
  incf 21h, f ; 1
  incf 21h, f ; 2
  decf 21h, f ; 1
  decfsz 21h, f ; 0
  goto errout
  decf 21h, f ; -1
  incfsz 21h, f
  goto errout

  call testfunc

  ;; test bitwise instructions
  bsf 21h, 0 ; 21h = 1
  btfss 21h, 0
  goto errout
  movf 21h, w ; w = 1
  bsf 21h, 1 ; 21h = 3
  xorwf 21h, w ; w = 2
  sublw 2h ; w = 0
  btfss STATUS, Z
  goto errout
  movf 21h, w ; w = 3
  bcf 21h, 0 ; 21h = 2
  andwf 21h, w ; w = 2
  iorlw 4h ; w = 6
  sublw 6h ; w = 0
  btfss STATUS, Z
  goto errout

  ;; exit normally
  retlw 0

testfunc
  return
  goto errout

errout
  ;; exit with error
  retlw 1

isrService
  btfss 22h, 0
  goto isrService
  retfie

  end
