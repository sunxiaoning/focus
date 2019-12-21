package pageutil

import (
	"github.com/jinzhu/gorm"
)

func PageQuery(db *gorm.DB, pageIndex int, pageSize int, count *int, results interface{}) {
	db.Count(count)
	db.Offset((pageIndex - 1) * pageSize).Limit(pageSize).Find(results)
}
