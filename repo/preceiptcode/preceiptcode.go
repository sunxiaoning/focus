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
	MaxAmount   = "9999.99"
)

func GetById(ctx context.Context, id int) *ppaytype.PReceiptCodeEntity {
	db := getDb(ctx)
	db = db.Where("id = ?", id)
	db = db.Where(normalQuery)
	var receiptCodeEntity ppaytype.PReceiptCodeEntity
	dbutil.NewDbExecutor(db.Find(&receiptCodeEntity))
	return &receiptCodeEntity
}

func GetByAccountIdAndAmount(ctx context.Context, amount string, accountId int) *ppaytype.PReceiptCodeEntity {
	db := getDb(ctx)
	db = db.Where("payee_amount = ? and payee_account_id = ?", amount, accountId)
	db = db.Where(normalQuery)
	var receiptCodeEntity ppaytype.PReceiptCodeEntity
	dbutil.NewDbExecutor(db.Find(&receiptCodeEntity))
	return &receiptCodeEntity
}

func Create(ctx context.Context, receiptCodeEntity *ppaytype.PReceiptCodeEntity) {
	db := getDb(ctx)
	dbutil.NewDbExecutor(db.Create(receiptCodeEntity))
}

func UpdateReceiptCodeUrl(ctx context.Context, id int, qrCodeUrl string) {
	db := getDb(ctx)
	db = db.Where("id = ?", id)
	db = db.Where(normalQuery)
	dbutil.NewDbExecutor(db.Update("receipt_code_url", qrCodeUrl))
}

func getDb(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value("tx").(*gorm.DB)
	if ok {
		tx = tx.Table(tabName)
		return tx
	}
	return cfg.FocusCtx.DB.Table(tabName)
}
