package preceiptaccountrepo

import (
	"context"
	"focus/cfg"
	ppaytype "focus/types/ppay"
	dbutil "focus/util/db"
	"github.com/jinzhu/gorm"
)

const (
	tabName     = "personal_receipt_account"
	normalQuery = "status = 1"
)

func GetById(ctx context.Context, id int) *ppaytype.PReceiptAccountEntity {
	db := getDb(ctx)
	db = db.Where("id = ?", id)
	db = db.Where(normalQuery)
	var receiptAccountEntity ppaytype.PReceiptAccountEntity
	dbutil.NewDbExecutor(db.Find(&receiptAccountEntity))
	return &receiptAccountEntity
}

func getDb(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value("tx").(*gorm.DB)
	if ok {
		tx = tx.Table(tabName)
		return tx
	}
	return cfg.FocusCtx.DB.Table(tabName)
}
