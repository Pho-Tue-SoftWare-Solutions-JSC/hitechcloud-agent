package service

import (
	"context"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/Pho-Tue-SoftWare-Solutions-JSC/hitechcloud-agent/app/dto/request"
	"github.com/Pho-Tue-SoftWare-Solutions-JSC/hitechcloud-agent/app/dto/response"
	"github.com/Pho-Tue-SoftWare-Solutions-JSC/hitechcloud-agent/app/model"
	"github.com/Pho-Tue-SoftWare-Solutions-JSC/hitechcloud-agent/app/repo"
	"github.com/Pho-Tue-SoftWare-Solutions-JSC/hitechcloud-agent/global"
	"github.com/Pho-Tue-SoftWare-Solutions-JSC/hitechcloud-agent/utils/cmd"
)

type IHostingAccountService interface {
	Create(req request.HostingAccountCreate) error
	GetByUsername(username string) (response.HostingAccountInfo, error)
	List() ([]response.HostingAccountInfo, error)
	Update(username string, req request.HostingAccountUpdate) error
	Suspend(username string) error
	Unsuspend(username string) error
	Terminate(username string) error
	ChangePassword(username string, req request.HostingAccountPassword) error
	GetStats(username string) (response.HostingAccountStats, error)
}

func NewIHostingAccountService() IHostingAccountService {
	return &HostingAccountService{}
}

type HostingAccountService struct{}

var hostingAccountRepo = repo.NewIHostingAccountRepo()

func (s *HostingAccountService) Create(req request.HostingAccountCreate) error {
	// Check if username already exists
	_, err := hostingAccountRepo.GetFirst(hostingAccountRepo.WithByUsername(req.Username))
	if err == nil {
		return fmt.Errorf("username %s already exists", req.Username)
	}

	shell := "/usr/sbin/nologin"
	if req.ShellAccess {
		shell = "/bin/bash"
	}
	homeDir := fmt.Sprintf("/home/%s", req.Username)
	phpVersion := req.PHPVersion
	if phpVersion == "" {
		phpVersion = "8.3"
	}

	// 1. Create Linux user
	createUserCmd := fmt.Sprintf("useradd -m -d %s -s %s %s", homeDir, shell, req.Username)
	if _, err := cmd.ExecWithTimeOut(createUserCmd, 10); err != nil {
		return fmt.Errorf("failed to create user: %v", err)
	}

	// Set password
	setPassCmd := fmt.Sprintf("echo '%s:%s' | chpasswd", req.Username, req.Password)
	if _, err := cmd.ExecWithTimeOut(setPassCmd, 10); err != nil {
		_ = cmd.ExecWithTimeOut(fmt.Sprintf("userdel -r %s", req.Username), 10)
		return fmt.Errorf("failed to set password: %v", err)
	}

	// 2. Create directory structure
	dirs := []string{
		path.Join(homeDir, "public_html"),
		path.Join(homeDir, "domains"),
		path.Join(homeDir, "logs"),
		path.Join(homeDir, "tmp"),
		path.Join(homeDir, "ssl"),
		path.Join(homeDir, "backups"),
	}
	for _, d := range dirs {
		if err := os.MkdirAll(d, 0755); err != nil {
			_ = cmd.ExecWithTimeOut(fmt.Sprintf("userdel -r %s", req.Username), 10)
			return fmt.Errorf("failed to create directory %s: %v", d, err)
		}
	}
	// Set ownership
	chownCmd := fmt.Sprintf("chown -R %s:%s %s", req.Username, req.Username, homeDir)
	_, _ = cmd.ExecWithTimeOut(chownCmd, 10)

	// Create default index.html
	indexPath := path.Join(homeDir, "public_html", "index.html")
	indexContent := fmt.Sprintf("<html><body><h1>Welcome to %s</h1><p>Hosted by HiTechCloud</p></body></html>", req.Domain)
	_ = os.WriteFile(indexPath, []byte(indexContent), 0644)
	_, _ = cmd.ExecWithTimeOut(fmt.Sprintf("chown %s:%s %s", req.Username, req.Username, indexPath), 5)

	// 3. Create PHP-FPM pool
	if err := s.createPHPPool(req.Username, phpVersion, homeDir); err != nil {
		_ = cmd.ExecWithTimeOut(fmt.Sprintf("userdel -r %s", req.Username), 10)
		return fmt.Errorf("failed to create PHP pool: %v", err)
	}

	// 4. Create Nginx vhost
	if err := s.createNginxVhost(req.Username, req.Domain, homeDir, phpVersion); err != nil {
		s.removePHPPool(req.Username, phpVersion)
		_ = cmd.ExecWithTimeOut(fmt.Sprintf("userdel -r %s", req.Username), 10)
		return fmt.Errorf("failed to create nginx vhost: %v", err)
	}

	// 5. Set disk quota
	if req.DiskQuota > 0 {
		s.setDiskQuota(req.Username, req.DiskQuota)
	}

	// 6. Get UID/GID
	uid, gid := s.getUserIDs(req.Username)

	// 7. Reload services
	_, _ = cmd.ExecWithTimeOut("systemctl reload php"+phpVersion+"-fpm 2>/dev/null || true", 10)
	_, _ = cmd.ExecWithTimeOut("nginx -t && systemctl reload nginx", 10)

	// 8. Save to database
	account := &model.HostingAccount{
		Username:      req.Username,
		Domain:        req.Domain,
		Password:      "***",
		Package:       req.Package,
		Status:        model.AccountStatusActive,
		Shell:         shell,
		HomeDir:       homeDir,
		PHPVersion:    phpVersion,
		DiskQuota:     req.DiskQuota,
		BandwidthCap:  req.BandwidthCap,
		MaxDomains:    req.MaxDomains,
		MaxDatabases:  req.MaxDatabases,
		MaxFTP:        req.MaxFTP,
		MaxCronjobs:   req.MaxCronjobs,
		MaxEmail:      req.MaxEmail,
		SSLEnabled:    req.SSLEnabled,
		BackupEnabled: req.BackupEnabled,
		UID:           uid,
		GID:           gid,
		Remark:        req.Remark,
	}
	return hostingAccountRepo.Create(context.Background(), account)
}

