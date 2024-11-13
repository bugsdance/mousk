package mousectl

import (
	"mousk/common/logger"
	"syscall"
	"unsafe"
)

// 定义常量
const (
	INPUT_MOUSE          = 0
	MOUSEEVENTF_LEFTDOWN = 0x0002
	MOUSEEVENTF_LEFTUP   = 0x0004
)

func MouseLeftUp() {
	sendMouseInput(MOUSEEVENTF_LEFTUP)
}

func MouseLeftDown() {
	sendMouseInput(MOUSEEVENTF_LEFTDOWN)
}

func sendMouseInput(flags uint32) {
	var input input
	input.Type = INPUT_MOUSE
	input.Mi.DwFlags = flags

	_, _, err := sendInput.Call(
		uintptr(1),
		uintptr(unsafe.Pointer(&input)),
		unsafe.Sizeof(input),
	)
	if err != syscall.Errno(0) {
		logger.Infof("", "Error calling SendInput: %v", err)
	}
}
