package client

import (
	"fmt"

	"github.com/Pho-Tue-SoftWare-Solutions-JSC/hitechcloud-agent/global"
)

func loadParamFromVars(key string, vars map[string]interface{}) string {
	if _, ok := vars[key]; !ok {
		if key != "bucket" && key != "port" && key != "authMode" && key != "passPhrase" {
			global.LOG.Errorf("load param %s from vars failed, err: not exist!", key)
		}
		return ""
	}

	return fmt.Sprintf("%v", vars[key])
}
