//go:build darwin

package display

import (
	"os"
	"syscall"
	"unsafe"
)

// RawState is the saved terminal state used for restoration.
type RawState = syscall.Termios

// isATTY checks if the given file descriptor is a terminal (best-effort on darwin).
func isATTY(fd uintptr) bool {
	// Best-effort: query termios; if it works, it's a TTY.
	var tio syscall.Termios
	_, _, e := syscall.Syscall6(syscall.SYS_IOCTL, fd, uintptr(syscall.TIOCGETA), uintptr(unsafe.Pointer(&tio)), 0, 0, 0)
	return e == 0
}

// enableRawMode puts the terminal into raw mode and returns the previous state.
func enableRawMode(fd int) (*RawState, error) {
	var orig syscall.Termios
	if _, _, err := syscall.Syscall6(syscall.SYS_IOCTL, uintptr(fd), uintptr(syscall.TIOCGETA), uintptr(unsafe.Pointer(&orig)), 0, 0, 0); err != 0 {
		return nil, os.NewSyscallError("ioctl TIOCGETA", err)
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

	if _, _, err := syscall.Syscall6(syscall.SYS_IOCTL, uintptr(fd), uintptr(syscall.TIOCSETA), uintptr(unsafe.Pointer(&raw)), 0, 0, 0); err != 0 {
		return nil, os.NewSyscallError("ioctl TIOCSETA", err)
	}
	return &orig, nil
}

// restoreRawMode restores a previously saved terminal state.
func restoreRawMode(fd int, state *RawState) {
	if state == nil {
		return
	}
	syscall.Syscall6(syscall.SYS_IOCTL, uintptr(fd), uintptr(syscall.TIOCSETA), uintptr(unsafe.Pointer(state)), 0, 0, 0)
}
