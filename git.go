package hackflow

import (
	"fmt"
	"os/exec"
)

type Git struct {
	name     string
	execPath string
}
type CloneConfig struct {
	Url      string
	SavePath string
	Depth    int
}

func (g *Git) Name() string {
	return g.name
}
func (g *Git) ExecPath() (string, error) {
	return g.execPath, nil
}
func (g *Git) download() error {
	return nil
}

//Clone 克隆远程仓库
func (g *Git) Clone(config CloneConfig) error {
	args := []string{"clone"}
	if config.Depth > 0 {
		args = append(args, "--depth", fmt.Sprintf("%d", config.Depth))
	}
	if config.Url != "" {
		args = append(args, config.Url)
	}
	if config.SavePath != "" {
		args = append(args, config.SavePath)
	}
	logger.Debugf("execPath:%s,args:%s", g.execPath, args)
	output, err := exec.Command(g.execPath, args...).CombinedOutput()
	if err != nil {
		logger.Errorf("exec.Command failed,err:%v,output:%s", err, output)
		return err
	}
	logger.Debugf("exec.Command success,output:%s", output)
	return nil
}

//GetGit 获取Git对象
func GetGit() *Git {
	if tool := container.Get(GIT); tool != nil {
		return tool.(*Git)
	}
	tool := &Git{
		name:     GIT,
		execPath: "git",
	}
	container.Set(tool)
	return tool
}