func (s *HostingAccountService) GetByUsername(username string) (response.HostingAccountInfo, error) {
	account, err := hostingAccountRepo.GetFirst(hostingAccountRepo.WithByUsername(username))
	if err != nil {
		return response.HostingAccountInfo{}, fmt.Errorf("account not found: %s", username)
	}
	return s.toAccountInfo(account), nil
}

func (s *HostingAccountService) List() ([]response.HostingAccountInfo, error) {
	accounts, err := hostingAccountRepo.List()
	if err != nil {
		return nil, err
	}
	var result []response.HostingAccountInfo
	for _, a := range accounts {
		result = append(result, s.toAccountInfo(a))
	}
	return result, nil
}

func (s *HostingAccountService) Update(username string, req request.HostingAccountUpdate) error {
	account, err := hostingAccountRepo.GetFirst(hostingAccountRepo.WithByUsername(username))
	if err != nil {
		return fmt.Errorf("account not found: %s", username)
	}

	if req.Package != "" {
		account.Package = req.Package
	}
	if req.DiskQuota > 0 {
		account.DiskQuota = req.DiskQuota
		s.setDiskQuota(username, req.DiskQuota)
	}
	if req.BandwidthCap >= 0 {
		account.BandwidthCap = req.BandwidthCap
	}
	if req.MaxDomains > 0 {
		account.MaxDomains = req.MaxDomains
	}
	if req.MaxDatabases > 0 {
		account.MaxDatabases = req.MaxDatabases
	}
	if req.MaxFTP > 0 {
		account.MaxFTP = req.MaxFTP
	}
	if req.MaxCronjobs > 0 {
		account.MaxCronjobs = req.MaxCronjobs
	}
	account.MaxEmail = req.MaxEmail
	account.SSLEnabled = req.SSLEnabled
	account.BackupEnabled = req.BackupEnabled
	account.Remark = req.Remark

	// Update PHP version if changed
	if req.PHPVersion != "" && req.PHPVersion != account.PHPVersion {
		s.removePHPPool(username, account.PHPVersion)
		s.createPHPPool(username, req.PHPVersion, account.HomeDir)
		s.updateNginxPHP(username, account.Domain, req.PHPVersion)
		account.PHPVersion = req.PHPVersion
		_, _ = cmd.ExecWithTimeOut("systemctl reload php"+req.PHPVersion+"-fpm 2>/dev/null || true", 10)
		_, _ = cmd.ExecWithTimeOut("nginx -t && systemctl reload nginx", 10)
	}

	// Update shell access
	if req.ShellAccess {
		account.Shell = "/bin/bash"
	} else {
		account.Shell = "/usr/sbin/nologin"
	}
	_, _ = cmd.ExecWithTimeOut(fmt.Sprintf("usermod -s %s %s", account.Shell, username), 10)

	return hostingAccountRepo.Save(context.Background(), &account)
}

