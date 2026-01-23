//go:build unix

package input

import (
	"os"
	"syscall"
	"unsafe"
)

// termios represents the terminal settings
type termios struct {
	Iflag  uint32
	Oflag  uint32
	Cflag  uint32
	Lflag  uint32
	Cc     [20]byte
	Ispeed uint32
	Ospeed uint32
}

// getTermios gets the current terminal settings
func getTermios(fd uintptr) (*termios, error) {
	var t termios
	_, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		fd,
		syscall.TCGETS,
		uintptr(unsafe.Pointer(&t)),
	)
	if errno != 0 {
		return nil, errno
	}
	return &t, nil
}

// setTermios sets the terminal settings
func setTermios(fd uintptr, t *termios) error {
	_, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		fd,
		syscall.TCSETS,
		uintptr(unsafe.Pointer(t)),
	)
	if errno != 0 {
		return errno
	}
	return nil
}

// setRawMode sets the terminal to raw mode and returns a restore function
func setRawMode() (restore func(), err error) {
	fd := os.Stdin.Fd()

	oldTermios, err := getTermios(fd)
	if err != nil {
		return nil, err
	}

	newTermios := *oldTermios
	// Disable canonical mode and echo
	newTermios.Lflag &^= syscall.ICANON | syscall.ECHO

	if err := setTermios(fd, &newTermios); err != nil {
		return nil, err
	}

	return func() {
		setTermios(fd, oldTermios)
	}, nil
}

// ReadKeyNonBlocking attempts to read a key without blocking
// Returns the key pressed or empty string if no key available
func ReadKeyNonBlocking() string {
	// Set stdin to non-blocking
	fd := int(os.Stdin.Fd())

	// Save current flags
	flags, err := syscall.Fcntl(fd, syscall.F_GETFL, 0)
	if err != nil {
		return ""
	}

	// Set non-blocking
	syscall.Fcntl(fd, syscall.F_SETFL, flags|syscall.O_NONBLOCK)
	defer syscall.Fcntl(fd, syscall.F_SETFL, flags)

	buf := make([]byte, 1)
	n, _ := os.Stdin.Read(buf)
	if n > 0 {
		return string(buf[0])
	}
	return ""
}

// StartKeyListener starts listening for key presses and sends them to the channel
// Call the returned stop function to clean up
func StartKeyListener(keyChan chan<- string) (stop func()) {
	restore, err := setRawMode()
	if err != nil {
		return func() {}
	}

	done := make(chan struct{})

	go func() {
		buf := make([]byte, 1)
		for {
			select {
			case <-done:
				return
			default:
				// Set non-blocking read with timeout
				fd := int(os.Stdin.Fd())
				flags, _ := syscall.Fcntl(fd, syscall.F_GETFL, 0)
				syscall.Fcntl(fd, syscall.F_SETFL, flags|syscall.O_NONBLOCK)

				n, _ := os.Stdin.Read(buf)

				syscall.Fcntl(fd, syscall.F_SETFL, flags)

				if n > 0 {
					select {
					case keyChan <- string(buf[0]):
					default:
					}
				}
			}
		}
	}()

	return func() {
		close(done)
		if restore != nil {
			restore()
		}
	}
}
