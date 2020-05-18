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

//Database 配置文件数据库对象 定义了数据库的基础信息
type DatabaseStruct struct {
	Type    string `json:"type"`
	Host    string `json:"hostname"`
	Name    string `json:"database"`
	User    string `json:"username"`
	Pwd     string `json:"password"`
	Port    string `json:"hostport"`
	Charset string `json:"charset"`
	Prefix  string `json:"prefix"`
	AuthCode string `json:"auth_code"`
}

//ConfigDefault 定义了配置文件初始结构
type ConfigDefaultStruct struct {
	App      AppStruct
	Template TemplateMapStruct `json:"template"`
	Database DatabaseStruct
}

//ConfigData 定义一个空结构体
type ConfigDataStruct struct {
}

var configData ConfigDataStruct
var config = ConfigDefaultStruct{}

//定义空结构体
func (conf *ConfigDataStruct) init(filePath string, v interface{}) {

	//ReadFile函数会读取文件的全部内容，并将结果以[]byte类型返回
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("ReadFile err", err)
		return
	}

	fmt.Println("data",string(data[:]))

	//读取的数据为json格式，需要进行解码
	err = json.Unmarshal(data, v)
	if err != nil {
		return
	}
}

func Initialize(filePath string) {
	configData.init(filePath, &config)
	initDefault()
}

func Conf() *ConfigDefaultStruct {

	//if attr == "" {
	//	panic(errors.New("属性不能为空！"))
	//}
	//attrArr := strings.Split(attr,".")
	//fmt.Println("attrArr",attrArr)
	//refc := reflect.ValueOf(config)
	//
	//for k,v := range attrArr{
	//	fmt.Println(k,v)
	//	refc = refc.FieldByName(v)
	//	fmt.Println(v,refc.Type())
	//}
	//
	//fmt.Println("value：",refc)
	//
	//return refc.String()

	return &config

}


