package arch

import (
	"debug/pe"
	"errors"
	"os"
	"strings"
)

// CPU is type of CPU architecture.
type CPU int

const (
	// X86 means Intel x86 (32 bit).
	X86 CPU = iota + 1

	// AMD64 means AMD/Intel 64 bit.
	AMD64
)

// ErrorUnknownArch is returned when failed to deetect architecture.
var ErrorUnknownArch = errors.New("unknown architecture")

// OS returns architecture of operating system.
func OS() (CPU, error) {
	v, ok := os.LookupEnv("PROCESSOR_ARCHITECTURE")
	if !ok {
		return 0, ErrorUnknownArch
	}
	switch strings.ToUpper(v) {
	case "X86":
		return X86, nil
	case "AMD64":
		return AMD64, nil
	default:
		return 0, ErrorUnknownArch
	}
}

// Exe returns architecture of execute file.
func Exe(name string) (CPU, error) {
	f, err := pe.Open(name)
	if err != nil {
		if os.IsNotExist(err) {
			return OS()
		}
		return 0, err
	}
	defer f.Close()

	switch f.FileHeader.Machine {
	case 0x014c:
		return X86, nil
	case 0x8664:
		return AMD64, nil
	}
	return 0, ErrorUnknownArch
}
