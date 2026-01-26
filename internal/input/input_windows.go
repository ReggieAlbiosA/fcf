//go:build windows

package input

import (
	"syscall"
	"unsafe"
)

var (
	kernel32                          = syscall.NewLazyDLL("kernel32.dll")
	procGetStdHandle                  = kernel32.NewProc("GetStdHandle")
	procReadConsoleInput              = kernel32.NewProc("ReadConsoleInputW")
	procGetNumberOfConsoleInputEvents = kernel32.NewProc("GetNumberOfConsoleInputEvents")
	procSetConsoleMode                = kernel32.NewProc("SetConsoleMode")
	procGetConsoleMode                = kernel32.NewProc("GetConsoleMode")
	procFlushConsoleInputBuffer       = kernel32.NewProc("FlushConsoleInputBuffer")
)

const (
	stdInputHandle = ^uintptr(0) - 10 + 1 // STD_INPUT_HANDLE = -10
	enableEchoInput = 0x0004
	enableLineInput = 0x0002
	enableProcessedInput = 0x0001
	keyEvent = 0x0001
)

type inputRecord struct {
	EventType uint16
	_         uint16
	Event     [16]byte
}

type keyEventRecord struct {
	KeyDown         int32
	RepeatCount     uint16
	VirtualKeyCode  uint16
	VirtualScanCode uint16
	Char            uint16
	ControlKeyState uint32
}

func getStdHandle(handle uintptr) uintptr {
	ret, _, _ := procGetStdHandle.Call(handle)
	return ret
}

// FlushStdin discards any pending input in the console's input buffer.
// This should be called after using raw mode to ensure no leftover keypresses
// interfere with subsequent line-based input.
func FlushStdin() {
	handle := getStdHandle(stdInputHandle)
	if handle != 0 {
		procFlushConsoleInputBuffer.Call(handle)
	}
}

// ReadKeyNonBlocking attempts to read a key without blocking
func ReadKeyNonBlocking() string {
	handle := getStdHandle(stdInputHandle)
	if handle == 0 {
		return ""
	}

	// Check if there are events available
	var numEvents uint32
	ret, _, _ := procGetNumberOfConsoleInputEvents.Call(handle, uintptr(unsafe.Pointer(&numEvents)))
	if ret == 0 || numEvents == 0 {
		return ""
	}

	// Read the input
	var ir inputRecord
	var numRead uint32
	ret, _, _ = procReadConsoleInput.Call(
		handle,
		uintptr(unsafe.Pointer(&ir)),
		1,
		uintptr(unsafe.Pointer(&numRead)),
	)

	if ret == 0 || numRead == 0 {
		return ""
	}

	if ir.EventType == keyEvent {
		keyEvent := (*keyEventRecord)(unsafe.Pointer(&ir.Event[0]))
		if keyEvent.KeyDown != 0 && keyEvent.Char != 0 {
			return string(rune(keyEvent.Char))
		}
	}

	return ""
}

// StartKeyListener starts listening for key presses and sends them to the channel
func StartKeyListener(keyChan chan<- string) (stop func()) {
	handle := getStdHandle(stdInputHandle)

	// Save old console mode
	var oldMode uint32
	procGetConsoleMode.Call(handle, uintptr(unsafe.Pointer(&oldMode)))

	// Set new mode without line input and echo
	newMode := oldMode &^ (enableLineInput | enableEchoInput)
	procSetConsoleMode.Call(handle, uintptr(newMode))

	done := make(chan struct{})

	go func() {
		for {
			select {
			case <-done:
				return
			default:
				key := ReadKeyNonBlocking()
				if key != "" {
					select {
					case keyChan <- key:
					default:
					}
				}
			}
		}
	}()

	return func() {
		close(done)
		// Restore old console mode
		procSetConsoleMode.Call(handle, uintptr(oldMode))
	}
}
