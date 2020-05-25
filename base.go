package cmf

import (
	"github.com/gincmf/cmf/router"
	"github.com/gincmf/cmf/util"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var Db *gorm.DB


//定义数据库类
var err error

//2020 05 09
//初始化默认设置
func initDefault() {

	config := Conf()

	//初始化配置信息
	TemplateMap.Theme = config.Template.Theme
	TemplateMap.ThemePath = config.Template.ThemePath
	TemplateMap.Glob = config.Template.Glob
	TemplateMap.Static = config.Template.Static

	dbType := config.Database.Type
	dbUser := config.Database.User
	dbPwd := config.Database.Pwd
	dbHost := config.Database.Host
	dbPort := config.Database.Port
	dbName := config.Database.Name
	dbCharset := config.Database.Charset

	router.AppPort = &config.App.Port
	router.AuthCode  = &config.Database.AuthCode

	util.AuthCode =  &config.Database.AuthCode
	//连接sql
	Db, err = gorm.Open(dbType, dbUser+":"+dbPwd+"@tcp("+dbHost+":"+dbPort+")/"+dbName+"?charset="+dbCharset)

	Db.SingularTable(true)

	gorm.DefaultTableNameHandler = func (db *gorm.DB, defaultTableName string) string  {
		return config.Database.Prefix + defaultTableName
	}

	router.Db = Db
	if err != nil {
		panic(err)
		defer Db.Close()
	}
}