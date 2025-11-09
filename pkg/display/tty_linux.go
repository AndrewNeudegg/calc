//go:build linux

package display

import (
	"os"
	"syscall"
	"unsafe"
)

// Termios structure for Linux
type Termios struct {
	Iflag  uint32
	Oflag  uint32
	Cflag  uint32
	Lflag  uint32
	Line   uint8
	Cc     [32]uint8
	Ispeed uint32
	Ospeed uint32
}

// RawState is the saved terminal state used for restoration.
type RawState = Termios

// Linux termios constants
const (
	TCGETS  = 0x5401
	TCSETS  = 0x5402
	BRKINT  = 0x0002
	ICRNL   = 0x0100
	INPCK   = 0x0010
	ISTRIP  = 0x0020
	IXON    = 0x0400
	OPOST   = 0x0001
	CS8     = 0x0030
	ECHO    = 0x0008
	ICANON  = 0x0002
	IEXTEN  = 0x8000
	ISIG    = 0x0001
	VMIN    = 6
	VTIME   = 5
)

// isATTY checks if the given file descriptor is a terminal.
func isATTY(fd uintptr) bool {
	var termios Termios
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, fd, TCGETS, uintptr(unsafe.Pointer(&termios)))
	return errno == 0
}

// enableRawMode puts the terminal into raw mode and returns the previous state.
func enableRawMode(fd int) (*RawState, error) {
	var orig Termios
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), TCGETS, uintptr(unsafe.Pointer(&orig)))
	if errno != 0 {
		return nil, os.NewSyscallError("ioctl TCGETS", errno)
	}

	raw := orig
	// Input flags - clear: BRKINT, ICRNL, INPCK, ISTRIP, IXON
	raw.Iflag &^= BRKINT | ICRNL | INPCK | ISTRIP | IXON
	// Output flags - clear: OPOST
	raw.Oflag &^= OPOST
	// Control flags - set: CS8
	raw.Cflag |= CS8
	// Local flags - clear: ECHO, ICANON, IEXTEN, ISIG
	raw.Lflag &^= ECHO | ICANON | IEXTEN | ISIG
	// Control chars
	raw.Cc[VMIN] = 1
	raw.Cc[VTIME] = 0

	_, _, errno = syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), TCSETS, uintptr(unsafe.Pointer(&raw)))
	if errno != 0 {
		return nil, os.NewSyscallError("ioctl TCSETS", errno)
	}
	return &orig, nil
}

// restoreRawMode restores a previously saved terminal state.
func restoreRawMode(fd int, state *RawState) {
	if state == nil {
		return
	}
	syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), TCSETS, uintptr(unsafe.Pointer(state)))
}
