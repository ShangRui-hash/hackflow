package hackflow

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/cavaliercoder/grab"
	"github.com/serkanalgur/phpfuncs"
)

type Python struct {
	BaseTool
}

func newPython() Tool {
	return &Python{
		BaseTool: BaseTool{
			name: PYTHON,
			desp: "python 解释器",
		},
	}
}

func GetPython() *Python {
	return container.Get(PYTHON).(*Python)
}

func (p *Python) Download() (string, error) {
	url := "https://mirrors.huaweicloud.com/python"
	version := "3.9.9"
	var name string
	switch runtime.GOOS {
	case "darwin":
		name = fmt.Sprintf("python-%s-macos11.pkg", version)
	case "linux":
		name = fmt.Sprintf("Python-%s.tgz", version)
	case "windows":
		name = fmt.Sprintf("python-%s-amd64.exe", version)
	}
	url = fmt.Sprintf("%s/%s/%s", url, version, name)
	packagePath := filepath.Join(SavePath, name)
	if !phpfuncs.FileExists(packagePath) {
		// create client
		client := grab.NewClient()
		req, _ := grab.NewRequest(packagePath, url)

		// start download
		logger.Debugf("Downloading %v...\n", req.URL())
		resp := client.Do(req)
		logger.Debugf("  %v\n", resp.HTTPResponse.Status)

		// start UI loop
		t := time.NewTicker(500 * time.Millisecond)
		defer t.Stop()

	Loop:
		for {
			select {
			case <-t.C:
				logger.Debugf("\r  transferred %v / %v bytes (%.2f%%)\n",
					resp.BytesComplete(),
					resp.Size,
					100*resp.Progress())

			case <-resp.Done:
				// download is complete
				break Loop
			}
		}

		// check for errors
		if err := resp.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "Download failed: %v\n", err)
			os.Exit(1)
		}
		logger.Debugf("Download saved to %v \n", resp.Filename)
	}
	if err := os.Chmod(packagePath, 0777); err != nil {
		logger.Error("os.Chmod failed,err:", err)
		return "", nil
	}
	if err := exec.Command(packagePath).Run(); err != nil {
		logger.Error("exec.Command.Run failed,err:", err)
		return "", nil
	}
	return "python", nil
}

func (p *Python) ExecPath() (string, error) {
	return p.BaseTool.ExecPath(p.Download)
}

func (p *Python) Run(name string, args ...string) error {
	args = append([]string{name}, args...)
	if err := TryExec("python3", args...); err == nil {
		return CmdExec("python3", args...)
	}
	if err := TryExec("python", args...); err != nil {
		logger.Errorf("CmdExec python failed,err:%v,args:%v\n", err, args)
		return err
	}
	return CmdExec("python", args...)
}
