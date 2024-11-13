package monitor

import (
	"mousk/common/logger"
	"mousk/infra/base"
	"syscall"
	"unsafe"
)

var (
	enumDisplayMonitors = base.User32.NewProc("EnumDisplayMonitors")
	getMonitorInfoW     = base.User32.NewProc("GetMonitorInfoW")
)

const (
	SM_CXSCREEN = 0
	SM_CYSCREEN = 1
)

type MONITORINFOEX struct {
	Size    uint32
	Monitor RECT
	Work    RECT
	Flags   uint32
}

type RECT struct {
	Left   int32
	Top    int32
	Right  int32
	Bottom int32
}

func GetMonitors() []MONITORINFOEX {
	// 显示器信息获取
	var monitors []MONITORINFOEX

	// 定义回调函数
	callback := syscall.NewCallback(func(hMonitor uintptr, hdc uintptr, lprc uintptr, dwData uintptr) uintptr {
		var mi MONITORINFOEX
		mi.Size = uint32(unsafe.Sizeof(mi))
		if ret, _, _ := enumDisplayMonitors.Call(hdc, uintptr(unsafe.Pointer(&RECT{})), syscall.NewCallback(func(lprc uintptr, hdc uintptr, lprcClip uintptr, dwData uintptr) uintptr {
			return 1
		}), 0); ret != 0 {
			if ret, _, _ := syscall.Syscall(getMonitorInfoW.Addr(), 2, hMonitor, uintptr(unsafe.Pointer(&mi)), 0); ret != 0 {
				monitors = append(monitors, mi)
			}
		}
		return 1
	})

	// 调用 EnumDisplayMonitors 函数
	if ret, _, _ := enumDisplayMonitors.Call(0, 0, callback, 0); ret == 0 {
		logger.Infof("", "EnumDisplayMonitors failed")
		return nil
	}

	// 输出显示器数量和范围信息
	logger.Infof("", "Number of monitors: %d", len(monitors))
	for i, monitor := range monitors {
		logger.Infof("", "Monitor %d:", i+1)
		logger.Infof("", "    Left:   %d", monitor.Monitor.Left)
		logger.Infof("", "    Top:    %d", monitor.Monitor.Top)
		logger.Infof("", "    Right:  %d", monitor.Monitor.Right)
		logger.Infof("", "    Bottom: %d", monitor.Monitor.Bottom)
	}
	return monitors

}
