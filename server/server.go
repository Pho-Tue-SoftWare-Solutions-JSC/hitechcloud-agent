package server

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/Pho-Tue-SoftWare-Solutions-JSC/hitechcloud-agent/app/repo"
	"github.com/Pho-Tue-SoftWare-Solutions-JSC/hitechcloud-agent/constant"
	"github.com/Pho-Tue-SoftWare-Solutions-JSC/hitechcloud-agent/cron"
	"github.com/Pho-Tue-SoftWare-Solutions-JSC/hitechcloud-agent/global"
	"github.com/Pho-Tue-SoftWare-Solutions-JSC/hitechcloud-agent/i18n"
	"github.com/Pho-Tue-SoftWare-Solutions-JSC/hitechcloud-agent/init/app"
	"github.com/Pho-Tue-SoftWare-Solutions-JSC/hitechcloud-agent/init/business"
	"github.com/Pho-Tue-SoftWare-Solutions-JSC/hitechcloud-agent/init/cache"
	"github.com/Pho-Tue-SoftWare-Solutions-JSC/hitechcloud-agent/init/db"
	"github.com/Pho-Tue-SoftWare-Solutions-JSC/hitechcloud-agent/init/dir"
	"github.com/Pho-Tue-SoftWare-Solutions-JSC/hitechcloud-agent/init/firewall"
	"github.com/Pho-Tue-SoftWare-Solutions-JSC/hitechcloud-agent/init/hook"
	"github.com/Pho-Tue-SoftWare-Solutions-JSC/hitechcloud-agent/init/lang"
	"github.com/Pho-Tue-SoftWare-Solutions-JSC/hitechcloud-agent/init/log"
	"github.com/Pho-Tue-SoftWare-Solutions-JSC/hitechcloud-agent/init/migration"
	"github.com/Pho-Tue-SoftWare-Solutions-JSC/hitechcloud-agent/init/router"
	"github.com/Pho-Tue-SoftWare-Solutions-JSC/hitechcloud-agent/init/validator"
	"github.com/Pho-Tue-SoftWare-Solutions-JSC/hitechcloud-agent/init/viper"
	"github.com/Pho-Tue-SoftWare-Solutions-JSC/hitechcloud-agent/utils/encrypt"
	"github.com/Pho-Tue-SoftWare-Solutions-JSC/hitechcloud-agent/utils/re"
)

func Start() {
	re.Init()
	viper.Init()
	dir.Init()
	log.Init()
	db.Init()
	migration.Init()
	i18n.Init()
	cache.Init()
	app.Init()
	lang.Init()
	validator.Init()
	cron.Run()
	hook.Init()
	go firewall.Init()
	InitOthers()

	rootRouter := router.Routers()

	server := &http.Server{
		Handler: rootRouter,
	}

	if global.CONF.Base.Mode != "stable" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	if global.IsMaster {
		_ = os.Remove("/etc/HiTechCloud/agent.sock")
		_ = os.Mkdir("/etc/HiTechCloud", constant.DirPerm)
		listener, err := net.Listen("unix", "/etc/HiTechCloud/agent.sock")
		if err != nil {
			panic(err)
		}
		business.Init()
		_ = server.Serve(listener)
		return
	}

	port := global.CONF.Base.Port
	if port == "" {
		port = "8443"
	}
	server.Addr = fmt.Sprintf("0.0.0.0:%s", port)

	business.Init()

	// Try HTTPS with SSL cert from settings, fall back to self-signed or HTTP
	settingRepo := repo.NewISettingRepo()
	certItem, certErr := settingRepo.Get(settingRepo.WithByKey("ServerCrt"))
	keyItem, keyErr := settingRepo.Get(settingRepo.WithByKey("ServerKey"))

	if certErr == nil && keyErr == nil && certItem.Value != "" && keyItem.Value != "" {
		cert, _ := encrypt.StringDecrypt(certItem.Value)
		key, _ := encrypt.StringDecrypt(keyItem.Value)
		if cert != "" && key != "" {
			tlsCert, err := tls.X509KeyPair([]byte(cert), []byte(key))
			if err == nil {
				server.TLSConfig = &tls.Config{
					Certificates: []tls.Certificate{tlsCert},
					MinVersion:   tls.VersionTLS12,
				}
				global.LOG.Infof("HiTechCloud Agent listening at https://0.0.0.0:%s", port)
				if err := server.ListenAndServeTLS("", ""); err != nil {
					panic(err)
				}
				return
			}
			global.LOG.Warnf("Failed to load TLS cert, falling back to HTTP: %s", err)
		}
	}

	// Fallback: plain HTTP (should be behind reverse proxy in production)
	global.LOG.Infof("HiTechCloud Agent listening at http://0.0.0.0:%s", port)
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
