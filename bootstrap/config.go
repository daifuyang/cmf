package bootstrap

import (
	"encoding/json"
	"fmt"
	"github.com/gincmf/cmf/data"
	"github.com/gincmf/cmf/util"
	"io/ioutil"
)

//ConfigData 定义一个空结构体
type ConfigDataStruct struct {
}

var configData ConfigDataStruct

var tempConfig = &data.TempConfig{}
var config = &data.ConfigDefault{}

//定义空结构体
func (conf *ConfigDataStruct) init(filePath string, v interface{}) {

	//ReadFile函数会读取文件的全部内容，并将结果以[]byte类型返回
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("ReadFile err", err)
		return
	}

	//读取的数据为json格式，需要进行解码
	err = json.Unmarshal(data, v)
	if err != nil {
		return
	}
}

func Initialize(filePath string) {

	configData.init(filePath, &tempConfig)
	configData.init(filePath, &config)

	if tempConfig.Database.Default == "" {
		panic("默认数据库类型不能为空！")
	}

	config.Redis = tempConfig.Database.Redis

	switch tempConfig.Database.Default {
	case "mysql":
		config.Database = tempConfig.Database.Mysql
		break
	default:
		panic("数据库类型不支持或不存在！")
		break
	}

	TemplateMap.Theme = config.Template.Theme
	TemplateMap.ThemePath = config.Template.ThemePath
	TemplateMap.Glob = config.Template.Glob
	TemplateMap.Static = config.Template.Static
	util.Conf = config
	//initDefault()
}

func Conf() *data.ConfigDefault {
	return config
}
