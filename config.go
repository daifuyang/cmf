package cmf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

//App 配置文件app对象 定义了系统的基础信息
type AppStruct struct {
	Port string
}

//Datebase 配置文件数据库对象 定义了数据库的基础信息
type DatebaseStruct struct {
	Type    string
	Host    string
	User    string
	Pwd     string
	Port    string
	Charset string
	Prefix  string
}

//ConfigDefault 定义了配置文件初始结构
type ConfigDefaultStruct struct {
	App      AppStruct
	Template TemplateMapStruct `json:"template"`
	Datebase DatebaseStruct
}

//ConfigData 定义一个空结构体
type ConfigDataStruct struct {
}

var ConfigData ConfigDataStruct
var Config = ConfigDefaultStruct{}

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
	ConfigData.init(filePath, &Config)

	//初始化配置信息
	TemplateMap.Theme = Config.Template.Theme
	TemplateMap.ThemePath = Config.Template.ThemePath
	TemplateMap.Glob = Config.Template.Glob
	TemplateMap.Static = Config.Template.Static
}
