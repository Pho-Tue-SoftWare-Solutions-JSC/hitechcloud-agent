package migrations

import (
	"github.com/Pho-Tue-SoftWare-Solutions-JSC/hitechcloud-agent/app/model"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

var AddHostingAccountTable = &gormigrate.Migration{
	ID: "20260325-add-hosting-account",
	Migrate: func(tx *gorm.DB) error {
		return tx.AutoMigrate(&model.HostingAccount{})
	},
}
