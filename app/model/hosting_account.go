package model

type HostingAccount struct {
	BaseModel
	Username     string `gorm:"uniqueIndex;not null;size:32" json:"username"`
	Domain       string `gorm:"not null" json:"domain"`
	Password     string `gorm:"not null" json:"-"`
	Package      string `gorm:"not null" json:"package"`
	Status       string `gorm:"not null;default:active" json:"status"` // active, suspended, terminated
	Shell        string `gorm:"default:/usr/sbin/nologin" json:"shell"`
	HomeDir      string `json:"homeDir"`
	PHPVersion   string `gorm:"default:8.3" json:"phpVersion"`
	DiskQuota    int64  `json:"diskQuota"`    // MB, 0 = unlimited
	BandwidthCap int64  `json:"bandwidthCap"` // GB/month, 0 = unlimited
	MaxDomains   int    `gorm:"default:1" json:"maxDomains"`
	MaxDatabases int    `gorm:"default:1" json:"maxDatabases"`
	MaxFTP       int    `gorm:"default:1" json:"maxFtp"`
	MaxCronjobs  int    `gorm:"default:5" json:"maxCronjobs"`
	MaxEmail     int    `gorm:"default:0" json:"maxEmail"`
	SSLEnabled   bool   `gorm:"default:true" json:"sslEnabled"`
	BackupEnabled bool  `gorm:"default:true" json:"backupEnabled"`
	UID          int    `json:"uid"`
	GID          int    `json:"gid"`
	Remark       string `json:"remark"`
}

func (h HostingAccount) TableName() string {
	return "hosting_accounts"
}

// Account status constants
const (
	AccountStatusActive     = "active"
	AccountStatusSuspended  = "suspended"
	AccountStatusTerminated = "terminated"
)
