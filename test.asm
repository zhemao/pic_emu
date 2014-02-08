list p=16f74
title "add test"

#include <p16f74.inc>

  org 00h
  goto main

  org 04h
  goto errout

  org 05h

main
  clrw ; w = 0, Z is set
  movlw 7 ; w = 7
  addlw 5 ; w = 12, Z is unset
  btfsc STATUS, Z
  goto errout
  movwf 21h
  movlw 5 ; w = 5
  addwf 21h, w ; w = 17, DC is set
  btfss STATUS, DC
  goto errout
  sublw 11h
  btfss STATUS, Z
  goto errout
  movlw 3
  movwf 21h
  bcf STATUS, C
  bcf STATUS, Z
  rrf 21h, f
  rlf 21h, f
  sublw 3
  btfss STATUS, Z
  goto errout
  retlw 0

errout
  retlw 1
  end
