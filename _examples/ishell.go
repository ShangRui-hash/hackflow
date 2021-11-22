package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"

	"github.com/ShangRui-hash/hackflow"
	"github.com/abiosoft/ishell"
	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
)

func main() {
	shell := ishell.New()
	shell.SetPrompt("hackflow > ")
	tools := hackflow.GetAllTools()
	for i := range tools {
		shell.AddCmd(&ishell.Cmd{
			Name: tools[i].Name(),
			Help: tools[i].Desp(),
			Func: func(tool hackflow.Tool) func(c *ishell.Context) {
				return func(c *ishell.Context) {
					execPath, err := tool.ExecPath()
					if err != nil {
						logrus.Error("tools[i].ExecPath failed,err:", err)
						return
					}
					fmt.Println(c.RawArgs)
					if len(c.RawArgs) > 1 {
						CmdExec(execPath, c.RawArgs[1:]...)
					} else {
						CmdExec(execPath)
					}

				}
			}(tools[i]),
		})
	}
	//添加命令
	shell.AddCmd(&ishell.Cmd{
		Name: "tools",
		Help: "获取工具列表",
		Func: func(c *ishell.Context) {
			lines := [][]string{}
			tools := hackflow.GetAllTools()
			for i := range tools {
				lines = append(lines, []string{tools[i].Name(), tools[i].Desp()})
			}
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Name", "Desp"})
			table.SetAutoWrapText(false)
			table.SetAutoFormatHeaders(true)
			table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
			table.SetAlignment(tablewriter.ALIGN_LEFT)
			table.SetCenterSeparator("")
			table.SetColumnSeparator("")
			table.SetRowSeparator("")
			table.SetHeaderLine(false)
			table.SetBorder(false)
			table.SetTablePadding("\t") // pad with tabs
			table.SetNoWhiteSpace(true)
			table.AppendBulk(lines)
			table.Render()
			fmt.Println()
		},
	})
	shell.Run()
}

func CmdExec(name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func(ctx context.Context) {
	LOOP:
		for {
			select {
			case <-ctx.Done():
				break LOOP
			case <-c:
				cmd.Process.Release()
				cmd.Process.Kill()
				break LOOP
			}
		}
	}(ctx)
	if err := cmd.Run(); err != nil {
		logrus.Error("cmd.Run failed,err:", err)
		return
	}
}
