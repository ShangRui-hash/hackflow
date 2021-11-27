package hackflow

import (
	"io"
	"os/exec"
	"regexp"
	"strings"

	"github.com/serkanalgur/phpfuncs"
)

type DirSearch struct {
	BaseTool
}

func newDirSearch() Tool {
	return &DirSearch{
		BaseTool: BaseTool{
			name: DIRSEARCH,
			desp: "目录扫描工具",
		},
	}
}

func GetDirSearch() *DirSearch {
	return container.Get(DIRSEARCH).(*DirSearch)
}

type DirSearchResult struct {
	URL         string
	RespSize    string
	RespCode    string
	RedirectURL string
}

var resultReg = regexp.MustCompile(`^\[\d{2}:\d{2}:\d{2}\]\s{1}(\d{3})\s{1}-\s{2}([0-9A-Z]{4,5})\s+-\s+(.+)`)

//ExecPath 返回执行路径
func (d *DirSearch) ExecPath() (string, error) {
	return d.BaseTool.ExecPath(d.Download)
}

//download 下载源码
func (d *DirSearch) Download() (string, error) {
	execPath := SavePath + "/dirsearch/dirsearch.py"
	if !phpfuncs.FileExists(d.execPath) {
		logger.Debug("正在下载dirsearch")
		output, err := exec.Command("pip3", "install", "dirsearch", "--target="+SavePath).CombinedOutput()
		if err != nil {
			logger.Error("pip3 install failed,err:", err, "output:", string(output))
			return "", err
		}
		logger.Debug("下载源码完成，开始下载依赖:", string(output))
		output, err = exec.Command("pip3", "install", "-r", SavePath+"/dirsearch/requirements.txt").CombinedOutput()
		if err != nil {
			logger.Error("pip3 install -r failed,err:", err, "output:", string(output))
			return "", err
		}
		logger.Debug("依赖下载完成:", string(output))
	}
	return execPath, nil
}

type DirSearchConfig struct {
	HTTPMethod      string   `flag:"-m"`
	FullURL         bool     `flag:"--full-url"`
	RandomAgent     bool     `flag:"--random-agent"`
	RemoveExtension bool     `flag:"--remove-extensions"`
	EXT             []string `flag:"-e"`
	Subdirs         []string `flag:"--subdirs"`
}

func (d *DirSearch) Run(reader io.Reader, config DirSearchConfig) (io.Reader, error) {
	execPath, err := d.ExecPath()
	if err != nil {
		return nil, err
	}
	args := []string{execPath, "--no-color", "--stdin"}
	args = append(args, parseConfig(config)...)
	cmd := exec.Command("python3", args...)
	cmd.Stdin = reader
	logger.Debug(cmd.String())
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	//3.执行命令
	go func() {
		if err := cmd.Start(); err != nil {
			logger.Error("Execute failed when Start:" + err.Error())
		}
		if err := cmd.Wait(); err != nil {
			logger.Error("Execute failed when Wait:" + err.Error())
		}
	}()

	//1.获取标准输出
	return stdout, nil
}

func (d *DirSearch) ParseResult(line string) (result *DirSearchResult) {
	if resultReg.MatchString(line) {
		submatch := resultReg.FindStringSubmatch(line)
		result = &DirSearchResult{
			RespSize: submatch[2],
			RespCode: submatch[1],
		}
		if strings.Contains(submatch[3], "->") {
			urls := strings.Split(submatch[3], "->")
			result.URL = urls[0]
			result.RedirectURL = urls[1]
		} else {
			result.URL = submatch[3]
		}
	}
	return result
}
