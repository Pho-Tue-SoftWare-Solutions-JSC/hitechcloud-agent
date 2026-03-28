package service

import (
	"os"
	"sort"
	"strings"
	"time"

	"github.com/Pho-Tue-SoftWare-Solutions-JSC/hitechcloud-agent/global"
)

type LogService struct{}

type ILogService interface {
	ListSystemLogFile() ([]string, error)
}

func NewILogService() ILogService {
	return &LogService{}
}

func (u *LogService) ListSystemLogFile() ([]string, error) {
	var listFile []string
	files, err := os.ReadDir(global.Dir.LogDir)
	if err != nil {
		return nil, err
	}
	listMap := make(map[string]struct{})
	for _, item := range files {
		if item.IsDir() || !strings.HasPrefix(item.Name(), "HiTechCloud") {
			continue
		}
		if item.Name() == "HiTechCloud.log" || item.Name() == "HiTechCloud-Core.log" {
			itemName := time.Now().Format("2006-01-02")
			if _, ok := listMap[itemName]; ok {
				continue
			}
			listMap[itemName] = struct{}{}
			listFile = append(listFile, itemName)
			continue
		}
		itemFileName := strings.TrimPrefix(item.Name(), "HiTechCloud-Core-")
		itemFileName = strings.TrimPrefix(itemFileName, "HiTechCloud-")
		itemFileName = strings.TrimSuffix(itemFileName, ".gz")
		itemFileName = strings.TrimSuffix(itemFileName, ".log")
		if len(itemFileName) == 0 {
			continue
		}
		if _, ok := listMap[itemFileName]; ok {
			continue
		}
		listMap[itemFileName] = struct{}{}
		listFile = append(listFile, itemFileName)
	}
	if len(listFile) < 2 {
		return listFile, nil
	}
	sort.Slice(listFile, func(i, j int) bool {
		return listFile[i] > listFile[j]
	})

	return listFile, nil
}
