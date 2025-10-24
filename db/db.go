package db

import (
	_ "time/tzdata" // 导入时区数据

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"moul.io/zapgorm2"
)

type DBManager struct {
	dbMap map[string]*gorm.DB
}

func NewDBManager() *DBManager {
	return &DBManager{dbMap: make(map[string]*gorm.DB)}
}

func (dbManager *DBManager) InitDB(name string, dbType string, dsn string) (*gorm.DB, error) {
	logger := zapgorm2.New(zap.L())
	logger.SetAsDefault()

	var dialector gorm.Dialector = nil

	switch dbType {
	case "mysql":
		dialector = mysql.Open(dsn)
	case "pgsql":
		dialector = postgres.Open(dsn)
	default:
		panic("unsuport sql drive")
	}
	db, error := gorm.Open(dialector, &gorm.Config{Logger: logger})

	if error != nil {
		return nil, error
	}

	dbManager.dbMap[name] = db
	return db, nil
}

func (dbManager *DBManager) GetDB(name string) *gorm.DB {
	if _, ok := dbManager.dbMap[name]; ok {
		return dbManager.dbMap[name]
	}
	return nil
}

func (dbManager *DBManager) GetDefalutDB() *gorm.DB {
	return dbManager.GetDB("defalut")
}

func (dbManager *DBManager) GetDefaultDB() *gorm.DB {
	return dbManager.GetDB("default")
}

func (dbManager *DBManager) Close() {
	for _, db := range dbManager.dbMap {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}
}
