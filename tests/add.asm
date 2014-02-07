list p=16f74
title "add test"

  org 00h
  goto main

  org 04h
  goto isrService

  org 05h

main
  movlw 7 ; w = 7
  addlw 5 ; w = 12
  movwf 21h
  movlw 5 ; w = 5
  addwf 21h, w ; w = 17

isrService
  goto isrService

  end
