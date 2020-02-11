package preceiptcoderepo

import (
	"context"
	"focus/cfg"
	ppaytype "focus/types/ppay"
	dbutil "focus/util/db"
	"github.com/jinzhu/gorm"
)

const (
	tabName     = "personal_receipt_code"
	normalQuery = "status = 1"
	maxAmount   = "9999.99"
)

func GetByAccountIdAndAmount(ctx context.Context, amount string, accountId int) *ppaytype.PReceiptCodeEntity {
	db := getDb(ctx)
	db = db.Where("payee_amount = ? and payee_account_id = ?", amount, accountId)
	db = db.Where(normalQuery)
	var receiptCodeEntity ppaytype.PReceiptCodeEntity
	dbutil.NewDbExecutor(db.Find(&receiptCodeEntity))
	if amount != maxAmount && receiptCodeEntity.ID == 0 {
		return GetByAccountIdAndAmount(ctx, maxAmount, accountId)
	}
	return &receiptCodeEntity
}

func getDb(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value("tx").(*gorm.DB)
	if ok {
		tx = tx.Table(tabName)
		return tx
	}
	return cfg.FocusCtx.DB.Table(tabName)
}
