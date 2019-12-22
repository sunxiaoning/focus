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
type TFun func(ctx context.Context) (res TFunRes, err error)

func (txManager *txManager) RunTx(ctx context.Context, tFun TFun) (res TFunRes, err error) {
	defer func() {
		if r := recover(); r != nil {
			txManager.tx.Rollback()
			panic(r)
		}
	}()
	ctx = context.WithValue(ctx, "tx", txManager.tx)
	res, err = tFun(ctx)
	if err != nil {
		txManager.tx.Rollback()
	} else {
		txManager.tx.Commit()
	}
	return res, err
}
