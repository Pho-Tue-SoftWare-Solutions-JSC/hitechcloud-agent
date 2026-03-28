package cache

import (
	"github.com/Pho-Tue-SoftWare-Solutions-JSC/hitechcloud-agent/global"
	cachedb "github.com/Pho-Tue-SoftWare-Solutions-JSC/hitechcloud-agent/init/cache/db"
)

func Init() {
	global.CACHE = cachedb.NewCacheDB()
}
