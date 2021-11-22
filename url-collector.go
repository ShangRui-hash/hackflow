package hackflow

import (
	"bufio"
	"go/build"
	"os/exec"
	"strings"

	"github.com/sirupsen/logrus"
)

type UrlCollector struct {
	name     string
	execPath string
}

func GetUrlCollector() (*UrlCollector, error) {
	tool := container.Get(URL_COLLECTOR)
	if tool == nil {
		tool = &UrlCollector{
			name: URL_COLLECTOR,
		}
		if err := tool.download(); err != nil {
			return nil, err
		}
		container.Set(tool)
	}
	return tool.(*UrlCollector), nil
}

//Name 返回工具名称
func (u *UrlCollector) Name() string {
	return u.name
}

//ExecPath 返回执行路径,如果没有,则下载
func (u *UrlCollector) ExecPath() (string, error) {
	if len(u.execPath) == 0 {
		if err := u.download(); err != nil {
			logrus.Errorf("download %s failed,err:", err)
			return "", err
		}
	}
	return u.execPath, nil
}

//donwload 下载工具
func (u *UrlCollector) download() error {
	if err := GetGo().Install("github.com/ShangRui-hash/url-collector"); err != nil {
		logrus.Error("download url-collector failed:", err)
		return err
	}
	u.execPath = build.Default.GOPATH + "/bin/url-collector"
	return nil
}

//UrlCollectorCofnig 工具配置
type UrlCollectorCofnig struct {
	RoutineCount int    `flag:"-r"`
	InputFile    string `flag:"-i"`
	SearchEngine string `flag:"-e"`
	Keyword      string `flag:"-k"`
	OuputFormat  string `flag:"-f"`
	Proxy        string `flag:"-p"`
}

//Run 运行工具
func (u *UrlCollector) Run(config *UrlCollectorCofnig) (urlCh chan string, err error) {
	execPath, err := u.ExecPath()
	if err != nil {
		return nil, err
	}
	cmd := exec.Command(execPath, parseConfig(*config)...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		logger.Error("cmd.StdoutPipe failed,err:", err)
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		logger.Error("Execute failed when Start:", err)
		return nil, err
	}
	urlCh = make(chan string, 1024)
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			if strings.HasPrefix(scanner.Text(), "http") {
				urlCh <- scanner.Text()
			}
		}
		close(urlCh)
	}()
	return urlCh, nil
}
