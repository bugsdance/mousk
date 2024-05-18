package main

import (
	"fmt"
	"mousek/infra/keyboardctl"
	"mousek/infra/monitor"
	"mousek/infra/mousectl"
	"mousek/infra/util"
	"os"
	"time"
	"unsafe"
)

var mode = 0       // 0:normal, 1:control
var speedLevel = 1 // the speed of you mouse movement
var vkCodesMulitiSpeedLevelArr = []uint32{keyboardctl.VK_1, keyboardctl.VK_2, keyboardctl.VK_3, keyboardctl.VK_4, keyboardctl.VK_5}

const (
	ModeNormal  = 0
	ModeControl = 1
)

func main() {

	// monitors := monitor.GetMonitors()
	// for _, monitor := range monitors {
	// 	moveMouseAround(monitor.Monitor)
	// }
	// win+space : activate control mode
	vkCodesWinSpace := []uint32{keyboardctl.VK_LWIN, keyboardctl.VK_SPACE}
	startControlMode := func(wParam uintptr, vkCode, scanCode uint32) uintptr {
		fmt.Printf("current mode:%d", mode)
		fmt.Println()
		if mode == ModeControl {
			fmt.Println("already in control mode", time.Now())
		} else {
			mode = ModeControl
			fmt.Println("change to control mode", time.Now())
		}
		return 0
	}
	keyboardctl.RegisterOne(startControlMode, vkCodesWinSpace...)

	// when in ModeControl, 1\2\3\4...,control the speed of your mouse move
	vkCodesMulitiSpeedLevel := [][]uint32{{keyboardctl.VK_1}, {keyboardctl.VK_2}, {keyboardctl.VK_3}, {keyboardctl.VK_4}, {keyboardctl.VK_5}}
	speedLevelSwitch := func(wParam uintptr, vkCode, scanCode uint32) uintptr {
		fmt.Printf("current mode:%d,current speed:%d\n", mode, speedLevel)

		if mode != ModeControl {
			fmt.Printf("not in control mode, can not switch speed,mode:%d,current speed:%d\n", mode, speedLevel)
			return 0
		}
		if util.Contains[uint32](vkCodesMulitiSpeedLevelArr, uint32(vkCode)) {
			speedLevel = int(vkCode) - keyboardctl.VK_0
			fmt.Printf("change speed to :%d\n", speedLevel)
		}
		return 0
	}
	keyboardctl.RegisterMulti(speedLevelSwitch, vkCodesMulitiSpeedLevel...)

	keyboardctl.RawKeyboardListener(keyboardctl.LowLevelKeyboardCallback)

}

// 控制鼠标在指定显示器的四周移动
func moveMouseAround(monitor monitor.RECT) {
	x := int(monitor.Left)
	y := int(monitor.Top)

	/* 	width := int(monitor.Right - monitor.Left)
	   	height := int(monitor.Bottom - monitor.Top)
	*/
	// 向右移动到显示器右边缘
	for x < int(monitor.Right) {
		mousectl.MoveMouse(x, y)
		x += 10
		time.Sleep(5 * time.Millisecond)
	}

	// 向下移动到显示器底边缘
	for y < int(monitor.Bottom) {
		mousectl.MoveMouse(x, y)
		y += 10
		time.Sleep(5 * time.Millisecond)
	}

	// 向左移动到显示器左边缘
	for x > int(monitor.Left) {
		mousectl.MoveMouse(x, y)
		x -= 10
		time.Sleep(5 * time.Millisecond)
	}

	// 向上移动到显示器上边缘
	for y > int(monitor.Top) {
		mousectl.MoveMouse(x, y)
		y -= 10
		time.Sleep(5 * time.Millisecond)
	}
}

func Callback(nCode int, wParam uintptr, lParam uintptr) uintptr {
	if nCode >= 0 {
		kbdStruct := (*keyboardctl.KBDLLHOOKSTRUCT)(unsafe.Pointer(lParam))
		vkCode := kbdStruct.VkCode

		if wParam == keyboardctl.WM_KEYDOWN {
			keyboardctl.SetPressed(vkCode)
			fmt.Printf("Key pressed (VK code): %x\n", vkCode)
		} else if wParam == keyboardctl.WM_KEYUP {
			keyboardctl.SetReleased(vkCode)
			fmt.Printf("Key released (VK code): %x\n", vkCode)
		}

		// 检查是否同时按下了 Ctrl、Shift 和 A 键
		// if keyboardctl.Pressed(keyboardctl._VK_CTRL) && keyboardctl.Pressed(keyboardctl._VK_SHIFT) && keyboardctl.Pressed(keyboardctl.VK_A) {
		// 	fmt.Println("Ctrl+Shift+A keys pressed simultaneously")
		// }

		// 如果按下了 'Q' 键，退出程序
		if keyboardctl.Pressed(keyboardctl.VK_Q) {
			os.Exit(0)
		}
		return 1
	}
	return 0
}
