package service

import (
	"context"

	"github.com/Pho-Tue-SoftWare-Solutions-JSC/hitechcloud-agent/constant"
	"github.com/Pho-Tue-SoftWare-Solutions-JSC/hitechcloud-agent/global"
	"gorm.io/gorm"
)

func getTxAndContext() (tx *gorm.DB, ctx context.Context) {
	tx = global.DB.Begin()
	ctx = context.WithValue(context.Background(), constant.DB, tx)
	return
}
