package mousectl

import (
	"mousk/common/logger"
	"mousk/infra/base"
	"unsafe"
)

var (
	setCursorPos = base.User32.NewProc("SetCursorPos")
	sendInput    = base.User32.NewProc("SendInput")

	getCursorPos = base.User32.NewProc("GetCursorPos")
)

type MoveSpeedType int
type MoveDirection string

const (
	DirectionUp    MoveDirection = "up"
	DirectionDown  MoveDirection = "down"
	DirectionLeft  MoveDirection = "left"
	DirectionRight MoveDirection = "right"

	SpeedFast MoveSpeedType = 2
	SpeedSlow MoveSpeedType = 1
)

const (
	inputMouse          = 0
	mouseMove           = 0x0001
	mouseEventLeftDown  = 0x0002
	mouseEventLeftUp    = 0x0004
	mouseEventRightDown = 0x0008
	mouseEventRightUp   = 0x0010
	mouseEventFV        = 0x8000
	absolute            = 0x8000
	inputMouseFV        = 0
	move                = 1
	mouseEventF         = 0
)

type input struct {
	Type uint32
	Mi   mouseInput
}

type mouseInput struct {
	Dx          int32
	Dy          int32
	MouseData   uint32
	DwFlags     uint32
	Time        uint32
	DwExtraInfo uintptr
}

// 设置鼠标位置
func SetMousePos(x, y int) {
	setCursorPos.Call(uintptr(x), uintptr(y))
}

func MoveMouse(dx, dy int32) {
	var inputs [1]input
	inputs[0] = input{
		Type: inputMouse,
		Mi: mouseInput{
			Dx:      dx,
			Dy:      dy,
			DwFlags: mouseMove,
			Time:    0,
		},
	}
	// sendInput.Call(
	// 	uintptr(unsafe.Sizeof(input{})),
	// 	uintptr(unsafe.Pointer(&inputs[0])),
	// 	uintptr(unsafe.Sizeof(input{})),
	// )
	ret, _, err := sendInput.Call(
		uintptr(1),
		uintptr(unsafe.Pointer(&inputs[0])),
		unsafe.Sizeof(inputs[0]),
	)

	if ret == 0 {
		logger.Infof("", "SendInput failed: %v", err)
	}
}

func MoveMouseCtrl(direction MoveDirection, speedType MoveSpeedType) {
	GetMousePos()
	dx := 0
	dy := 0
	speed := int(speedType) * base.GetMoveSpeed()
	logger.Infof("", "mouse move speedType:%d, speedLevel:%d, direction:%s", speedType, base.GetMoveSpeed(), direction)
	switch direction {
	case DirectionUp:
		dx, dy = 0, -1*speed
	case DirectionDown:
		dx, dy = 0, 1*speed
	case DirectionLeft:
		dx, dy = -1*speed, 0
	case DirectionRight:
		dx, dy = 1*speed, 0
	default:
		logger.Infof("", "mouse move move direction undefined:%s", direction)
	}
	MoveMouse(int32(dx), int32(dy))
}

var moveCount = 0

func GetMousePos() (int32, int32) {
	var pos struct {
		x, y int32
	}
	getCursorPos.Call(uintptr(unsafe.Pointer(&pos)))
	moveCount += 1
	logger.Infof("", "cursor position, x:%d, y:%d, move count:%d", pos.x, pos.y, moveCount)
	return pos.x, pos.y
}
