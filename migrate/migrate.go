package migrate

type Migrate interface {
	AutoMigrate()
}

type MigrateStruct struct {

}