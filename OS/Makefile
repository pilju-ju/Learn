ARCH = armv7-a
MCPU = cortex-a8

CC = arm-none-eabi-gcc
AS = arm-none-eabi-as
LD = arm-none-eabi-ld
OC = arm-none-eabi-objcopy

LINKER_SCRIPT = ./pilos.ld
MAP_FILE = build/pilos.map

ASM_SRCS = $(wildcard boot/*.S)
ASM_OBJS = $(patsubst boot/%.S, build/%.o, $(ASM_SRCS))

C_SRCS = $(wildcard boot/*.c)
C_OBJS= $(patsubst boot/%.c, build/%.o, $(C_SRCS))

INC_DIRS = include

pilos = build/pilos.axf
pilos_bin = build/pilos.bin

.PHONY: all clean run debug gdb

all: $(pilos)

clean:
	@rm -fr build 

run: $(pilos)
	qemu-system-arm -M realview-pb-a8 -kernel $(pilos)

debug: $(pilos)
	qemu-system-arm -M realview-pb-a8 -kernel $(pilos) -S -gdb tcp::1234,ipv4


gdb:
	 gdb-multiarch 

$(pilos):$(ASM_OBJS) $(C_OBJS) $(LINKER_SCRIPT)
	$(LD) -n -T $(LINKER_SCRIPT) -o $(pilos) $(ASM_OBJS) $(C_OBJS) -Map=$(MAP_FILE)
	$(OC) -O binary $(pilos) $(pilos_bin)

build/%.os: $(ASM_SRCS)
	mkdir -p $(shell dirname $@)
	$(CC) -march=$(ARCH) -mcpu=$(MCPU) -I $(INC_DIRS) -c -g -o $@ $<


build/%.o: $(C_SRCS)
	mkdir -p $(shell dirname $@)
	$(CC) -march=$(ARCH) -mcpu=$(MCPU) -I $(INC_DIRS) -c -g -o $@ $<
