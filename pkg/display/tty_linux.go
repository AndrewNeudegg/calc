//go:build linux

package display

import (
	"os"
	"syscall"
	"unsafe"
)

// RawState is the saved terminal state used for restoration.
type RawState = syscall.Termios

// Linux termios constants
const (
	TCGETS = syscall.TCGETS
	TCSETS = syscall.TCSETS
)

// isATTY checks if the given file descriptor is a terminal.
func isATTY(fd uintptr) bool {
	var termios syscall.Termios
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, fd, TCGETS, uintptr(unsafe.Pointer(&termios)))
	return errno == 0
}

// enableRawMode puts the terminal into raw mode and returns the previous state.
func enableRawMode(fd int) (*RawState, error) {
	var orig syscall.Termios
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), TCGETS, uintptr(unsafe.Pointer(&orig)))
	if errno != 0 {
		return nil, os.NewSyscallError("ioctl TCGETS", errno)
	}

	raw := orig
	// Input flags - clear: BRKINT, ICRNL, INPCK, ISTRIP, IXON
	raw.Iflag &^= syscall.BRKINT | syscall.ICRNL | syscall.INPCK | syscall.ISTRIP | syscall.IXON
	// Output flags - clear: OPOST
	raw.Oflag &^= syscall.OPOST
	// Control flags - set: CS8
	raw.Cflag |= syscall.CS8
	// Local flags - clear: ECHO, ICANON, IEXTEN, ISIG
	raw.Lflag &^= syscall.ECHO | syscall.ICANON | syscall.IEXTEN | syscall.ISIG
	// Control chars
	raw.Cc[syscall.VMIN] = 1
	raw.Cc[syscall.VTIME] = 0

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
