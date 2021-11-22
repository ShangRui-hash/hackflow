package main

import (
	"github.com/abiosoft/ishell"
)

func main() {
	shell := ishell.New()
	shell.AddCmd(&ishell.Cmd{
		Name: "tools",
		Help: "获取工具列表",
		Func: func(c *ishell.Context) {
			c.Println("tools")
		},
	})
	shell.Run()
}
