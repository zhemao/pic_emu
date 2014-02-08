test: pic_emu test.hex
	./pic_emu test.hex

pic_emu: *.go
	go build

test.hex: test.asm
	gpasm test.asm