func (s *HostingAccountService) Suspend(username string) error {
	account, err := hostingAccountRepo.GetFirst(hostingAccountRepo.WithByUsername(username))
	if err != nil {
		return fmt.Errorf("account not found: %s", username)
	}
	if account.Status == model.AccountStatusSuspended {
		return nil
	}

	// Disable Linux user
	_, _ = cmd.ExecWithTimeOut(fmt.Sprintf("usermod -L -s /usr/sbin/nologin %s", username), 10)

	// Stop PHP-FPM pool (rename config to .disabled)
	poolFile := fmt.Sprintf("/etc/php/%s/fpm/pool.d/%s.conf", account.PHPVersion, username)
	_, _ = cmd.ExecWithTimeOut(fmt.Sprintf("mv %s %s.disabled 2>/dev/null || true", poolFile, poolFile), 5)
	_, _ = cmd.ExecWithTimeOut("systemctl reload php"+account.PHPVersion+"-fpm 2>/dev/null || true", 10)

	// Update Nginx to return 403
	s.suspendNginxVhost(username, account.Domain)
	_, _ = cmd.ExecWithTimeOut("nginx -t && systemctl reload nginx", 10)

	account.Status = model.AccountStatusSuspended
	return hostingAccountRepo.Save(context.Background(), &account)
}

func (s *HostingAccountService) Unsuspend(username string) error {
	account, err := hostingAccountRepo.GetFirst(hostingAccountRepo.WithByUsername(username))
	if err != nil {
		return fmt.Errorf("account not found: %s", username)
	}
	if account.Status == model.AccountStatusActive {
		return nil
	}

	// Unlock Linux user
	_, _ = cmd.ExecWithTimeOut(fmt.Sprintf("usermod -U -s %s %s", account.Shell, username), 10)

	// Restore PHP-FPM pool
	poolFile := fmt.Sprintf("/etc/php/%s/fpm/pool.d/%s.conf", account.PHPVersion, username)
	_, _ = cmd.ExecWithTimeOut(fmt.Sprintf("mv %s.disabled %s 2>/dev/null || true", poolFile, poolFile), 5)
	_, _ = cmd.ExecWithTimeOut("systemctl reload php"+account.PHPVersion+"-fpm 2>/dev/null || true", 10)

	// Restore Nginx vhost
	s.unsuspendNginxVhost(username, account.Domain, account.HomeDir, account.PHPVersion)
	_, _ = cmd.ExecWithTimeOut("nginx -t && systemctl reload nginx", 10)

	account.Status = model.AccountStatusActive
	return hostingAccountRepo.Save(context.Background(), &account)
}

func (s *HostingAccountService) Terminate(username string) error {
	account, err := hostingAccountRepo.GetFirst(hostingAccountRepo.WithByUsername(username))
	if err != nil {
		return fmt.Errorf("account not found: %s", username)
	}

	// Remove Nginx vhost
	vhostFile := fmt.Sprintf("/etc/nginx/sites-enabled/%s.conf", account.Domain)
	vhostAvail := fmt.Sprintf("/etc/nginx/sites-available/%s.conf", account.Domain)
	_, _ = cmd.ExecWithTimeOut(fmt.Sprintf("rm -f %s %s", vhostFile, vhostAvail), 5)

	// Remove PHP-FPM pool
	s.removePHPPool(username, account.PHPVersion)

	// Reload services
	_, _ = cmd.ExecWithTimeOut("systemctl reload php"+account.PHPVersion+"-fpm 2>/dev/null || true", 10)
	_, _ = cmd.ExecWithTimeOut("nginx -t && systemctl reload nginx", 10)

	// Remove disk quota
	_, _ = cmd.ExecWithTimeOut(fmt.Sprintf("setquota -u %s 0 0 0 0 / 2>/dev/null || true", username), 5)

	// Remove Linux user and home directory
	_, _ = cmd.ExecWithTimeOut(fmt.Sprintf("pkill -u %s 2>/dev/null || true", username), 5)
	_, _ = cmd.ExecWithTimeOut(fmt.Sprintf("userdel -r %s 2>/dev/null || true", username), 10)

	// Remove from database
	account.Status = model.AccountStatusTerminated
	return hostingAccountRepo.Delete(context.Background(), hostingAccountRepo.WithByUsername(username))
}

