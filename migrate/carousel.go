package migrate

import "fmt"

type CarouselStruct struct {
	MigrateStruct
}

func (m *MigrateStruct) AutoMigrate()  {
	fmt.Println("轮播图数据迁移！")
}