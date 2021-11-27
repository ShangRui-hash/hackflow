package hackflow

import (
	"os/exec"
)

type BaseTool struct {
	name     string
	desp     string
	execPath string
}

func (b *BaseTool) Name() string {
	return b.name
}

func (b *BaseTool) Desp() string {
	return b.desp
}

func (b *BaseTool) ExecPath(download func() (string, error)) (string, error) {
	if b.execPath != "" {
		return b.execPath, nil
	}
	execPath, err := exec.LookPath(b.name)
	if err != nil {
		logger.Debugf("exec.LookPath(%s) failed,err:%v", b.name, err)
		if execPath, err = download(); err != nil {
			logger.Errorf("download %s failed,err:%v", b.Name(), err)
			return "", err
		}
	}
	b.execPath = execPath
	return b.execPath, nil
}
