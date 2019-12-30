package dbutil

import (
	"focus/types"
	"github.com/jinzhu/gorm"
)

type dbExecutor struct {
	DB *gorm.DB
}

func NewDbExecutor(db *gorm.DB) *dbExecutor {
	dbExe := &dbExecutor{db}
	dbExe.checkErr()
	return dbExe
}

func (dbExecutor *dbExecutor) PageQuery(pageIndex int, pageSize int, count *int, results interface{}) {
	dbExecutor.DB.Count(count).Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(results)
}

func (dbExecutor *dbExecutor) RowsAffected() int64 {
	return dbExecutor.DB.RowsAffected
}

func (dbExecutor *dbExecutor) checkErr() {
	if dbExecutor.DB.Error != nil && !gorm.IsRecordNotFoundError(dbExecutor.DB.Error) {
		panic(types.DbErr(dbExecutor.DB.Error))
	}
}
