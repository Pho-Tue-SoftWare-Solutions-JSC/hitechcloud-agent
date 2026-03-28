package dto

import (
	"github.com/Pho-Tue-SoftWare-Solutions-JSC/hitechcloud-agent/app/model"
)

type SearchTaskLogReq struct {
	Status string `json:"status"`
	Type   string `json:"type"`
	PageInfo
}

type TaskDTO struct {
	model.Task
}
