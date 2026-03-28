package app

import (
	"github.com/Pho-Tue-SoftWare-Solutions-JSC/hitechcloud-agent/utils/docker"
)

func Init() {
	go func() {
		_ = docker.CreateDefaultDockerNetwork()
	}()
}
