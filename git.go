package hackflow

import (
	"fmt"
	"os/exec"

	"github.com/serkanalgur/phpfuncs"
)

type Git struct {
	BaseTool
}

func newGit() Tool {
	return &Git{
		BaseTool: BaseTool{
			name:     "git",
			desp:     "git",
			execPath: "git",
		},
	}
}

//GetGit 获取Git对象
func GetGit() *Git {
	return container.Get(GIT).(*Git)
}

func (g *Git) ExecPath() (string, error) {
	return g.execPath, nil
}

func (g *Git) Download() (string, error) {
	return "", nil
}

type CloneConfig struct {
	Url      string
	SavePath string
	Depth    int
}

//Clone 克隆远程仓库
func (g *Git) Clone(config CloneConfig) error {
	args := []string{"clone"}
	if config.Depth != 0 {
		args = append(args, []string{"--depth", fmt.Sprintf("%d", config.Depth)}...)
	}
	if config.Url != "" {
		args = append(args, config.Url)
	}
	if config.SavePath != "" {
		args = append(args, config.SavePath)
		if phpfuncs.FileExists(config.SavePath) {
			return nil
		}
	}
	execPath, err := g.ExecPath()
	if err != nil {
		return err
	}
	output, err := exec.Command(execPath, args...).CombinedOutput()
	if err != nil {
		logger.Errorf("exec.Command failed,err:%v,output:%s", err, output)
		return err
	}
	return nil
}
