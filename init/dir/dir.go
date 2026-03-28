package dir

import (
	"path"

	"github.com/Pho-Tue-SoftWare-Solutions-JSC/hitechcloud-agent/global"
	"github.com/Pho-Tue-SoftWare-Solutions-JSC/hitechcloud-agent/utils/files"
)

func Init() {
	fileOp := files.NewFileOp()
	baseDir := global.CONF.Base.InstallDir
	_, _ = fileOp.CreateDirWithPath(true, path.Join(baseDir, "HiTechCloud/docker/compose/"))

	global.Dir.BaseDir, _ = fileOp.CreateDirWithPath(true, baseDir)
	global.Dir.DataDir, _ = fileOp.CreateDirWithPath(true, path.Join(baseDir, "HiTechCloud"))
	global.Dir.DbDir, _ = fileOp.CreateDirWithPath(true, path.Join(baseDir, "HiTechCloud/db"))
	global.Dir.LogDir, _ = fileOp.CreateDirWithPath(true, path.Join(baseDir, "HiTechCloud/log"))
	global.Dir.TaskDir, _ = fileOp.CreateDirWithPath(true, path.Join(baseDir, "HiTechCloud/log/task"))
	global.Dir.TmpDir, _ = fileOp.CreateDirWithPath(true, path.Join(baseDir, "HiTechCloud/tmp"))

	global.Dir.AppDir, _ = fileOp.CreateDirWithPath(true, path.Join(baseDir, "HiTechCloud/apps"))
	global.Dir.ResourceDir, _ = fileOp.CreateDirWithPath(true, path.Join(baseDir, "HiTechCloud/resource"))
	global.Dir.IconCacheDir, _ = fileOp.CreateDirWithPath(true, path.Join(baseDir, "HiTechCloud/resource/icon"))
	global.Dir.AppResourceDir, _ = fileOp.CreateDirWithPath(true, path.Join(baseDir, "HiTechCloud/resource/apps"))
	global.Dir.AppInstallDir, _ = fileOp.CreateDirWithPath(true, path.Join(baseDir, "HiTechCloud/apps"))
	global.Dir.LocalAppResourceDir, _ = fileOp.CreateDirWithPath(true, path.Join(baseDir, "HiTechCloud/resource/apps/local"))
	global.Dir.LocalAppInstallDir, _ = fileOp.CreateDirWithPath(true, path.Join(baseDir, "HiTechCloud/apps/local"))
	global.Dir.RemoteAppResourceDir, _ = fileOp.CreateDirWithPath(true, path.Join(baseDir, "HiTechCloud/resource/apps/remote"))
	global.Dir.CustomAppResourceDir, _ = fileOp.CreateDirWithPath(true, path.Join(baseDir, "HiTechCloud/resource/apps/custom"))
	global.Dir.OfflineAppResourceDir, _ = fileOp.CreateDirWithPath(true, path.Join(baseDir, "HiTechCloud/resource/offline"))
	global.Dir.RuntimeDir, _ = fileOp.CreateDirWithPath(true, path.Join(baseDir, "HiTechCloud/runtime"))
	global.Dir.RecycleBinDir, _ = fileOp.CreateDirWithPath(true, "/.HiTechCloud_clash")
	global.Dir.SSLLogDir, _ = fileOp.CreateDirWithPath(true, path.Join(baseDir, "HiTechCloud/log/ssl"))
	global.Dir.McpDir, _ = fileOp.CreateDirWithPath(true, path.Join(baseDir, "HiTechCloud/mcp"))
	global.Dir.ConvertLogDir, _ = fileOp.CreateDirWithPath(true, path.Join(baseDir, "HiTechCloud/log/convert"))
	global.Dir.TensorRTLLMDir, _ = fileOp.CreateDirWithPath(true, path.Join(baseDir, "HiTechCloud/ai/tensorrt_llm"))
	global.Dir.FirewallDir, _ = fileOp.CreateDirWithPath(true, path.Join(baseDir, "HiTechCloud/firewall"))
}
