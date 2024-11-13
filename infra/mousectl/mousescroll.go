package mousectl

import (
	"mousk/common/logger"
	"mousk/infra/base"
)

type ScrollDirection string

const (
	DirectionHorizontalRight ScrollDirection = "horizontal_right"
	DirectionHorizontalLeft  ScrollDirection = "horizontal_left"
	DirectionVerticalUp      ScrollDirection = "vertical_up"
	DirectionVerticalDown    ScrollDirection = "vertical_down"
)

const (
	MOUSEEVENTF_WHEEL  = 0x0800
	MOUSEEVENTF_HWHEEL = 0x1000
	WHEEL_DELTA        = 120
	// WHEEL_DELTA = 10
)

var (
	procMouseEvent = base.User32.NewProc("mouse_event")
)

func mouse_event(dwFlags, dx, dy, dwData, dwExtraInfo uintptr) {
	procMouseEvent.Call(dwFlags, dx, dy, dwData, dwExtraInfo)
}

// 模拟垂直滚动
func ScrollVertically(amount int) {
	mouse_event(MOUSEEVENTF_WHEEL, 0, 0, uintptr(amount*WHEEL_DELTA), 0)
}

// 模拟水平滚动
func ScrollHorizontally(amount int) {
	mouse_event(MOUSEEVENTF_HWHEEL, 0, 0, uintptr(amount*WHEEL_DELTA), 0)
}

func ScrollMouseCtrl(direction ScrollDirection, speedType MoveSpeedType) {
	speed := int(speedType) * base.GetScrollSpeed()
	logger.Infof("", "mouse scroll speedType:%d, speedLevel:%d, direction:%s", speedType, base.GetScrollSpeed(), direction)
	switch direction {
	case DirectionHorizontalLeft:
		ScrollHorizontally(-1 * speed)
	case DirectionHorizontalRight:
		ScrollHorizontally(1 * speed)
	case DirectionVerticalDown:
		ScrollVertically(-1 * speed)
	case DirectionVerticalUp:
		ScrollVertically(1 * speed)
	default:
		logger.Infof("", "mouse scroll scroll direction undefined:%s", direction)
	}
}
