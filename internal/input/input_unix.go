//go:build unix

package input

import (
	"os"
	"time"

	"golang.org/x/sys/unix"
)

// setRawMode sets the terminal to raw mode and returns a restore function
func setRawMode() (restore func(), err error) {
	fd := int(os.Stdin.Fd())

	oldTermios, err := unix.IoctlGetTermios(fd, ioctlReadTermios)
	if err != nil {
		return nil, err
	}

	newTermios := *oldTermios
	// Disable canonical mode and echo
	newTermios.Lflag &^= unix.ICANON | unix.ECHO
	// Set minimum bytes and timeout for read
	newTermios.Cc[unix.VMIN] = 0
	newTermios.Cc[unix.VTIME] = 0

	if err := unix.IoctlSetTermios(fd, ioctlWriteTermios, &newTermios); err != nil {
		return nil, err
	}

	return func() {
		unix.IoctlSetTermios(fd, ioctlWriteTermios, oldTermios)
	}, nil
}

// ReadKeyNonBlocking attempts to read a key without blocking
// Returns the key pressed or empty string if no key available
func ReadKeyNonBlocking() string {
	fd := int(os.Stdin.Fd())

	// Set non-blocking
	if err := unix.SetNonblock(fd, true); err != nil {
		return ""
	}
	defer unix.SetNonblock(fd, false)

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
		fd := int(os.Stdin.Fd())
		buf := make([]byte, 1)

		for {
			select {
			case <-done:
				return
			default:
				// Set non-blocking temporarily
				unix.SetNonblock(fd, true)
				n, _ := os.Stdin.Read(buf)
				unix.SetNonblock(fd, false)

				if n > 0 {
					select {
					case keyChan <- string(buf[0]):
					default:
					}
				} else {
					// Small sleep to avoid busy loop
					time.Sleep(50 * time.Millisecond)
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
