/**
** @创建时间: 2020/11/3 8:35 下午
** @作者　　: return
** @描述　　:
 */
package bootstrap

import (
	"fmt"
	"github.com/gincmf/cmf/data"
	"github.com/gincmf/cmf/model"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

var (
	db      *gorm.DB
	redisDb *redis.Client
)

func NewDb() *gorm.DB {
	if db == nil {
		config := Conf()
		dbName := config.Database.Name
		//创建不存在的数据库
		model.CreateTable(dbName, config)
		dsn := model.NewDsn(dbName, config)
		db = model.NewDb(dsn, config.Database.Prefix)
	}
	return db
}

func ManualDb(dbName string) *gorm.DB {
	config := Conf()
	dsn := model.NewDsn(dbName, config)
	db = model.NewDb(dsn, config.Database.Prefix)
	return db
}

func NewRedisDb() *redis.Client {
	if redisDb == nil {
		database := Conf().Redis
		empty := data.Redis{}
		if database != empty {
			if database.Host == "" {
				panic("redis host not empty")
			}

			if database.Port == "" {
				panic("redis port not empty")
			}
			redisDb = redis.NewClient(&redis.Options{
				Addr:     database.Host + ":" + database.Port,
				Password: database.Pwd,      // no password set
				DB:       database.Database, // use default DB
			})
			fmt.Println("RedisDb：", redisDb)
			result, err := redisDb.Ping().Result()
			if err != nil {
				panic(err.Error())
			}
			fmt.Println("redis连接状态：", result)
		}
	}
	return redisDb
}
