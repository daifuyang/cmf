/**
** @创建时间: 2020/10/4 9:13 下午
** @作者　　: return
** @描述　　:
 */
package data

//App 配置文件app对象 定义了系统的基础信息
type AppStruct struct {
	Port     string
	AppName  string
	AppDebug bool
	Domain   string
}

//Database 配置文件数据库对象 定义了数据库的基础信息
type DatabaseStruct struct {
	Type     string `json:"type"`
	Host     string `json:"hostname"`
	Name     string `json:"database"`
	User     string `json:"username"`
	Pwd      string `json:"password"`
	Port     string `json:"hostport"`
	Charset  string `json:"charset"`
	Prefix   string `json:"prefix"`
	AuthCode string `json:"authcode"`
}

//ConfigDefault 定义了配置文件初始结构
type ConfigDefaultStruct struct {
	App      AppStruct
	Template TemplateMapStruct `json:"template"`
	Database DatabaseStruct
	Token    string
}

//ConfigData 定义一个空结构体
type ConfigDataStruct struct {
}
