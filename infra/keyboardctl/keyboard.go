package keyboardctl

import (
	"encoding/json"
	"mousk/common/logger"
	"mousk/infra/base"
	"os"
	"syscall"
	"time"
	"unsafe"
)

var (
	keyPressedStates    map[uint32]*KeyState
	setWindowsHookEx    = base.User32.NewProc("SetWindowsHookExW")
	getMessageW         = base.User32.NewProc("GetMessageW")
	unhookWindowsHookEx = base.User32.NewProc("UnhookWindowsHookEx")
	// todo serilize to file config
	listeningKeyReference = make(map[uint32]*KeyReference, 0)
)

type KeyState struct {
	Pressed            bool
	LastReleased       *time.Time
	LastPressed        *time.Time
	SecondLastReleased *time.Time
	SecondLastPressed  *time.Time
}

type KBDLLHOOKSTRUCT struct {
	VkCode   uint32
	ScanCode uint32
}
type KeyReference struct {
	Count           int
	KeyCombinations []KeyCallback
}

type KeyCallback struct {
	LastTriggerTime    time.Time
	FirstClickKeys     []uint32  // for example. []{"ctrl","a"}
	SecondClickKeys    []uint32  // for example. []{"ctrl","a"}
	Cb                 Callback2 `json:"-"`
	CbPriority         int
	withReleaseEvent   bool
	effectOnNormalMode bool
}

type Callback HookProc

// type Callback2 HookProc
type Callback2 func(wParam uintptr, vkCode, scanCode uint32) uintptr
type HookProc func(nCode int, wParam uintptr, lParam uintptr) uintptr

func registerKeyListening(cb Callback2, cbPriority int, effectOnNormal bool, withReleaseEvent bool, firstClickVkCodes []uint32, secondClickVkCodes []uint32) {
	for _, v := range firstClickVkCodes {
		if listeningKeyReference[v] == nil {
			listeningKeyReference[v] = &KeyReference{}
		}
		listeningKeyReference[v].Count += 1
		listeningKeyReference[v].KeyCombinations = append(listeningKeyReference[v].KeyCombinations, KeyCallback{
			FirstClickKeys:     firstClickVkCodes,
			SecondClickKeys:    secondClickVkCodes,
			Cb:                 cb,
			CbPriority:         cbPriority,
			withReleaseEvent:   withReleaseEvent,
			effectOnNormalMode: effectOnNormal,
		})
	}
	for _, v := range secondClickVkCodes {
		if listeningKeyReference[v] == nil {
			listeningKeyReference[v] = &KeyReference{}
		}
		listeningKeyReference[v].Count += 1
		listeningKeyReference[v].KeyCombinations = append(listeningKeyReference[v].KeyCombinations, KeyCallback{
			FirstClickKeys:     firstClickVkCodes,
			SecondClickKeys:    secondClickVkCodes,
			Cb:                 cb,
			withReleaseEvent:   withReleaseEvent,
			effectOnNormalMode: effectOnNormal,
		})
	}
}

//lint:ignore U1000 Ignore unused function temporarily for debugging
func unRegisterKeyListening(vkCodes ...uint32) {
	for _, v := range vkCodes {
		listeningKeyReference[v].Count -= 1
		listeningKeyReference[v].KeyCombinations = append(listeningKeyReference[v].KeyCombinations, KeyCallback{
			FirstClickKeys: vkCodes,
		})
	}

}

func RegisterMulti(cb Callback2, priority int, mulitiVkCodes ...[]uint32) {
	// switch keyAction {
	// case WM_KEYDOWN:

	// 	break
	// default:
	// 	logger.Infof("","method not develop yet")
	// 	return
	// }
	for _, vkCodes := range mulitiVkCodes {
		registerKeyListening(cb, priority, false, false, vkCodes, nil)
	}
}

func RegisterNormal(cb Callback2, priority int, vkCodes ...uint32) {
	registerKeyListening(cb, priority, true, false, vkCodes, nil)
}

