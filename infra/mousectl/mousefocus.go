package mousectl

import (
	"mousk/infra/base"
	"syscall"
	"unicode/utf16"
	"unsafe"
)

var (
	procGetForegroundWindow = base.User32.NewProc("GetForegroundWindow")
	procGetFocus            = base.User32.NewProc("GetFocus")
	procGetClassNameW       = base.User32.NewProc("GetClassNameW")
	procGetGUIThreadInfo    = base.User32.NewProc("GetGUIThreadInfo")
	procGetCaretPos         = base.User32.NewProc("GetCaretPos")
)

// IsCaretInInputState checks if the caret (text cursor) is in an input state
// by using either GetGUIThreadInfo or GetCaretPos.
func IsCaretInInputState() bool {
	var guiThreadInfo GUITHREADINFO
	guiThreadInfo.cbSize = uint32(unsafe.Sizeof(guiThreadInfo))

	// Using the current thread's ID (0 for the current thread)
	threadID := uint32(0)

	// Try to get GUI thread information
	guiThreadSuccessRaw, _, _ := procGetGUIThreadInfo.Call(
		uintptr(threadID),
		uintptr(unsafe.Pointer(&guiThreadInfo)),
	)
	guiThreadSuccess := guiThreadSuccessRaw != 0
	// Check if the GUI thread information indicates the caret is visible
	if guiThreadSuccess && guiThreadInfo.hwndCaret != 0 {
		return true
	}

	// Get the caret position
	var caretPos POINT
	caretPosSuccessRaw, _, _ := procGetCaretPos.Call(uintptr(unsafe.Pointer(&caretPos)))
	caretPosSuccess := caretPosSuccessRaw != 0
	// Return true if either GUI thread info was successful or caret position was obtained
	return caretPosSuccess
}

type POINT struct {
	X int32
	Y int32
}

type GUITHREADINFO struct {
	cbSize        uint32
	flags         uint32
	hwndActive    syscall.Handle
	hwndFocus     syscall.Handle
	hwndCapture   syscall.Handle
	hwndMenuOwner syscall.Handle
	hwndMoveSize  syscall.Handle
	hwndCaret     syscall.Handle
	rcCaret       RECT
}

type RECT struct {
	Left   int32
	Top    int32
	Right  int32
	Bottom int32
}

func GetGUIThreadInfo(idThread uint32, lpgui *GUITHREADINFO) bool {
	ret, _, _ := procGetGUIThreadInfo.Call(
		uintptr(idThread),
		uintptr(unsafe.Pointer(lpgui)),
	)
	return ret != 0
}

// GetForegroundWindow retrieves the handle of the window that is currently in the foreground.
func GetForegroundWindow() uintptr {
	hwnd, _, _ := procGetForegroundWindow.Call()
	return hwnd
}

// GetFocus retrieves the handle of the window that has the keyboard focus, if the window is attached to the calling thread's message queue.
func GetFocus() uintptr {
	hwnd, _, _ := procGetFocus.Call()
	return hwnd
}

// GetClassName retrieves the class name of the specified window.
func GetClassName(hwnd uintptr) string {
	var className [256]uint16
	procGetClassNameW.Call(hwnd, uintptr(unsafe.Pointer(&className[0])), uintptr(len(className)))

	// Convert UTF-16 to string
	return string(utf16.Decode(className[:]))
}

// IsFocusInEditControl checks if the current focus is in an Edit control.
func IsFocusInEditControl() bool {
	hwndForeground := GetForegroundWindow()
	hwndFocus := GetFocus()

	if hwndFocus != 0 && hwndFocus == hwndForeground {
		className := GetClassName(hwndFocus)
		return className == "Edit"
	}

	return false
}
