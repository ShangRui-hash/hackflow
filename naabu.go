package hackflow

import (
	"bufio"
	"fmt"
	"go/build"
	"os/exec"

	"github.com/sirupsen/logrus"
)

type Naabu struct {
	BaseTool
}

func newNaabu() Tool {
	return &Naabu{
		BaseTool{
			name: NAABU,
			desp: "端口扫描、服务识别",
		},
	}
}

func GetNaabu() *Naabu {
	return container.Get(NAABU).(*Naabu)
}

func (n *Naabu) ExecPath() (string, error) {
	return n.BaseTool.ExecPath(n.Download)
}

func (n *Naabu) Download() (string, error) {
	if err := GetGo().Install("github.com/projectdiscovery/naabu/v2/cmd/naabu@latest"); err != nil {
		logrus.Error("download naabu failed:", err)
		return "", err
	}
	return build.Default.GOPATH + "/bin/naabu", nil
}

//NabbuRunConfig 工具运行配置
type NabbuRunConfig struct {
	RoutineCount int `flag:"-c"`
	HostCh       chan string
}

func (n *Naabu) Run(config *NabbuRunConfig) (resultCh chan string, err error) {
	execPath, err := n.ExecPath()
	if err != nil {
		logger.Error("naabu exec path failed:", err)
		return
	}
	args := append([]string{"-silent", "-json"}, parseConfig(*config)...)
	cmd := exec.Command(execPath, args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		logger.Error("cmd.StdoutPipe failed,err:", err)
		return
	}
	//获取标准输入
	stdin, err := cmd.StdinPipe()
	if err != nil {
		logger.Error("cmd.StdinPipe failed,err:", err)
		return
	}
	if config.HostCh != nil {
		//写入标准输入
		go func() {
			for domain := range config.HostCh {
				fmt.Fprintln(stdin, domain)
			}
			stdin.Close()
		}()
	}
	//运行
	if err := cmd.Start(); err != nil {
		logger.Error("Execute failed when Start:" + err.Error())
		return nil, err
	}
	//输出
	resultCh = make(chan string, 1024)

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			resultCh <- scanner.Text()
		}
		close(resultCh)
	}()

	return nil, nil
}
