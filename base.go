package cmf

import (
	"fmt"
	"github.com/gincmf/cmf/router"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var Db *gorm.DB


//定义数据库类
var err error

//2020 05 09
//初始化默认设置
func initDefault() {
	fmt.Println("核心初始化")
	//初始化配置信息
	TemplateMap.Theme = Config.Template.Theme
	TemplateMap.ThemePath = Config.Template.ThemePath
	TemplateMap.Glob = Config.Template.Glob
	TemplateMap.Static = Config.Template.Static

	dbType := Config.Datebase.Type
	dbUser := Config.Datebase.User
	dbPwd := Config.Datebase.Pwd
	dbHost := Config.Datebase.Host
	dbPort := Config.Datebase.Port
	dbName := Config.Datebase.Name
	dbCharset := Config.Datebase.Charset

	router.AppPort = &Config.App.Port

	fmt.Println("port",router.AppPort)

	//连接sql
	Db, err = gorm.Open(dbType, dbUser+":"+dbPwd+"@tcp("+dbHost+":"+dbPort+")/"+dbName+"?charset="+dbCharset)
	if err != nil {
		panic(err)
		defer Db.Close()
	}
}