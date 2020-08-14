// Example of a daemon with echo service
package main

import (
	"leo/cmd/leo/command"
	"os"

	"github.com/fatih/color"
)

func main() {
	if len(os.Args) == 1 {
		os.Args = append(os.Args, "cal")
	}
	if err := command.Command().Execute(); err != nil {
		color.Red("启动leo失败: %v\n", err)
	}
}
