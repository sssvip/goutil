package confutil

import (
	"github.com/olebedev/config"
	"github.com/sssvip/goutil/logutil"
)

var defaultConfig *config.Config
//默认读取项目启动文件同路径下的config.yml,读取active项的值 采用相应环境变量
func DefaultYmlConfig(findFileInParent bool, path ...string) *config.Config {
	configPath := "config.yml"
	if len(path) > 0 {
		configPath = path[0]
	}
	if defaultConfig == nil {
		allConfig := ConfigPath(configPath)
		if allConfig == nil {
			logutil.Error.Println("can not find file config.yml in root path")
			if findFileInParent {
				return DefaultYmlConfig(false, "../"+configPath)
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

//读取指定路径下的yaml配置文件
func ConfigPath(confPath string) *config.Config {
	cfg, err := config.ParseYamlFile(confPath)
	if err != nil {
		logutil.Error.Println(err, confPath)
		return nil
	}
	return cfg
}
