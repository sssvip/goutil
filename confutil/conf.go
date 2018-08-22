package confutil

import (
	"github.com/olebedev/config"
	"github.com/sssvip/goutil/logutil"
)

var defaultConfig *config.Config
//DefaultYmlConfig 默认读取项目启动文件同路径下的config.yml,读取active项的值 采用相应环境变量
func DefaultYmlConfig(findFileInParentLevel int, path ...string) *config.Config {
	configPath := "config.yml"
	if len(path) > 0 {
		configPath = path[0]
	}
	if defaultConfig == nil {
		allConfig := ConfigPath(configPath, false)
		if allConfig == nil {
			if findFileInParentLevel > 0 {
				return DefaultYmlConfig(findFileInParentLevel-1, "../"+configPath)
			} else {
				logutil.Error.Println("can not find file config.yml in root path")
			}
			return nil
		}
		var err error
		defaultConfig, err = allConfig.Get(allConfig.UString("active", "dev"))
		if err != nil {
			logutil.Error.Println(err)
			return nil
		}
	}
	return defaultConfig
}

//ConfigPath 读取指定路径下的yaml配置文件
func ConfigPath(confPath string, showLog ...bool) *config.Config {
	cfg, err := config.ParseYamlFile(confPath)
	if err != nil {
		if !(len(showLog) > 0 && showLog[0] == false) {
			logutil.Error.Println(err, confPath)
		}
		return nil
	}
	return cfg
}