func (s *HostingAccountService) ChangePassword(username string, req request.HostingAccountPassword) error {
	_, err := hostingAccountRepo.GetFirst(hostingAccountRepo.WithByUsername(username))
	if err != nil {
		return fmt.Errorf("account not found: %s", username)
	}

	setPassCmd := fmt.Sprintf("echo '%s:%s' | chpasswd", username, req.Password)
	if _, err := cmd.ExecWithTimeOut(setPassCmd, 10); err != nil {
		return fmt.Errorf("failed to change password: %v", err)
	}
	return nil
}

func (s *HostingAccountService) GetStats(username string) (response.HostingAccountStats, error) {
	account, err := hostingAccountRepo.GetFirst(hostingAccountRepo.WithByUsername(username))
	if err != nil {
		return response.HostingAccountStats{}, fmt.Errorf("account not found: %s", username)
	}

	stats := response.HostingAccountStats{
		Username:    username,
		DiskQuotaMB: account.DiskQuota,
		BandwidthCap: account.BandwidthCap,
	}

	// Get disk usage
	diskOut, err := cmd.ExecWithTimeOut(fmt.Sprintf("du -sm /home/%s 2>/dev/null | awk '{print $1}'", username), 10)
	if err == nil {
		diskOut = strings.TrimSpace(diskOut)
		if v, e := strconv.ParseInt(diskOut, 10, 64); e == nil {
			stats.DiskUsedMB = v
		}
	}

	return stats, nil
}

// --- Internal helpers ---

func (s *HostingAccountService) createPHPPool(username, phpVersion, homeDir string) error {
	poolConf := fmt.Sprintf(`[%s]
user = %s
group = %s
listen = /run/php/%s.sock
listen.owner = www-data
listen.group = www-data
listen.mode = 0660

pm = ondemand
pm.max_children = 5
pm.process_idle_timeout = 10s
pm.max_requests = 200

php_admin_value[open_basedir] = %s:/tmp:/usr/share/php
php_admin_value[upload_tmp_dir] = %s/tmp
php_admin_value[session.save_path] = %s/tmp
php_admin_value[error_log] = %s/logs/php_error.log
php_admin_value[disable_functions] = exec,passthru,shell_exec,system,proc_open,popen
php_admin_value[memory_limit] = 256M
php_admin_value[max_execution_time] = 30
php_admin_value[post_max_size] = 64M
php_admin_value[upload_max_filesize] = 64M
`, username, username, username, username, homeDir, homeDir, homeDir, homeDir)

	poolFile := fmt.Sprintf("/etc/php/%s/fpm/pool.d/%s.conf", phpVersion, username)
	return os.WriteFile(poolFile, []byte(poolConf), 0644)
}

func (s *HostingAccountService) removePHPPool(username, phpVersion string) {
	poolFile := fmt.Sprintf("/etc/php/%s/fpm/pool.d/%s.conf", phpVersion, username)
	_ = os.Remove(poolFile)
	_ = os.Remove(poolFile + ".disabled")
}

func (s *HostingAccountService) createNginxVhost(username, domain, homeDir, phpVersion string) error {
	vhostConf := fmt.Sprintf(`server {
    listen 80;
    listen [::]:80;
    server_name %s www.%s;

    root %s/public_html;
    index index.php index.html index.htm;

    access_log %s/logs/access.log;
    error_log %s/logs/error.log;

    # Security headers
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;

    # Deny access to hidden files
    location ~ /\. {
        deny all;
    }

    location / {
        try_files $uri $uri/ /index.php?$query_string;
    }

    location ~ \.php$ {
        fastcgi_pass unix:/run/php/%s.sock;
        fastcgi_param SCRIPT_FILENAME $realpath_root$fastcgi_script_name;
        include fastcgi_params;
        fastcgi_read_timeout 300;
    }

    # Static file caching
    location ~* \.(jpg|jpeg|gif|png|css|js|ico|xml|svg|woff|woff2|ttf|eot)$ {
        expires 30d;
        add_header Cache-Control "public, immutable";
    }

    # Deny access to sensitive files
    location ~* \.(env|git|htaccess|htpasswd|ini|log|sh|sql|conf|bak)$ {
        deny all;
    }

    client_max_body_size 64M;
}
`, domain, domain, homeDir, homeDir, homeDir, username)

	// Write vhost config
	vhostFile := fmt.Sprintf("/etc/nginx/sites-available/%s.conf", domain)
	if err := os.WriteFile(vhostFile, []byte(vhostConf), 0644); err != nil {
		return err
	}

	// Enable site
	enabledFile := fmt.Sprintf("/etc/nginx/sites-enabled/%s.conf", domain)
	return os.Symlink(vhostFile, enabledFile)
}

