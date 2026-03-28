package main

import (
	"fmt"
	"github.com/Pho-Tue-SoftWare-Solutions-JSC/hitechcloud-agent/cmd/server/cmd"
	"os"
)

// @title HiTechCloud
// @version 2.0
// @description Top-Rated Web-based Linux Server Management Tool
// @termsOfService http://swagger.io/terms/
// @license.name GPL-3.0
// @license.url https://www.gnu.org/licenses/gpl-3.0.html
// @BasePath /api/v2
// @schemes http https

// @securityDefinitions.apikey ApiKeyAuth
// @description Custom Token Format, Format: md5('HiTechCloud' + API-Key + UnixTimestamp).
// @description ```
// @description eg:
// @description curl -X GET "http://{host}:{port}/api/v2/toolbox/device/base" \
// @description -H "X-API-Key: <HiTechCloud_token>" \
// @description -H "X-API-Timestamp: <current_unix_timestamp>"
// @description ```
// @description - `X-API-Key` is the key for the panel API Key.
// @type apiKey
// @in Header
// @name X-API-Key
// @securityDefinitions.apikey Timestamp
// @type apiKey
// @in header
// @name X-API-Timestamp
// @description - `X-API-Timestamp` is the Unix timestamp of the current time in seconds.

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
