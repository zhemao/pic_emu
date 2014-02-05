list p=16f74
title "add test"

  org 00h
  goto main

  org 04h
  goto isrService

  org 05h

main
  clrw
  nop

isrService
  goto isrService

  end
