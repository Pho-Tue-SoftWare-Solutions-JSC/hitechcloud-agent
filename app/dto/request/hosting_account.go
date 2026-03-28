package request

type HostingAccountCreate struct {
	Username     string `json:"username" validate:"required,min=3,max=32,alphanum"`
	Password     string `json:"password" validate:"required,min=8"`
	Domain       string `json:"domain" validate:"required,fqdn"`
	Package      string `json:"package" validate:"required"`
	PHPVersion   string `json:"phpVersion"`
	DiskQuota    int64  `json:"diskQuota"`
	BandwidthCap int64  `json:"bandwidthCap"`
	MaxDomains   int    `json:"maxDomains"`
	MaxDatabases int    `json:"maxDatabases"`
	MaxFTP       int    `json:"maxFtp"`
	MaxCronjobs  int    `json:"maxCronjobs"`
	MaxEmail     int    `json:"maxEmail"`
	SSLEnabled   bool   `json:"sslEnabled"`
	BackupEnabled bool  `json:"backupEnabled"`
	ShellAccess  bool   `json:"shellAccess"`
	Remark       string `json:"remark"`
}

type HostingAccountUpdate struct {
	Package      string `json:"package"`
	PHPVersion   string `json:"phpVersion"`
	DiskQuota    int64  `json:"diskQuota"`
	BandwidthCap int64  `json:"bandwidthCap"`
	MaxDomains   int    `json:"maxDomains"`
	MaxDatabases int    `json:"maxDatabases"`
	MaxFTP       int    `json:"maxFtp"`
	MaxCronjobs  int    `json:"maxCronjobs"`
	MaxEmail     int    `json:"maxEmail"`
	SSLEnabled   bool   `json:"sslEnabled"`
	BackupEnabled bool  `json:"backupEnabled"`
	ShellAccess  bool   `json:"shellAccess"`
	Remark       string `json:"remark"`
}

type HostingAccountPassword struct {
	Password string `json:"password" validate:"required,min=8"`
}
