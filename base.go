package cmf

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

//定义数据库类

var Db *gorm.DB
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
	//连接sql
	Db, err = gorm.Open(dbType, dbUser+":"+dbPwd+"@tcp("+dbHost+":"+dbPort+")/"+dbName+"?charset="+dbCharset)
	if err != nil {
		fmt.Println("出错", err)
		return
	}
	autoMigrate()
	defer Db.Close()
}

//2020 05 09
//初始化数据库迁移
func autoMigrate(){
	
}