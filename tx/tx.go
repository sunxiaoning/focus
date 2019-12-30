package tx

import (
	"context"
	"focus/cfg"
	"github.com/jinzhu/gorm"
)

type txManager struct {
	tx *gorm.DB
}

func NewTxManager() *txManager {
	return &txManager{cfg.FocusCtx.DB.Begin()}
}

type TFunRes interface{}
type TFun func(ctx context.Context) (res TFunRes)

func (txManager *txManager) RunTx(ctx context.Context, tFun TFun) (res TFunRes) {
	defer func() {
		if r := recover(); r != nil {
			txManager.tx.Rollback()
			panic(r)
		}
	}()
	ctx = context.WithValue(ctx, "tx", txManager.tx)
	res = tFun(ctx)
	txManager.tx.Commit()
	return res
}
