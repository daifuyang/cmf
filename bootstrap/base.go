package bootstrap

import (
	"fmt"
	"github.com/gincmf/cmf/model"
	"github.com/gincmf/cmf/util"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

var Db *gorm.DB

var RedisDb *redis.Client

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
	util.Conf = config

	initClient()
	if dbHost != "" {
		dbName := config.Database.Name
		//创建不存在的数据库
		model.CreateTable(dbName,config)
		dsn := model.NewDsn(dbName,config)
		Db = model.NewDb(dsn,config.Database.Prefix)
		fmt.Println("创建数据库表成功！")
	}
}

func initClient() (err error) {
	RedisDb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err = RedisDb.Ping().Result()
	if err != nil {
		return err
	}
	return nil
}
