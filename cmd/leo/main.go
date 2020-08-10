// Example of a daemon with echo service
package main

import (
	"leo/cmd/leo/command"

	"github.com/fatih/color"
	"github.com/jung-kurt/gofpdf"
)

func main() {
	color.Green("Hello")
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Hello, world")
	pdf.OutputFileAndClose("hello.pdf")
	if err := command.Command().Execute(); err != nil {
		color.Red("启动leo失败: %v\n", err)
	}
}
