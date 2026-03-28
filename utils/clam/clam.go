package clam

import (
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/Pho-Tue-SoftWare-Solutions-JSC/hitechcloud-agent/app/model"
	"github.com/Pho-Tue-SoftWare-Solutions-JSC/hitechcloud-agent/app/repo"
	"github.com/Pho-Tue-SoftWare-Solutions-JSC/hitechcloud-agent/app/task"
	"github.com/Pho-Tue-SoftWare-Solutions-JSC/hitechcloud-agent/constant"
	"github.com/Pho-Tue-SoftWare-Solutions-JSC/hitechcloud-agent/global"
	"github.com/Pho-Tue-SoftWare-Solutions-JSC/hitechcloud-agent/i18n"
	"github.com/Pho-Tue-SoftWare-Solutions-JSC/hitechcloud-agent/utils/cmd"
	"github.com/Pho-Tue-SoftWare-Solutions-JSC/hitechcloud-agent/utils/controller"
	"github.com/robfig/cron/v3"
)

func AddScanTask(taskItem *task.Task, clam model.Clam, timeNow string) {
	taskItem.AddSubTask(i18n.GetWithName("Clamscan", clam.Path), func(t *task.Task) error {
		strategy := ""
		switch clam.InfectedStrategy {
		case "remove":
			strategy = "--remove"
		case "move", "copy":
			dir := path.Join(clam.InfectedDir, "HiTechCloud-infected", clam.Name, timeNow)
			taskItem.Log("infected dir: " + dir)
			if _, err := os.Stat(dir); err != nil {
				_ = os.MkdirAll(dir, os.ModePerm)
			}
			strategy = fmt.Sprintf("--%s=%s", clam.InfectedStrategy, dir)
		}
		taskItem.Logf("clamdscan --fdpass %s %s", strategy, clam.Path)
		mgr := cmd.NewCommandMgr(cmd.WithIgnoreExist1(), cmd.WithTimeout(time.Duration(clam.Timeout)*time.Second), cmd.WithTask(*taskItem))
		if err := mgr.RunBashCf("clamdscan --fdpass %s %s", strategy, clam.Path); err != nil {
			return fmt.Errorf("clamdscan failed, %v", err)
		}
		return nil
	}, nil)
}

func AnalysisFromLog(pathItem string, record *model.ClamRecord) {
	file, err := os.ReadFile(pathItem)
	if err != nil {
		return
	}
	lines := strings.Split(string(file), "\n")
	for _, line := range lines {
		if len(line) < 20 {
			continue
		}
		line = line[20:]
		switch {
		case strings.HasPrefix(line, "Infected files: "):
			record.InfectedFiles = strings.TrimPrefix(line, "Infected files: ")
		case strings.HasPrefix(line, "Total errors: "):
			record.TotalError = strings.TrimPrefix(line, "Total errors: ")
		case strings.HasPrefix(line, "Time: "):
			record.ScanTime = strings.TrimPrefix(line, "Time: ")
		}
	}
}

func CheckWithStopAll(withCheck bool, clamRepo repo.IClamRepo) bool {
	if withCheck {
		isExist, _ := controller.CheckExist("clam")
		if !isExist {
			return false
		}
		isActive, _ := controller.CheckActive("clam")
		if isActive {
			return true
		}
	}
	clams, _ := clamRepo.List(repo.WithByStatus(constant.StatusEnable))
	for i := 0; i < len(clams); i++ {
		global.Cron.Remove(cron.EntryID(clams[i].EntryID))
		_ = clamRepo.Update(clams[i].ID, map[string]interface{}{"status": constant.StatusDisable, "entry_id": 0})
	}
	return false
}
