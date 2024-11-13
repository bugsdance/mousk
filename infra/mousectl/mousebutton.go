package mousectl

import (
	"mousk/common/logger"
	"time"
	"unsafe"
)

func LeftClick() {
	var inputs [2]input

	inputs[0] = input{
		Type: inputMouse,
		Mi: mouseInput{
			DwFlags: mouseEventLeftDown,
		},
	}

	inputs[1] = input{
		Type: inputMouse,
		Mi: mouseInput{
			DwFlags: mouseEventLeftUp,
		},
	}

	ret, _, err := sendInput.Call(
		uintptr(len(inputs)),
		uintptr(unsafe.Pointer(&inputs[0])),
		unsafe.Sizeof(inputs[0]),
	)

	if ret == 0 {
		logger.Infof("", "SendInput failed: %v", err)
	}
}

func RightClick() {
	var inputs [2]input

	inputs[0] = input{
		Type: inputMouse,
		Mi: mouseInput{
			DwFlags: mouseEventRightDown,
		},
	}

	inputs[1] = input{
		Type: inputMouse,
		Mi: mouseInput{
			DwFlags: mouseEventRightUp,
		},
	}

	ret, _, err := sendInput.Call(
		uintptr(len(inputs)),
		uintptr(unsafe.Pointer(&inputs[0])),
		unsafe.Sizeof(inputs[0]),
	)

	if ret == 0 {
		logger.Infof("", "SendInput failed: %v", err)
	}
}

func LeftClickLongPress(duration time.Duration) {
	var inputs [1]input

	inputs[0] = input{
		Type: inputMouse,
		Mi: mouseInput{
			DwFlags: mouseEventLeftDown,
		},
	}

	ret, _, err := sendInput.Call(
		uintptr(1),
		uintptr(unsafe.Pointer(&inputs[0])),
		unsafe.Sizeof(inputs[0]),
	)

	if ret == 0 {
		logger.Infof("", "SendInput failed: %v", err)
	}

	time.Sleep(duration)

	inputs[0].Mi.DwFlags = mouseEventLeftUp

	ret, _, err = sendInput.Call(
		uintptr(1),
		uintptr(unsafe.Pointer(&inputs[0])),
		unsafe.Sizeof(inputs[0]),
	)

	if ret == 0 {
		logger.Infof("", "SendInput failed: %v", err)
	}
}

func RightClickLongPress(duration time.Duration) {
	var inputs [1]input
	inputs[0] = input{
		Type: inputMouse,
		Mi: mouseInput{
			DwFlags: mouseEventRightDown,
		},
	}

	ret, _, err := sendInput.Call(
		uintptr(1),
		uintptr(unsafe.Pointer(&inputs[0])),
		unsafe.Sizeof(inputs[0]),
	)

	if ret == 0 {
		logger.Infof("", "SendInput failed: %v", err)
	}

	time.Sleep(duration)

	inputs[0].Mi.DwFlags = mouseEventRightUp

	ret, _, err = sendInput.Call(
		uintptr(1),
		uintptr(unsafe.Pointer(&inputs[0])),
		unsafe.Sizeof(inputs[0]),
	)

	if ret == 0 {
		logger.Infof("", "SendInput failed: %v", err)
	}
}
