package bootstrap

import (
	"encoding/json"
	"fmt"
	"github.com/gincmf/cmf/data"
	"io/ioutil"
)

//ConfigData 定义一个空结构体
type ConfigDataStruct struct {
}

var configData ConfigDataStruct
var config = &data.ConfigDefaultStruct{}

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
	configData.init(filePath, &config)
	initDefault()
}

func Conf() *data.ConfigDefaultStruct {
	return config
}
