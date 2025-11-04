//go:build !darwin

package display

import "errors"

// RawState is a placeholder on non-darwin platforms
type RawState struct{}

func isATTY(fd uintptr) bool { return false }

func enableRawMode(fd int) (*RawState, error) {
	return nil, errors.New("raw mode unsupported on this platform")
}

func restoreRawMode(fd int, _ *RawState) {}
