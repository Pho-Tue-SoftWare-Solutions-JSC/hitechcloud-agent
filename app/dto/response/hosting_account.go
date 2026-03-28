package response

import "time"

type HostingAccountInfo struct {
	ID           uint      `json:"id"`
	Username     string    `json:"username"`
	Domain       string    `json:"domain"`
	Package      string    `json:"package"`
	Status       string    `json:"status"`
	PHPVersion   string    `json:"phpVersion"`
	HomeDir      string    `json:"homeDir"`
	DiskQuota    int64     `json:"diskQuota"`
	DiskUsed     int64     `json:"diskUsed"`
	BandwidthCap int64     `json:"bandwidthCap"`
	BandwidthUsed int64    `json:"bandwidthUsed"`
	MaxDomains   int       `json:"maxDomains"`
	MaxDatabases int       `json:"maxDatabases"`
	MaxFTP       int       `json:"maxFtp"`
	MaxCronjobs  int       `json:"maxCronjobs"`
	MaxEmail     int       `json:"maxEmail"`
	SSLEnabled   bool      `json:"sslEnabled"`
	BackupEnabled bool     `json:"backupEnabled"`
	UID          int       `json:"uid"`
	GID          int       `json:"gid"`
	Remark       string    `json:"remark"`
	CreatedAt    time.Time `json:"createdAt"`
}

type HostingAccountStats struct {
	Username      string  `json:"username"`
	DiskUsedMB    int64   `json:"diskUsedMB"`
	DiskQuotaMB   int64   `json:"diskQuotaMB"`
	BandwidthUsed int64   `json:"bandwidthUsedGB"`
	BandwidthCap  int64   `json:"bandwidthCapGB"`
	CPUUsage      float64 `json:"cpuUsage"`
	MemoryUsedMB  int64   `json:"memoryUsedMB"`
	WebsiteCount  int     `json:"websiteCount"`
	DatabaseCount int     `json:"databaseCount"`
	FTPCount      int     `json:"ftpCount"`
	CronjobCount  int     `json:"cronjobCount"`
}