func RegisterOne(cb Callback2, priority int, vkCodes ...uint32) {
	registerKeyListening(cb, priority, false, false, vkCodes, nil)
}

func RegisterDoubleClick(cb Callback2, priority int, firstClick []uint32, secondClick []uint32) {
	registerKeyListening(cb, priority, false, false, firstClick, secondClick)
}

func EffectOnNormalMode(vkCode uint32) bool {
	if listeningKeyReference[vkCode] != nil {
		// logger.Infof("", "effect on normal check:%+v, current:%d", listeningKeyReference[vkCode], vkCode)
		ref := listeningKeyReference[vkCode]
		for _, v := range ref.KeyCombinations {
			// logger.Infof("", "effect on normal check:%+v, %t", v.FirstClickKeys, StatusCheck(v.FirstClickKeys, 1, time.Second))
			if !v.effectOnNormalMode {
				continue
			}
			// logger.Infof("","11111:%+v, %t", v.FirstClickKeys, AllPressed(v.FirstClickKeys...))
			// if AllPressed(v.FirstClickKeys...) {
			if StatusCheckNew(v.FirstClickKeys, 1) {
				return true
			}
		}
	}
	// if AllPressed(VK_LWIN, VK_SPACE) {
	// 	return true
	// }
	// if AllPressed(VK_LALT, VK_0) {
	// 	return true
	// }
	// if AllPressed(VK_TAB, VK_SPACE) {
	// 	return true
	// }
	return false
}

func LowLevelKeyboardCallback(nCode int, wParam uintptr, lParam uintptr) uintptr {
	if nCode == 0 {
		kbdStruct := (*KBDLLHOOKSTRUCT)(unsafe.Pointer(lParam))
		scanCode := kbdStruct.ScanCode
		vkCode := kbdStruct.VkCode

		if wParam == WM_KEYDOWN || wParam == WM_SYSKEYDOWN {
			SetPressed(vkCode)
		} else if wParam == WM_KEYUP || wParam == WM_SYSKEYUP {
			SetReleased(vkCode)
			if vkCode == 160 {
				logger.Infof("", "shift dectect:%t", IsShiftPressed())
			}
		}

		// // 检查是否同时按下了 Ctrl、Shift 和 A 键
		// if Pressed(VK_LCONTROL) && Pressed(VK_LSHIFT) && Pressed(VK_A) {
		// 	logger.Infof("", "Ctrl+Shift+A keys pressed simultaneously")
		// 	os.Exit(0)
		// 	return 1
		// }

		if base.GetMode() == base.ModeNormal && !EffectOnNormalMode(vkCode) {
			logger.Infof("", "%d not in control mode, mode:%d, keystatus:%d", time.Now().UnixMilli(), base.GetMode(), vkCode)
			return 0
		}

		if listeningKeyReference[vkCode] != nil {
			ref := listeningKeyReference[vkCode]
			// TODO : ref.KeyCombinations中满足多个快捷键组合时，仅优先执行其中按键数量最多的
			// k | ctrl+k, 优先执行ctrl+k

			satisfiedCallback := make([]KeyCallback, 0)
			for _, v := range ref.KeyCombinations {
				// if AllPressed(v.FirstClickKeys...) {
				if StatusCheckNew(v.FirstClickKeys, 1) {
					satisfiedCallback = append(satisfiedCallback, v)
				}
				if v.withReleaseEvent && StatusCheckNew(v.FirstClickKeys, 0) {
					satisfiedCallback = append(satisfiedCallback, v)
				}
			}

			// return 1 is importantA
			// 如果没有匹配的快捷键，返回1
			if len(satisfiedCallback) == 0 {
				return 0
			}

			// order by priority
			mostKeyNumCallback := satisfiedCallback[0]

			for _, v := range satisfiedCallback {
				logger.Infof("", "all keycallback:%+v", GetNamesByCodes(v.FirstClickKeys))
				if len(v.FirstClickKeys) > len(mostKeyNumCallback.FirstClickKeys) {
					mostKeyNumCallback = v
				} else if len(v.FirstClickKeys) == len(mostKeyNumCallback.FirstClickKeys) {
					if v.CbPriority > mostKeyNumCallback.CbPriority {
						mostKeyNumCallback = v
					}

				} else {
				}
			}
			logger.Infof("", "most keycallback:%+v", GetNamesByCodes(mostKeyNumCallback.FirstClickKeys))
			// return mostKeyNumCallback.Cb(wParam, vkCode, scanCode)
			mostKeyNumCallback.Cb(wParam, vkCode, scanCode)
			return 1
		}
		return 0
	}
	return 0
}

