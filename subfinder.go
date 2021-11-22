package hackflow

import (
	"bufio"
	"fmt"
	"go/build"
	"os/exec"

	"github.com/sirupsen/logrus"
)

type Subfinder struct {
	name     string
	execPath string
}

// Name 获取工具名称
func (s *Subfinder) Name() string {
	return s.name
}

//ExecPath 获取工具执行路径,如果不存在则下载
func (s *Subfinder) ExecPath() (string, error) {
	if s.execPath == "" {
		if err := s.download(); err != nil {
			logger.Errorf("download %s failed,err:%v", s.Name(), err)
			return "", err
		}
	}
	return s.execPath, nil
}

func (s *Subfinder) download() error {
	if err := GetGo().Install("github.com/projectdiscovery/subfinder/v2/cmd/subfinder@latest"); err != nil {
		return err
	}
	s.execPath = build.Default.GOPATH + "/bin/subfinder"
	return nil
}

func GetSubfinder() *Subfinder {
	if tool := container.Get(SUBFINDER); tool != nil {
		return tool.(*Subfinder)
	}
	container.Set(&Subfinder{
		name: SUBFINDER,
	})
	return container.Get(SUBFINDER).(*Subfinder)
}

type SubfinderRunConfig struct {
	Proxy        string `flag:"-proxy"`
	Domain       string `flag:"-d"`
	RoutineCount int    `flag:"-t"`
	DomainCh     chan string
}

func (s *Subfinder) Run(config *SubfinderRunConfig) (subdomainCh chan string, err error) {
	args := append([]string{"-silent", "-nW"}, parseConfig(*config)...)
	execPath, err := s.ExecPath()
	if err != nil {
		return nil, err
	}
	cmd := exec.Command(execPath, args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		logrus.Error("cmd.StdoutPipe failed,err:", err)
		return nil, err
	}
	stdin, err := cmd.StdinPipe()
	if err != nil {
		logrus.Error("cmd.StdinPipe failed,err:", err)
		return nil, err
	}
	if config.DomainCh != nil {
		go func() {
			for domain := range config.DomainCh {
				fmt.Fprintln(stdin, domain)
			}
			stdin.Close()
		}()
	}
	if err := cmd.Start(); err != nil {
		logrus.Error("Execute failed when Start:", err)
		return nil, err
	}
	logger.Debugf("%s 启动成功\n", s.name)
	subdomainCh = make(chan string, 1024)
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			subdomainCh <- scanner.Text()
		}
		close(subdomainCh)
	}()
	return subdomainCh, nil
}
