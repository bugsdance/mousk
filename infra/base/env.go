package base

import "mousk/infra/config"

var (
	mode              = 0 // 0:normal, 1:control
	moveSpeedLevel    = 3 // the speed of you mouse movement
	scrollSpeedLevel  = 3 // the speed of you mouse scroll
	optionViewVisable = false
)

const (
	ModeNormal  = 0
	ModeControl = 1
)

func ToggleOptionViewState() bool {
	optionViewVisable = !optionViewVisable
	return optionViewVisable
}

func SetMoveSpeedLevel(speedLevelArg int) {
	// speedLevel = speedLevelArg
	// speedLevel = int(math.Round(math.Log10(1000 * math.Pow(float64(speedLevelArg), 10))))
	// speedLevel = int(math.Log(1000 * math.Pow(float64(speedLevelArg), 5)))
	// moveSpeedLevel = int(math.Pow(float64(speedLevelArg), 2))
	moveSpeedLevel = speedLevelArg
}

func SetScrollSpeedLevel(speedLevelArg int) {
	// speedLevel = speedLevelArg
	// speedLevel = int(math.Round(math.Log10(1000 * math.Pow(float64(speedLevelArg), 10))))
	// speedLevel = int(math.Log(1000 * math.Pow(float64(speedLevelArg), 5)))
	scrollSpeedLevel = speedLevelArg
}

func GetMoveSpeedLevel() int {
	return moveSpeedLevel
}

func GetMoveSpeed() int {
	return config.GetSettings().PresetFunc.MouseMove.SpeedLevel.LevelSwitch[moveSpeedLevel-1]
}

func GetScrollSpeedLevel() int {
	return scrollSpeedLevel
}

func GetScrollSpeed() int {
	return config.GetSettings().PresetFunc.MouseScroll.SpeedLevel.LevelSwitch[scrollSpeedLevel-1]
}

func SetMode(modeArg int) {
	mode = modeArg
}

func GetMode() int {
	return mode
}

func GetModeDesc() string {
	switch mode {
	case 0:
		return "normal"
	case 1:
		return "control"
	}
	return ""
}