func KeyboardCallback(nCode int, wParam uintptr, lParam uintptr) uintptr {
	if nCode >= 0 {
		kbdStruct := (*KBDLLHOOKSTRUCT)(unsafe.Pointer(lParam))
		vkCode := kbdStruct.VkCode

		if wParam == WM_KEYDOWN {
			SetPressed(vkCode)
			logger.Infof("", "Key pressed (VK code): %d", vkCode)
		} else if wParam == WM_KEYUP {
			SetReleased(vkCode)
			logger.Infof("", "Key released (VK code): %d", vkCode)
		}

		// 检查是否同时按下了 Ctrl、Shift 和 A 键
		if Pressed(VK_LCONTROL) && Pressed(VK_LSHIFT) && Pressed(VK_A) {
			logger.Infof("", "Ctrl+Shift+A keys pressed simultaneously")
		}

		// 如果按下了 'Q' 键，退出程序
		if Pressed(VK_Q) {
			os.Exit(0)
		}
		return 1
	}
	return 0
}

func RawKeyboardListener(cb Callback) {

	keyPressedStates = make(map[uint32]*KeyState)

	hookProc := HookProc(cb)

	hookID, _, _ := setWindowsHookEx.Call(
		uintptr(WH_KEYBOARD_LL),
		syscall.NewCallback(hookProc),
		0,
		0,
	)

	if hookID == 0 {
		logger.Infof("", "Failed to set hook")
		return
	}

	logger.Infof("", "Hook set, waiting for events...")

	defer unhookWindowsHookEx.Call(hookID)

	var msg uintptr
	for {
		_, _, _ = getMessageW.Call(uintptr(unsafe.Pointer(&msg)), 0, 0, 0)
	}

}

func Pressed(vkCode uint32) bool {
	// logger.Infof("",vkCode, keyPressedStates[vkCode])
	if keyPressedStates[vkCode] == nil {
		keyPressedStates[vkCode] = nilKeyState()
	}
	return keyPressedStates[vkCode].Pressed
}

// deprecated
func AllPressed(vkCodes ...uint32) bool {
	if vkCodes == nil {
		return true
	}
	allPressed := true
	for _, v := range vkCodes {
		logger.Infof("", "keys %d , %t", v, Pressed(v))
		allPressed = allPressed && Pressed(v)
	}
	return allPressed
}

func SetPressed(vkCode uint32) {
	if keyPressedStates[vkCode] == nil {
		keyPressedStates[vkCode] = nilKeyState()
	}
	keyPressedStates[vkCode].Pressed = true

	currTime := time.Now()
	keyPressedStates[vkCode].LastPressed = &currTime
	keyPressedStates[vkCode].SecondLastPressed = keyPressedStates[vkCode].LastPressed
	logger.Infof("", "Key pressed (VK code): %d, last pressed: %s", vkCode, currTime)
}

func nilKeyState() *KeyState {
	return &KeyState{
		Pressed:            false,
		LastReleased:       nil,
		LastPressed:        nil,
		SecondLastPressed:  nil,
		SecondLastReleased: nil,
	}
}

func SetReleased(vkCode uint32) {
	if keyPressedStates[vkCode] == nil {
		keyPressedStates[vkCode] = nilKeyState()
	}
	keyPressedStates[vkCode].Pressed = false
	currTime := time.Now()
	keyPressedStates[vkCode].SecondLastReleased = keyPressedStates[vkCode].LastReleased
	keyPressedStates[vkCode].LastReleased = &currTime
	logger.Infof("", "Key released (VK code): %d, last released: %s", vkCode, currTime)
}

