//go:build unix

package input

import (
	"os"
	"time"

	"golang.org/x/sys/unix"
)

// FlushStdin discards any pending input in the terminal's input buffer.
// This should be called after using raw mode or direct fd reads to ensure
// no leftover keypresses interfere with subsequent line-based input.
func FlushStdin() {
	fd := int(os.Stdin.Fd())

	// Method 1: Use tcflush via ioctl to discard data received but not yet read
	// Using Syscall directly for reliability
	unix.Syscall(unix.SYS_IOCTL, uintptr(fd), uintptr(unix.TCFLSH), uintptr(unix.TCIFLUSH))

	// Method 2: Also manually drain any remaining data in non-blocking mode
	// Set non-blocking temporarily
	unix.SetNonblock(fd, true)
	buf := make([]byte, 256)
	for {
		n, err := unix.Read(fd, buf)
		if n <= 0 || err != nil {
			break
		}
	}
	unix.SetNonblock(fd, false)

	// Small delay to let terminal fully settle after mode changes
	time.Sleep(50 * time.Millisecond)
}

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
	fd := int(os.Stdin.Fd())

	// Check if stdin is a terminal
	if _, err := unix.IoctlGetTermios(fd, ioctlReadTermios); err != nil {
		// Not a terminal, can't listen for keys
		return func() {}
	}

	restore, err := setRawMode()
	if err != nil {
		return func() {}
	}

	done := make(chan struct{})

	go func() {
		buf := make([]byte, 1)

		// Set non-blocking mode for the duration of listening
		unix.SetNonblock(fd, true)
		defer unix.SetNonblock(fd, false)

		for {
			select {
			case <-done:
				return
			default:
				// Use unix.Read directly on file descriptor to avoid
				// conflicts with bufio.Reader wrapping os.Stdin
				n, err := unix.Read(fd, buf)
				if err == unix.EAGAIN || err == unix.EWOULDBLOCK {
					// No data available, sleep briefly
					time.Sleep(30 * time.Millisecond)
					continue
				}
				if n > 0 {
					select {
					case keyChan <- string(buf[0]):
					default:
					}
				} else {
					// Sleep to avoid busy loop
					time.Sleep(30 * time.Millisecond)
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
