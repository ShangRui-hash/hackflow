package hackflow

import (
	"os"
	"os/exec"

	"github.com/sirupsen/logrus"
)

type Go struct {
	name     string
	execPath string
	desp     string
}

func newGo() Tool {
	return &Go{
		name:     "go",
		execPath: "go",
		desp:     "go 工具链",
	}
}

func (g *Go) Name() string {
	return g.name
}

func (g *Go) Desp() string {
	return g.desp
}

func (g *Go) ExecPath() (string, error) {
	return g.execPath, nil
}
func (g *Go) download() error {
	return nil
}

//Install go install
func (g *Go) Install(url string) error {
	logger.Debugf("go install %s ...\n", url)
	output, err := exec.Command(g.execPath, "install", "-v", url).CombinedOutput()
	if err != nil {
		logger.Errorf("go install %s error: %s", url, string(output))
		return err
	}
	logger.Debug("go install finished ", string(output))
	return nil
}

//Mod go mod
func (g *Go) Mod(path, name string) error {
	if err := os.Chdir(path); err != nil {
		logrus.Error("os.Chdir error:", err)
		return err
	}
	output, err := exec.Command(g.execPath, "mod", name).CombinedOutput()
	if err != nil {
		logrus.Errorf("go mod %s error: %s", name, string(output))
	}
	logrus.Debug(string(output))
	return nil
}

type BuildConfig struct {
	Path       string
	OutputFile string
	Files      []string
}

func (g *Go) Build(config BuildConfig) error {
	if err := os.Chdir(config.Path); err != nil {
		logrus.Error("os.Chdir error:", err)
		return err
	}
	args := []string{"build"}
	if config.OutputFile != "" {
		args = append(args, "-o", config.OutputFile)
	}
	if len(config.Files) > 0 {
		args = append(args, config.Files...)
	}
	output, err := exec.Command("go", args...).CombinedOutput()
	if err != nil {
		logrus.Errorf("go build %s error: %s", config.OutputFile, string(output))
		return err
	}
	logrus.Debug(string(output))
	return nil
}

//GetGo 获取go对象
func GetGo() *Go {
	tool := container.Get(GO)
	if tool == nil {
		tool = &Go{
			name:     GO,
			execPath: "go",
		}
		container.Set(tool)
	}
	return tool.(*Go)
}
