package firewall

import (
	"fmt"
	"os"

	"github.com/Pho-Tue-SoftWare-Solutions-JSC/hitechcloud-agent/app/dto"
	"github.com/Pho-Tue-SoftWare-Solutions-JSC/hitechcloud-agent/app/repo"
	"github.com/Pho-Tue-SoftWare-Solutions-JSC/hitechcloud-agent/app/service"
	"github.com/Pho-Tue-SoftWare-Solutions-JSC/hitechcloud-agent/constant"
	"github.com/Pho-Tue-SoftWare-Solutions-JSC/hitechcloud-agent/global"
	"github.com/Pho-Tue-SoftWare-Solutions-JSC/hitechcloud-agent/utils/firewall"
	firewallClient "github.com/Pho-Tue-SoftWare-Solutions-JSC/hitechcloud-agent/utils/firewall/client"
	"github.com/Pho-Tue-SoftWare-Solutions-JSC/hitechcloud-agent/utils/firewall/client/iptables"
)

func Init() {
	if !needInit() {
		return
	}
	InitPingStatus()
	global.LOG.Info("initializing firewall settings...")
	client, err := firewall.NewFirewallClient()
	if err != nil {
		return
	}
	clientName := client.Name()

	settingRepo := repo.NewISettingRepo()
	if clientName == "ufw" || clientName == "iptables" {
		if err := iptables.LoadRulesFromFile(iptables.FilterTab, iptables.ChainHiTechCloudForward, iptables.ForwardFileName); err != nil {
			global.LOG.Errorf("load forward rules from file failed, err: %v", err)
			return
		}
		if err := iptables.LoadRulesFromFile(iptables.NatTab, iptables.ChainHiTechCloudPreRouting, iptables.ForwardFileName1); err != nil {
			global.LOG.Errorf("load prerouting rules from file failed, err: %v", err)
			return
		}
		if err := iptables.LoadRulesFromFile(iptables.NatTab, iptables.ChainHiTechCloudPostRouting, iptables.ForwardFileName2); err != nil {
			global.LOG.Errorf("load postrouting rules from file failed, err: %v", err)
			return
		}
		global.LOG.Infof("loaded iptables rules for forward from file successfully")

		iptablesForwardStatus, _ := settingRepo.GetValueByKey("IptablesForwardStatus")
		if iptablesForwardStatus == constant.StatusEnable {
			if err := firewallClient.EnableIptablesForward(); err != nil {
				global.LOG.Errorf("enable iptables forward failed, err: %v", err)
				return
			}
		}
	}

	if clientName != "iptables" {
		return
	}
	if err := iptables.LoadRulesFromFile(iptables.FilterTab, iptables.ChainHiTechCloudBasicBefore, iptables.BasicBeforeFileName); err != nil {
		global.LOG.Errorf("load basic before rules from file failed, err: %v", err)
		return
	}
	if err := iptables.LoadRulesFromFile(iptables.FilterTab, iptables.ChainHiTechCloudBasic, iptables.BasicFileName); err != nil {
		global.LOG.Errorf("load basic rules from file failed, err: %v", err)
		return
	}
	if err := iptables.LoadRulesFromFile(iptables.FilterTab, iptables.ChainHiTechCloudBasicAfter, iptables.BasicAfterFileName); err != nil {
		global.LOG.Errorf("load basic after rules from file failed, err: %v", err)
		return
	}
	panelPort := service.LoadPanelPort()
	if len(panelPort) == 0 {
		global.LOG.Errorf("find HiTechCloud service port failed")
		return
	}
	if err := iptables.AddRule(iptables.FilterTab, iptables.ChainHiTechCloudBasicBefore, fmt.Sprintf("-p tcp -m tcp --dport %v -j ACCEPT", panelPort)); err != nil {
		global.LOG.Errorf("add port accept rule %v failed, err: %v", panelPort, err)
		return
	}
	global.LOG.Infof("loaded iptables rules for basic from file successfully")
	iptablesService := service.IptablesService{}
	iptablesStatus, _ := settingRepo.GetValueByKey("IptablesStatus")
	if iptablesStatus == constant.StatusEnable {
		if err := iptablesService.Operate(dto.IptablesOp{Operate: "bind-base-without-init"}); err != nil {
			global.LOG.Errorf("bind base chains failed, err: %v", err)
			return
		}
	}

	if err := iptables.LoadRulesFromFile(iptables.FilterTab, iptables.ChainHiTechCloudInput, iptables.InputFileName); err != nil {
		global.LOG.Errorf("load input rules from file failed, err: %v", err)
		return
	}
	if err := iptables.LoadRulesFromFile(iptables.FilterTab, iptables.ChainHiTechCloudOutput, iptables.OutputFileName); err != nil {
		global.LOG.Errorf("load output rules from file failed, err: %v", err)
		return
	}
	global.LOG.Infof("loaded iptables rules for input and output from file successfully")
	iptablesInputStatus, _ := settingRepo.GetValueByKey("IptablesInputStatus")
	if iptablesInputStatus == constant.StatusEnable {
		if err := iptablesService.Operate(dto.IptablesOp{Name: iptables.ChainHiTechCloudInput, Operate: "bind"}); err != nil {
			global.LOG.Errorf("bind input chains failed, err: %v", err)
			return
		}
	}
	iptablesOutputStatus, _ := settingRepo.GetValueByKey("IptablesOutputStatus")
	if iptablesOutputStatus == constant.StatusEnable {
		if err := iptablesService.Operate(dto.IptablesOp{Name: iptables.ChainHiTechCloudOutput, Operate: "bind"}); err != nil {
			global.LOG.Errorf("bind output chains failed, err: %v", err)
			return
		}
	}
}

func needInit() bool {
	file, err := os.OpenFile("/run/HiTechCloud_boot_mark", os.O_RDWR|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		if os.IsExist(err) {
			return false
		}
		global.LOG.Errorf("check boot mark file failed: %v", err)
		return true
	}
	defer file.Close()
	fmt.Fprintf(file, "Boot Mark for HiTechCloud\n")
	return true
}

func InitPingStatus() {
	global.LOG.Info("initializing ban ping status from settings...")
	status := firewall.LoadPingStatus()
	statusInDB, _ := repo.NewISettingRepo().GetValueByKey("BanPing")
	if statusInDB == status {
		return
	}

	enable := "1"
	if statusInDB == constant.StatusDisable {
		enable = "0"
	}
	if err := firewall.UpdatePingStatus(enable); err != nil {
		global.LOG.Errorf("initialize ping status failed: %v", err)
	}
}
