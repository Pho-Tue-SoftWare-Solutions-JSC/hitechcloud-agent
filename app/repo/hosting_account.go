package repo

import (
	"context"

	"github.com/Pho-Tue-SoftWare-Solutions-JSC/hitechcloud-agent/app/model"
	"github.com/Pho-Tue-SoftWare-Solutions-JSC/hitechcloud-agent/global"
	"gorm.io/gorm"
)

type IHostingAccountRepo interface {
	WithByUsername(username string) DBOption
	WithByStatus(status string) DBOption
	WithByDomain(domain string) DBOption

	Create(ctx context.Context, account *model.HostingAccount) error
	Save(ctx context.Context, account *model.HostingAccount) error
	GetFirst(opts ...DBOption) (model.HostingAccount, error)
	List(opts ...DBOption) ([]model.HostingAccount, error)
	Page(page, size int, opts ...DBOption) (int64, []model.HostingAccount, error)
	Delete(ctx context.Context, opts ...DBOption) error
}

func NewIHostingAccountRepo() IHostingAccountRepo {
	return &HostingAccountRepo{}
}

type HostingAccountRepo struct{}

func (h *HostingAccountRepo) WithByUsername(username string) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("username = ?", username)
	}
}

func (h *HostingAccountRepo) WithByStatus(status string) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("status = ?", status)
	}
}

func (h *HostingAccountRepo) WithByDomain(domain string) DBOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("domain = ?", domain)
	}
}

func (h *HostingAccountRepo) Create(ctx context.Context, account *model.HostingAccount) error {
	return getTx(ctx).Create(account).Error
}

func (h *HostingAccountRepo) Save(ctx context.Context, account *model.HostingAccount) error {
	return getTx(ctx).Save(account).Error
}

func (h *HostingAccountRepo) GetFirst(opts ...DBOption) (model.HostingAccount, error) {
	var account model.HostingAccount
	db := getDb(opts...)
	if err := db.First(&account).Error; err != nil {
		return account, err
	}
	return account, nil
}

func (h *HostingAccountRepo) List(opts ...DBOption) ([]model.HostingAccount, error) {
	var accounts []model.HostingAccount
	db := getDb(opts...)
	if err := db.Find(&accounts).Error; err != nil {
		return nil, err
	}
	return accounts, nil
}

func (h *HostingAccountRepo) Page(page, size int, opts ...DBOption) (int64, []model.HostingAccount, error) {
	var accounts []model.HostingAccount
	db := getDb(opts...)
	var count int64
	db = db.Model(&model.HostingAccount{})
	db.Count(&count)
	if err := db.Offset((page - 1) * size).Limit(size).Find(&accounts).Error; err != nil {
		return 0, nil, err
	}
	return count, accounts, nil
}

func (h *HostingAccountRepo) Delete(ctx context.Context, opts ...DBOption) error {
	return getTx(ctx, opts...).Delete(&model.HostingAccount{}).Error
}

func getDb(opts ...DBOption) *gorm.DB {
	db := global.DB
	for _, opt := range opts {
		db = opt(db)
	}
	return db
}

func getTx(ctx context.Context, opts ...DBOption) *gorm.DB {
	tx, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		tx = global.DB
	}
	for _, opt := range opts {
		tx = opt(tx)
	}
	return tx
}