func (s *HostingAccountService) suspendNginxVhost(username, domain string) {
	suspendConf := fmt.Sprintf(`server {
    listen 80;
    listen [::]:80;
    server_name %s www.%s;
    return 403;
}
`, domain, domain)
	vhostFile := fmt.Sprintf("/etc/nginx/sites-available/%s.conf", domain)
	_ = os.WriteFile(vhostFile, []byte(suspendConf), 0644)
}

func (s *HostingAccountService) unsuspendNginxVhost(username, domain, homeDir, phpVersion string) {
	_ = s.createNginxVhost(username, domain, homeDir, phpVersion)
}

func (s *HostingAccountService) updateNginxPHP(username, domain, phpVersion string) {
	// Read current config and update PHP socket path
	vhostFile := fmt.Sprintf("/etc/nginx/sites-available/%s.conf", domain)
	content, err := os.ReadFile(vhostFile)
	if err != nil {
		return
	}
	// Replace socket path
	updated := strings.ReplaceAll(string(content),
		fmt.Sprintf("unix:/run/php/%s.sock", username),
		fmt.Sprintf("unix:/run/php/%s.sock", username))
	_ = os.WriteFile(vhostFile, []byte(updated), 0644)
}

func (s *HostingAccountService) setDiskQuota(username string, quotaMB int64) {
	// softlimit = quota, hardlimit = quota + 10%
	hard := quotaMB + quotaMB/10
	quotaCmd := fmt.Sprintf("setquota -u %s %dM %dM 0 0 / 2>/dev/null || true", username, quotaMB, hard)
	_, _ = cmd.ExecWithTimeOut(quotaCmd, 5)
}

func (s *HostingAccountService) getUserIDs(username string) (int, int) {
	uidOut, _ := cmd.ExecWithTimeOut(fmt.Sprintf("id -u %s", username), 5)
	gidOut, _ := cmd.ExecWithTimeOut(fmt.Sprintf("id -g %s", username), 5)
	uid, _ := strconv.Atoi(strings.TrimSpace(uidOut))
	gid, _ := strconv.Atoi(strings.TrimSpace(gidOut))
	return uid, gid
}

func (s *HostingAccountService) toAccountInfo(account model.HostingAccount) response.HostingAccountInfo {
	info := response.HostingAccountInfo{
		ID:            account.ID,
		Username:      account.Username,
		Domain:        account.Domain,
		Package:       account.Package,
		Status:        account.Status,
		PHPVersion:    account.PHPVersion,
		HomeDir:       account.HomeDir,
		DiskQuota:     account.DiskQuota,
		BandwidthCap:  account.BandwidthCap,
		MaxDomains:    account.MaxDomains,
		MaxDatabases:  account.MaxDatabases,
		MaxFTP:        account.MaxFTP,
		MaxCronjobs:   account.MaxCronjobs,
		MaxEmail:      account.MaxEmail,
		SSLEnabled:    account.SSLEnabled,
		BackupEnabled: account.BackupEnabled,
		UID:           account.UID,
		GID:           account.GID,
		Remark:        account.Remark,
		CreatedAt:     account.CreatedAt,
	}
	// Get current disk usage
	diskOut, err := cmd.ExecWithTimeOut(fmt.Sprintf("du -sm /home/%s 2>/dev/null | awk '{print $1}'", account.Username), 10)
	if err == nil {
		diskOut = strings.TrimSpace(diskOut)
		if v, e := strconv.ParseInt(diskOut, 10, 64); e == nil {
			info.DiskUsed = v
		}
	}
	return info
}
