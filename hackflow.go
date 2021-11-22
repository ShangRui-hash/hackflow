package hackflow

import (
	"fmt"
	"go/build"
	"reflect"

	"github.com/sirupsen/logrus"
)

const (
	GO            = "go"
	GIT           = "git"
	SQLMAP        = "sqlmap"
	URL_COLLECTOR = "url_collector"
	DIRSEARCH     = "dirsearch"
	ONE_FOR_ALL   = "one_for_all"
	KSUBDOMAIN    = "ksubdomain"
	SUBFINDER     = "subfinder"
	HTTPX         = "httpx"
)

var (
	container *Container
	logger    = logrus.New()
	SavePath  = build.Default.GOPATH + "/hackflow"
)

type Container struct {
	allTools []string
	tools    map[string]Tool
}

type Tool interface {
	Name() string
	ExecPath() (string, error)
	download() error
}

func init() {
	tools := make(map[string]Tool)
	container = &Container{
		tools: tools,
		allTools: []string{
			SQLMAP,
			URL_COLLECTOR,
			DIRSEARCH,
			ONE_FOR_ALL,
			KSUBDOMAIN,
			SUBFINDER,
			HTTPX,
		},
	}
}

func (c *Container) Set(tool Tool) {
	c.tools[tool.Name()] = tool
}

func (c *Container) Get(name string) Tool {
	tool, exist := c.tools[name]
	if !exist {
		return nil
	}
	return tool
}

//GetAllTools 获取全部工具
func GetAllTools() (tools []Tool) {
	for i := range container.allTools {
		tools = append(tools, container.Get(container.allTools[i]))
	}
	return tools
}

//SetDebug 是否开启debug模式
func SetDebug(debug bool) {
	if debug {
		logger.SetLevel(logrus.DebugLevel)
	} else {
		logger.SetLevel(logrus.InfoLevel)
	}
}

//parseConfig 解析config
func parseConfig(config interface{}) (args []string) {
	typeof := reflect.TypeOf(config)
	valueof := reflect.ValueOf(config)
	for i := 0; i < typeof.NumField(); i++ {
		field := typeof.Field(i)
		switch valueof.Field(i).Type().Kind() {
		case reflect.Bool:
			if valueof.Field(i).Bool() {
				args = append(args, field.Tag.Get("flag"))
			}
		case reflect.Int:
			if valueof.Field(i).Int() != 0 {
				args = append(args, []string{field.Tag.Get("flag"), fmt.Sprintf("%d", valueof.Field(i).Int())}...)
			}
		case reflect.String:
			if valueof.Field(i).String() != "" {
				args = append(args, []string{field.Tag.Get("flag"), valueof.Field(i).String()}...)
			}
		}
	}
	logger.Debug("parse config to args:", args)
	return args
}
