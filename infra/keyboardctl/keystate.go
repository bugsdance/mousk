package keyboardctl

import (
	"fmt"
	"mousk/infra/base"
	"time"
)

var (
	procGetAsyncKeyState = base.User32.NewProc("GetAsyncKeyState")
)

func IsShiftPressed() bool {
	v, _, _ := procGetAsyncKeyState.Call(uintptr(VK_SHIFT))
	return v != 0
}

func main() {
	for {
		if IsShiftPressed() {
			fmt.Println("Shift key is pressed")
		} else {
			fmt.Println("Shift key is not pressed")
		}
		time.Sleep(100 * time.Millisecond) // Check every 100 milliseconds
	}
}
