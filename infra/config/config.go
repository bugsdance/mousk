package config

import (
	"fmt"
	"log"
	"mousk/common/logger"
	"os"
	"os/user"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

var settingsVar Settings

type Mode int

var (
	ModeLoadFromDefault Mode = Mode(0)
	ModeLoadFromUser    Mode = Mode(1)
)

var defaultFilePath = "./conf/default.toml"

// var filePath = "./confs.toml"

// 配置项枚举
const (
	MouseSlowSpeed = iota
	MouseSlowUp
	MouseSlowDown
	MouseSlowLeft
	MouseSlowRight
	MouseFastSpeed
	MouseFastUp
	MouseFastDown
	MouseFastLeft
	MouseFastRight

	KeyDown = iota
	KeyHold
	KeyUp

	MouseUp
	MouseHold
	MouseDown
	MouseMove
	MouseDrag
	MouseWheel

	WheelUp
	WheelDown
)

// type Settings struct {
// 	Mouse struct {
// 		Fast MouseConfig
// 		Slow MouseConfig
// 	}
// }

type Settings struct {
	StartOnSystemUp bool `toml:"StartOnSystemUp"`
	PresetFunc      struct {
		ForceQuit         []string `toml:"ForceQuit"`
		ActiveMode        []string `toml:"ActiveMode"`
		ResetSetting      []string `toml:"ResetSetting"`
		ToggleControlMode []string `toml:"ToggleControlMode"`
		TmpQuitMode       []string `toml:"TmpQuitMode"`
		QuitMode          []string `toml:"QuitMode"`
		OpenSetting       []string `toml:"OpenSetting"`
		MouseMove         struct {
			Fast struct {
				Down  []string `toml:"Down"`
				Left  []string `toml:"Left"`
				Right []string `toml:"Right"`
				Up    []string `toml:"Up"`
			} `toml:"Fast"`
			Slow struct {
				Down  []string `toml:"Down"`
				Left  []string `toml:"Left"`
				Right []string `toml:"Right"`
				Up    []string `toml:"Up"`
			} `toml:"Slow"`
			SpeedLevel struct {
				LevelSwitch []int `toml:"LevelSwitch"`
			} `toml:"SpeedLevel"`
		} `toml:"MouseMove"`
		MouseScroll struct {
			Fast struct {
				Down  []string `toml:"Down"`
				Left  []string `toml:"Left"`
				Right []string `toml:"Right"`
				Up    []string `toml:"Up"`
			} `toml:"Fast"`
			Slow struct {
				Down  []string `toml:"Down"`
				Left  []string `toml:"Left"`
				Right []string `toml:"Right"`
				Up    []string `toml:"Up"`
			} `toml:"Slow"`
			SpeedLevel struct {
				LevelSwitch []int `toml:"LevelSwitch"`
			} `toml:"SpeedLevel"`
		} `toml:"MouseScroll"`
		MouseLeftButtonClick struct {
			Primary   string `toml:"Primary"`
			Secondary string `toml:"Secondary"`
		} `toml:"MouseLeftButtonClick"`
		MouseRightButtonClick struct {
			Primary   string `toml:"Primary"`
			Secondary string `toml:"Secondary"`
		} `toml:"MouseRightButtonClick"`
		MouseLeftButtonHold struct {
			Primary   string `toml:"Primary"`
			Secondary string `toml:"Secondary"`
		} `toml:"MouseLeftButtonHold"`
	} `toml:"PresetFunc"`
}

func Init() {
	// TODO
	// if InitConfigFile() != nil {
	// 	log.Fatalf("init config file error")
	// 	return
	// }

	if LoadSettingsFromFile(ModeLoadFromDefault) != nil {
		log.Fatalf("load config file error")
		return
	}
}

func LoadSettingsFromFile(mode Mode) error {
	var filePath string
	var err error
	if mode == 0 {
		filePath = getEmbededConfigPath()
	} else if mode == 1 {
		_, filePath, err = getUserConfigPath()
	}
	if err != nil {
		return err
	}
	if _, err := toml.DecodeFile(filePath, &settingsVar); err != nil {
		return fmt.Errorf("decoding TOML:%+v", err)
	}

	// 输出生成的结构体
	logger.Infof("", "%#v", settingsVar)
	return nil
}

func GetSettings() *Settings {
	return &settingsVar
}

func WriteSettings() error {
	_, filePath, err := getUserConfigPath()
	if err != nil {
		return err
	}
	// 打开文件
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("open TOML:%+v", err)
	}
	defer file.Close()
	// 编码为TOML格式并写入文件
	if err := toml.NewEncoder(file).Encode(settingsVar); err != nil {
		return fmt.Errorf("writeback TOML:%+v", err)
	}
	return nil
}

func getEmbededConfigPath() string {
	return defaultFilePath
}

func getUserConfigPath() (string, string, error) {
	// 获取当前用户的主目录
	usr, err := user.Current()
	if err != nil {
		return "", "", err
	}

	// 构建文件路径
	configDir := filepath.Join(usr.HomeDir, ".config", "mousk")
	filePath := filepath.Join(configDir, "confs.toml")
	return configDir, filePath, nil
}

func InitConfigFile() error {
	// 构建文件路径
	configDir, filePath, err := getUserConfigPath()
	if err != nil {
		return err
	}
	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// 创建目录
		err := os.MkdirAll(configDir, 0755)
		if err != nil {
			return fmt.Errorf("creating directories:%+v", err)
		}

		// 创建文件
		file, err := os.Create(filePath)
		if err != nil {
			return fmt.Errorf("creating file:%+v", err)
		}
		defer file.Close()

		err = RestoreSettings()
		if err != nil {
			return err
		}
		logger.Infof("", "Config file created at:", filePath)
	} else {
		err = LoadSettingsFromFile(ModeLoadFromUser)
		if err != nil {
			return err
		}
		logger.Infof("", "Config file already exists at:", filePath)
	}
	return nil
}

func RestoreSettings() error {

	if err := LoadSettingsFromFile(ModeLoadFromDefault); err != nil {
		return fmt.Errorf("decoding TOML:%+v", err)
	}
	// 输出生成的结构体
	logger.Infof("", "restore setting:%#v", settingsVar)
	if err := WriteSettings(); err != nil {
		return err
	}
	return nil
}
