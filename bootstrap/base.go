package bootstrap

import (
	"database/sql"
	"fmt"
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

	dbHost := config.Database.Host

	util.AuthCode = &config.Database.AuthCode

	if dbHost != "" {
		dbType := config.Database.Type
		dbUser := config.Database.User
		dbPwd := config.Database.Pwd
		dbPort := config.Database.Port
		dbName := config.Database.Name
		dbCharset := config.Database.Charset

		//创建不存在的数据库
		dataSource := dbUser + ":" + dbPwd + "@tcp(" + dbHost + ":" + dbPort + ")/"
		tempDb, tempErr := sql.Open(dbType, dataSource)
		if tempErr != nil {
			panic(err)
		}

		defer tempDb.Close()

		_, err = tempDb.Exec("CREATE DATABASE IF NOT EXISTS " + dbName)
		if err != nil {
			panic(err)
		}
		tempDb.Close()

		fmt.Println("创建数据库表成功！")

		//连接sql
		Db, err = gorm.Open(dbType, dataSource+dbName+"?charset="+dbCharset)

		if err != nil {
			panic(err)
			defer Db.Close()
		}

		Db.SingularTable(true)

		gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
			return config.Database.Prefix + defaultTableName
		}
	}
}
