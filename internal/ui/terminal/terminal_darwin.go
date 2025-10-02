package terminal

import (
	"os"
	"syscall"
	"unsafe"
)

type termios struct {
	Iflag  uint64
	Oflag  uint64
	Cflag  uint64
	Lflag  uint64
	Cc     [20]byte
	Ispeed uint64
	Ospeed uint64
}

var oldTermState *termios

func EnableRawMode() error {
	var oldState termios

	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(syscall.TIOCGETA),
		uintptr(unsafe.Pointer(&oldState)))

	if errno != 0 {
		return errno
	}

	newState := oldState
	newState.Lflag &^= syscall.ECHO | syscall.ICANON
	newState.Cc[syscall.VMIN] = 1
	newState.Cc[syscall.VTIME] = 0

	_, _, errno = syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(syscall.TIOCSETA),
		uintptr(unsafe.Pointer(&newState)))

	if errno != 0 {
		return errno
	}

	oldTermState = &oldState
	return nil
}

func DisableRawMode() error {
	if oldTermState == nil {
		return nil
	}

	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(syscall.TIOCSETA),
		uintptr(unsafe.Pointer(oldTermState)))

	if errno != 0 {
		return errno
	}
	return nil
}

func ReadInput() ([]byte, error) {
	buf := make([]byte, 32)
	n, err := os.Stdin.Read(buf)
	if err != nil {
		return nil, err
	}
	return buf[:n], nil
}
