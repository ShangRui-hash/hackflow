package hackflow

import (
	"bufio"
	"go/build"
	"os/exec"
	"strings"

	"github.com/sirupsen/logrus"
)

type UrlCollector struct {
	BaseTool
}

func newUrlCollector() Tool {
	return &UrlCollector{
		BaseTool{
			name: URL_COLLECTOR,
			desp: "谷歌、百度、必应搜索引擎采集工具",
		},
	}
}

func GetUrlCollector() *UrlCollector {
	return container.Get(URL_COLLECTOR).(*UrlCollector)
}

func (u *UrlCollector) ExecPath() (string, error) {
	return u.BaseTool.ExecPath(u.Download)
}

func (u *UrlCollector) Download() (string, error) {
	if err := GetGo().Install("github.com/ShangRui-hash/url-collector"); err != nil {
		logrus.Error("download url-collector failed:", err)
		return "", err
	}
	return build.Default.GOPATH + "/bin/url-collector", nil
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