func RegisterWithReleaseEventMulti(cb Callback2, prority int, mulitiVkCodes ...[]uint32) {
	// switch keyAction {
	// case WM_KEYDOWN:

	// 	break
	// default:
	// 	logger.Infof("","method not develop yet")
	// 	return
	// }
	for _, vkCodes := range mulitiVkCodes {
		registerKeyListening(cb, prority, false, true, vkCodes, nil)
	}
}

func StatusCheckNew(vkCodes []uint32, pressed int) bool {
	// logger.Infof("", "key status check param:%+v", GetNamesByCodes(vkCodes))
	if vkCodes == nil {
		return true
	}
	for _, v := range vkCodes {
		keyState, ok := keyPressedStates[v]
		// logger.Infof("", "key status check param:%+v, key:%s, keystate:%+v", GetNamesByCodes(vkCodes), GetNameByCode(v), keyState)
		if !ok {
			keyState = nilKeyState()
		}
		if pressed == 1 && !keyState.Pressed {
			return false
		}
		if pressed == 0 && keyState.Pressed {
			return false
		}
	}
	logger.Infof("", "key status check:%+v", GetNamesByCodes(vkCodes))
	return true
}

// Deprecated!!!
func StatusCheck(vkCodes []uint32, pressed int, durationBetween time.Duration) bool {
	// logger.Infof("", "key status check param:%+v", GetNamesByCodes(vkCodes))
	if vkCodes == nil {
		return true
	}

	var maxLastReleasedTime *time.Time = nil
	var minLastReleasedTime *time.Time = nil
	var maxLastPressedTime *time.Time = nil
	var minLastPressedTime *time.Time = nil
	for _, v := range vkCodes {
		keyState, ok := keyPressedStates[v]
		// logger.Infof("", "key status check param:%+v, key:%s, keystate:%+v", GetNamesByCodes(vkCodes), GetNameByCode(v), keyState)
		if !ok {
			keyState = nilKeyState()
		}
		if pressed == 1 {
			if !keyState.Pressed || keyState.LastPressed == nil {
				return false
			}

			if maxLastPressedTime == nil {
				maxLastPressedTime = keyState.LastPressed
			} else {
				if maxLastPressedTime.Sub(*keyState.LastPressed) < 0 {
					maxLastPressedTime = keyState.LastPressed
				}
			}
			if minLastPressedTime == nil {
				minLastPressedTime = keyState.LastPressed
			} else {
				if minLastPressedTime.Sub(*keyState.LastPressed) > 0 {
					minLastPressedTime = keyState.LastPressed
				}
			}
		} else {
			if keyState.Pressed || keyState.LastReleased == nil {
				return false
			}

			if maxLastReleasedTime == nil {
				maxLastReleasedTime = keyState.LastReleased
			} else {
				if maxLastReleasedTime.Sub(*keyState.LastReleased) < 0 {
					maxLastReleasedTime = keyState.LastReleased
				}
			}
			if minLastReleasedTime == nil {
				minLastReleasedTime = keyState.LastReleased
			} else {
				if minLastReleasedTime.Sub(*keyState.LastReleased) > 0 {
					minLastReleasedTime = keyState.LastReleased
				}
			}
		}

	}
	if pressed == 1 {
		durationPressedVal := minLastPressedTime.Sub(*maxLastPressedTime)
		if durationPressedVal < 0 {
			durationPressedVal = -1 * durationPressedVal
		}
		if durationPressedVal > durationBetween {
			return false
		}
	} else {
		durationReleasedVal := minLastReleasedTime.Sub(*maxLastReleasedTime)
		if durationReleasedVal < 0 {
			durationReleasedVal = -1 * durationReleasedVal
		}
		if durationReleasedVal > durationBetween {
			return false
		}
	}

	logger.Infof("", "key status check:%+v", GetNamesByCodes(vkCodes))
	return true
}

func PrintAllKeys() string {
	keysRaw, _ := json.Marshal(listeningKeyReference)
	return string(keysRaw)
}
